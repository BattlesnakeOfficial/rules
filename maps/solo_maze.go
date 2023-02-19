package maps

import (
	"fmt"
	"strconv"

	"bufio"
	"bytes"
	"log"
	"os"

	"github.com/BattlesnakeOfficial/rules"
)

// When this is flipped to `true` TWO things happen
//  1. More println style debugging is done
//  2. We print out the current game board in between each room sub-division,
//     and wait for the CLI User to hit enter to sub-divide the next room. This
//     allows you to see the maze get generated in realtime, which was super useful
//     while debugging issues in the maze generation
const DEBUG_MAZE_GENERATION = false

const INITIAL_MAZE_SIZE = 7

const TURNS_AT_MAX_SIZE = 5

const EVIL_MODE_DISTANCE_TO_FOOD = 5

const MAX_TRIES = 100

type SoloMazeMap struct{}

func init() {
	mazeMap := SoloMazeMap{}
	globalRegistry.RegisterMap(mazeMap.ID(), mazeMap)
}

func (m SoloMazeMap) ID() string {
	return "solo_maze"
}

func (m SoloMazeMap) Meta() Metadata {
	return Metadata{
		Name:        "Solo Maze",
		Description: "Solo Maze where you need to find the food",
		Author:      "coreyja",
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  1,
		BoardSizes: FixedSizes(
			Dimensions{7, 7},
			Dimensions{11, 11},
			Dimensions{19, 19},
			Dimensions{19, 21},
			Dimensions{25, 25},
		),
		Tags: []string{TAG_EXPERIMENTAL, TAG_FOOD_PLACEMENT, TAG_HAZARD_PLACEMENT, TAG_SNAKE_PLACEMENT},
	}
}

func (m SoloMazeMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	if len(initialBoardState.Snakes) != 1 {
		return rules.RulesetError("This map requires exactly one snake")
	}

	if initialBoardState.Width < INITIAL_MAZE_SIZE || initialBoardState.Height < INITIAL_MAZE_SIZE {
		return rules.RulesetError(
			fmt.Sprintf("This map requires a board size of at least %dx%d", INITIAL_MAZE_SIZE, INITIAL_MAZE_SIZE))
	}

	return m.CreateMaze(initialBoardState, settings, editor, 0)
}

func maxBoardSize(boardState *rules.BoardState) int {
	return min(boardState.Width, boardState.Height-2)
}

func gameNeedsToEndSoon(maxBoardSize int, currentLevel int64) bool {
	return currentLevel-TURNS_AT_MAX_SIZE > int64(maxBoardSize-INITIAL_MAZE_SIZE)
}

func (m SoloMazeMap) CreateMaze(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor, currentLevel int64) error {
	rand := settings.GetRand(initialBoardState.Turn)

	// Make sure the actual maze size can always fit in the CreateBoard
	// This means that when you get to 'max' size each level stops making
	// the maze bigger
	actualBoardSize := INITIAL_MAZE_SIZE + currentLevel
	maxBoardSize := maxBoardSize(initialBoardState)

	if actualBoardSize > int64(maxBoardSize) {
		actualBoardSize = int64(maxBoardSize)
	}

	me := initialBoardState.Snakes[0]

	mazeBoardState := rules.NewBoardState(int(actualBoardSize), int(actualBoardSize))
	tempBoardState := initialBoardState.Clone()

	topRightCorner := rules.Point{X: int(actualBoardSize) - 1, Y: int(actualBoardSize) - 1}

	editor.ClearHazards()

	m.WriteBitState(initialBoardState, currentLevel, editor)

	m.SubdivideRoom(mazeBoardState, rand, rules.Point{X: 0, Y: 0}, topRightCorner, make([]int, 0), make([]int, 0), 0)

	for _, point := range removeDuplicateValues(mazeBoardState.Hazards) {
		adjusted := m.AdjustPosition(point, int(actualBoardSize), initialBoardState.Height, initialBoardState.Width)
		editor.AddHazard(adjusted)
		tempBoardState.Hazards = append(tempBoardState.Hazards, adjusted)
	}

	// Since we reserve the bottom row of the board for state,
	// AND we center the maze within the board we know there will
	// always be a `y: -1` that we can put the tail into
	snake_head_position := rules.Point{X: 0, Y: 0}
	snake_tail_position := rules.Point{X: 0, Y: -1}

	snakeBody := []rules.Point{
		snake_head_position,
	}
	for i := 0; i <= int(currentLevel)+1; i++ {
		snakeBody = append(snakeBody, snake_tail_position)
	}

	adjustedSnakeBody := make([]rules.Point, len(snakeBody))
	for i, point := range snakeBody {
		adjustedSnakeBody[i] = m.AdjustPosition(point, int(actualBoardSize), initialBoardState.Height, initialBoardState.Width)
	}
	editor.PlaceSnake(me.ID, adjustedSnakeBody, 100)
	tempBoardState.Snakes[0].Body = adjustedSnakeBody

	/// Pick random food spawn point
	m.PlaceFood(tempBoardState, settings, editor, currentLevel)

	// Fill outside of the board with walls
	xAdjust := int((initialBoardState.Width - int(actualBoardSize)) / 2)
	yAdjust := int((initialBoardState.Height - int(actualBoardSize)) / 2)
	for x := 0; x < initialBoardState.Width; x++ {
		for y := 1; y < initialBoardState.Height; y++ {
			if x < xAdjust || y < yAdjust || x >= xAdjust+int(actualBoardSize) || y >= yAdjust+int(actualBoardSize) {
				editor.AddHazard(rules.Point{X: x, Y: y})
			}
		}
	}

	return nil
}

func (m SoloMazeMap) PlaceFood(boardState *rules.BoardState, settings rules.Settings, editor Editor, currentLevel int64) {
	actualBoardSize := INITIAL_MAZE_SIZE + currentLevel
	maxBoardSize := maxBoardSize(boardState)
	if actualBoardSize > int64(maxBoardSize) {
		actualBoardSize = int64(maxBoardSize)
	}
	meBody := boardState.Snakes[0].Body
	myHead := meBody[0]

	foodPlaced := false
	tries := 0
	// We want to place a random food, but we also want an escape hatch for if the algo gets stuck in a loop
	// trying to place a food.
	for !foodPlaced && tries < MAX_TRIES {
		tries++
		rand := settings.GetRand(boardState.Turn + tries)

		foodSpawnPoint := rules.Point{X: rand.Intn(int(actualBoardSize)), Y: rand.Intn(int(actualBoardSize))}
		adjustedFood := m.AdjustPosition(foodSpawnPoint, int(actualBoardSize), boardState.Height, boardState.Width)

		minDistanceFromFood := min(EVIL_MODE_DISTANCE_TO_FOOD, int(actualBoardSize/2))
		if !containsPoint(boardState.Hazards, adjustedFood) && !containsPoint(meBody, adjustedFood) && manhattanDistance(adjustedFood, myHead) >= minDistanceFromFood {
			editor.AddFood(adjustedFood)
			foodPlaced = true
		}
	}
}

func (m SoloMazeMap) PreUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return nil
}

func (m SoloMazeMap) PostUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	currentLevel, e := m.ReadBitState(lastBoardState)
	if e != nil {
		return e
	}

	if len(lastBoardState.Food) == 0 {
		currentLevel += 1
		m.WriteBitState(lastBoardState, currentLevel, editor)

		// This will create a new maze
		return m.CreateMaze(lastBoardState, settings, editor, currentLevel)
	}

	maxBoardSize := maxBoardSize(lastBoardState)

	food := lastBoardState.Food[0]

	meBody := lastBoardState.Snakes[0].Body
	myHead := meBody[0]

	if gameNeedsToEndSoon(maxBoardSize, currentLevel) && manhattanDistance(myHead, food) < EVIL_MODE_DISTANCE_TO_FOOD {
		editor.RemoveFood(food)

		m.PlaceFood(lastBoardState, settings, editor, currentLevel)
	}

	return nil
}

// Mostly based off this algorithm from Wikipedia: https://en.wikipedia.org/wiki/Maze_generation_algorithm#Recursive_division_method
func (m SoloMazeMap) SubdivideRoom(tempBoardState *rules.BoardState, rand rules.Rand, lowPoint rules.Point, highPoint rules.Point, disAllowedHorizontal []int, disAllowedVertical []int, depth int) bool {
	didSubdivide := false

	if DEBUG_MAZE_GENERATION {
		log.Print("\n\n\n")
		log.Printf("Subdividing room from %v to %v", lowPoint, highPoint)
		log.Printf("disAllowedVertical %v", disAllowedVertical)
		log.Printf("disAllowedHorizontal %v", disAllowedHorizontal)
		printMap(tempBoardState)
		fmt.Print("Press 'Enter' to continue...")
		_, e := bufio.NewReader(os.Stdin).ReadBytes('\n')
		if e != nil {
			log.Fatal(e)
		}
	}

	verticalWallPosition := -1
	horizontalWallPosition := -1
	newVerticalWall := make([]rules.Point, 0)
	newHorizontalWall := make([]rules.Point, 0)

	if highPoint.X-lowPoint.X <= 2 && highPoint.Y-lowPoint.Y <= 2 {
		return false
	}

	verticalChoices := make([]int, 0)
	for i := lowPoint.X + 1; i < highPoint.X-1; i++ {
		if !contains(disAllowedVertical, i) {
			verticalChoices = append(verticalChoices, i)
		}
	}
	if len(verticalChoices) > 0 {
		verticalWallPosition = verticalChoices[rand.Intn(len(verticalChoices))]
		if DEBUG_MAZE_GENERATION {
			log.Printf("drawing Vertical Wall at %v\n", verticalWallPosition)
		}

		for y := lowPoint.Y; y <= highPoint.Y; y++ {
			newVerticalWall = append(newVerticalWall, rules.Point{X: verticalWallPosition, Y: y})
		}

		didSubdivide = true
	}

	/// We can only draw a horizontal wall if there is enough space
	horizontalChoices := make([]int, 0)
	for i := lowPoint.Y + 1; i < highPoint.Y-1; i++ {
		if !contains(disAllowedHorizontal, i) {
			horizontalChoices = append(horizontalChoices, i)
		}
	}
	if len(horizontalChoices) > 0 {
		horizontalWallPosition = horizontalChoices[rand.Intn(len(horizontalChoices))]
		if DEBUG_MAZE_GENERATION {
			log.Printf("drawing horizontal Wall at %v\n", horizontalWallPosition)
		}

		for x := lowPoint.X; x <= highPoint.X; x++ {
			newHorizontalWall = append(newHorizontalWall, rules.Point{X: x, Y: horizontalWallPosition})
		}

		didSubdivide = true
	}

	/// Here we make cuts in the walls
	if len(newVerticalWall) > 1 && len(newHorizontalWall) > 1 {
		if DEBUG_MAZE_GENERATION {
			log.Print("Need to cut with both walls")
		}
		intersectionPoint := rules.Point{X: verticalWallPosition, Y: horizontalWallPosition}

		newNewVerticalWall, verticalHoles := cutHoles(newVerticalWall, intersectionPoint, rand)
		newVerticalWall = newNewVerticalWall

		for _, hole := range verticalHoles {
			disAllowedHorizontal = append(disAllowedHorizontal, hole.Y)
		}
		if DEBUG_MAZE_GENERATION {
			log.Printf("Vertical Cuts are at %v\n", verticalHoles)
		}

		newNewHorizontalWall, horizontalHoles := cutHoleSingle(newHorizontalWall, intersectionPoint, rand)
		newHorizontalWall = newNewHorizontalWall
		for _, hole := range horizontalHoles {
			disAllowedVertical = append(disAllowedVertical, hole.X)
		}
		if DEBUG_MAZE_GENERATION {
			log.Printf("Horizontal Cuts are at %v\n", horizontalHoles)
		}
	} else if len(newVerticalWall) > 1 {
		if DEBUG_MAZE_GENERATION {
			log.Print("Only a vertical wall needs cut")
		}
		segmentToRemove := rand.Intn(len(newVerticalWall) - 1)
		hole := newVerticalWall[segmentToRemove]
		newVerticalWall = remove(newVerticalWall, segmentToRemove)

		disAllowedHorizontal = append(disAllowedHorizontal, hole.Y)
		if DEBUG_MAZE_GENERATION {
			log.Printf("Cuts are at %v from index %v", hole, segmentToRemove)
		}
	} else if len(newHorizontalWall) > 1 {
		if DEBUG_MAZE_GENERATION {
			log.Print("Only a horizontal wall needs cut")
		}
		segmentToRemove := rand.Intn(len(newHorizontalWall) - 1)
		hole := newHorizontalWall[segmentToRemove]
		newHorizontalWall = remove(newHorizontalWall, segmentToRemove)

		disAllowedVertical = append(disAllowedVertical, hole.X)
		if DEBUG_MAZE_GENERATION {
			log.Printf("Cuts are at %v from index %v", hole, segmentToRemove)
		}
	}

	tempBoardState.Hazards = append(tempBoardState.Hazards, newVerticalWall...)
	tempBoardState.Hazards = append(tempBoardState.Hazards, newHorizontalWall...)

	/// We have both so need 4 sub-rooms
	if verticalWallPosition != -1 && horizontalWallPosition != -1 {
		m.SubdivideRoom(tempBoardState, rand, rules.Point{X: lowPoint.X, Y: lowPoint.Y}, rules.Point{X: verticalWallPosition, Y: horizontalWallPosition}, disAllowedHorizontal, disAllowedVertical, depth+1)
		m.SubdivideRoom(tempBoardState, rand, rules.Point{X: verticalWallPosition + 1, Y: lowPoint.Y}, rules.Point{X: highPoint.X, Y: horizontalWallPosition}, disAllowedHorizontal, disAllowedVertical, depth+1)

		m.SubdivideRoom(tempBoardState, rand, rules.Point{X: lowPoint.X, Y: horizontalWallPosition + 1}, rules.Point{X: verticalWallPosition, Y: highPoint.Y}, disAllowedHorizontal, disAllowedVertical, depth+1)
		m.SubdivideRoom(tempBoardState, rand, rules.Point{X: verticalWallPosition + 1, Y: horizontalWallPosition + 1}, rules.Point{X: highPoint.X, Y: highPoint.Y}, disAllowedHorizontal, disAllowedVertical, depth+1)
	} else if verticalWallPosition != -1 {
		m.SubdivideRoom(tempBoardState, rand, rules.Point{X: lowPoint.X, Y: lowPoint.Y}, rules.Point{X: verticalWallPosition, Y: highPoint.Y}, disAllowedHorizontal, disAllowedVertical, depth+1)
		m.SubdivideRoom(tempBoardState, rand, rules.Point{X: verticalWallPosition + 1, Y: lowPoint.Y}, rules.Point{X: highPoint.X, Y: highPoint.Y}, disAllowedHorizontal, disAllowedVertical, depth+1)
	} else if horizontalWallPosition != -1 {
		m.SubdivideRoom(tempBoardState, rand, rules.Point{X: lowPoint.X, Y: lowPoint.Y}, rules.Point{X: highPoint.X, Y: horizontalWallPosition}, disAllowedHorizontal, disAllowedVertical, depth+1)
		m.SubdivideRoom(tempBoardState, rand, rules.Point{X: lowPoint.X, Y: horizontalWallPosition + 1}, rules.Point{X: highPoint.X, Y: highPoint.Y}, disAllowedHorizontal, disAllowedVertical, depth+1)
	}

	return didSubdivide
}

//////// Maze Helpers ////////

func (m SoloMazeMap) AdjustPosition(mazePosition rules.Point, actualBoardSize int, boardHeight int, boardWidth int) rules.Point {

	xAdjust := int((boardWidth - actualBoardSize) / 2)
	yAdjust := int((boardHeight - actualBoardSize) / 2)

	if DEBUG_MAZE_GENERATION {
		fmt.Printf("currentLevel: %v, boardHeight: %v, boardWidth: %v, xAdjust: %v, yAdjust: %v\n", actualBoardSize, boardHeight, boardWidth, xAdjust, yAdjust)
	}

	return rules.Point{X: mazePosition.X + xAdjust, Y: mazePosition.Y + yAdjust}
}

func (m SoloMazeMap) ReadBitState(boardState *rules.BoardState) (int64, error) {
	row := 0
	width := boardState.Width

	stringBits := ""

	for i := 0; i < width; i++ {
		if containsPoint(boardState.Hazards, rules.Point{X: i, Y: row}) {
			stringBits += "1"
		} else {
			stringBits += "0"
		}
	}

	return strconv.ParseInt(stringBits, 2, 64)
}

func (m SoloMazeMap) WriteBitState(boardState *rules.BoardState, state int64, editor Editor) {
	width := boardState.Width

	stringBits := strconv.FormatInt(state, 2)
	paddingBits := fmt.Sprintf("%0*s", width, stringBits)

	for i, c := range paddingBits {
		point := rules.Point{X: i, Y: 0}

		if c == '1' {
			editor.AddHazard(point)
		} else {
			editor.RemoveHazard(point)
		}
	}
}

// Return value is first the wall that has been cut, the second is the holes we cut out
func cutHoles(s []rules.Point, intersection rules.Point, rand rules.Rand) ([]rules.Point, []rules.Point) {
	holes := make([]rules.Point, 0)

	index := pos(s, intersection)

	if index != 0 {
		firstSegmentToRemove := rand.Intn(index)
		holes = append(holes, s[firstSegmentToRemove])
		s = remove(s, firstSegmentToRemove)
	}

	index = pos(s, intersection)
	if index != len(s)-1 {
		secondSegmentToRemove := rand.Intn(len(s)-index-1) + index + 1

		holes = append(holes, s[secondSegmentToRemove])
		s = remove(s, secondSegmentToRemove)
	}

	return s, holes
}

func cutHoleSingle(s []rules.Point, intersection rules.Point, rand rules.Rand) ([]rules.Point, []rules.Point) {
	holes := make([]rules.Point, 0)

	index := pos(s, intersection)

	if index != 0 {
		firstSegmentToRemove := rand.Intn(index)
		holes = append(holes, s[firstSegmentToRemove])
		s = remove(s, firstSegmentToRemove)

		return s, holes
	}

	if index != len(s)-1 {
		secondSegmentToRemove := rand.Intn(len(s)-index-1) + index + 1

		holes = append(holes, s[secondSegmentToRemove])
		s = remove(s, secondSegmentToRemove)
	}

	return s, holes
}

//////// Golang Helpers ////////

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsPoint(s []rules.Point, e rules.Point) bool {
	for _, a := range s {
		if a.X == e.X && a.Y == e.Y {
			return true
		}
	}
	return false
}

func remove(s []rules.Point, i int) []rules.Point {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func pos(s []rules.Point, e rules.Point) int {
	for i, a := range s {
		if a == e {
			return i
		}
	}
	return -1
}

func removeDuplicateValues(hazards []rules.Point) []rules.Point {
	keys := make(map[rules.Point]bool)
	uniqueList := []rules.Point{}

	for _, entry := range hazards {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueList = append(uniqueList, entry)
		}
	}
	return uniqueList
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func manhattanDistance(a, b rules.Point) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

//// DEBUGING HELPERS ////

// This mostly copy pasted from the CLI which prints out the boardState
// Copied here to not create a circular dependency
// Removed some of the color picking logic to simplify things
func printMap(boardState *rules.BoardState) {
	var o bytes.Buffer
	o.WriteString(fmt.Sprintf("Turn: %v\n", boardState.Turn))
	board := make([][]string, boardState.Width)
	for i := range board {
		board[i] = make([]string, boardState.Height)
	}
	for y := int(0); y < boardState.Height; y++ {
		for x := int(0); x < boardState.Width; x++ {
			board[x][y] = "◦"
		}
	}
	for _, oob := range boardState.Hazards {
		board[oob.X][oob.Y] = "░"
	}
	for _, f := range boardState.Food {
		board[f.X][f.Y] = "⚕"
	}
	o.WriteString(fmt.Sprintf("Food ⚕: %v\n", boardState.Food))
	for _, s := range boardState.Snakes {
		for _, b := range s.Body {
			if b.X >= 0 && b.X < boardState.Width && b.Y >= 0 && b.Y < boardState.Height {
				board[b.X][b.Y] = string("*")
			}
		}
	}
	for y := boardState.Height - 1; y >= 0; y-- {
		for x := int(0); x < boardState.Width; x++ {
			o.WriteString(board[x][y])
		}
		o.WriteString("\n")
	}
	log.Print(o.String())
}
