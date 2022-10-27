package rules

import "fmt"

// BoardState represents the internal state of a game board.
// NOTE: use NewBoardState to construct these to ensure fields are initialized
// correctly and that tests are resilient to changes to this type.
type BoardState struct {
	Turn    int
	Height  int
	Width   int
	Food    []Point
	Snakes  []Snake
	Hazards []Point

	// Generic game-level state for maps and rules stages to persist data between turns.
	GameState map[string]string

	// Numeric state keyed to specific points, also persisted between turns.
	PointState map[Point]int
}

type Point struct {
	X     int `json:"X"`
	Y     int `json:"Y"`
	TTL   int `json:"TTL,omitempty"`
	Value int `json:"Value,omitempty"`
}

// Makes it easier to copy sample points out of Go logs and test failures.
func (p Point) GoString() string {
	if p.TTL != 0 || p.Value != 0 {
		return fmt.Sprintf("{X:%d, Y:%d, TTL:%d, Value:%d}", p.X, p.Y, p.TTL, p.Value)
	}
	return fmt.Sprintf("{X:%d, Y:%d}", p.X, p.Y)
}

type Snake struct {
	ID               string
	Body             []Point
	Health           int
	EliminatedCause  string
	EliminatedOnTurn int
	EliminatedBy     string
}

// NewBoardState returns an empty but fully initialized BoardState
func NewBoardState(width, height int) *BoardState {
	return &BoardState{
		Turn:       0,
		Height:     height,
		Width:      width,
		Food:       []Point{},
		Snakes:     []Snake{},
		Hazards:    []Point{},
		GameState:  map[string]string{},
		PointState: map[Point]int{},
	}
}

// Clone returns a deep copy of prevState that can be safely modified without affecting the original
func (prevState *BoardState) Clone() *BoardState {
	nextState := &BoardState{
		Turn:       prevState.Turn,
		Height:     prevState.Height,
		Width:      prevState.Width,
		Food:       append([]Point{}, prevState.Food...),
		Snakes:     make([]Snake, len(prevState.Snakes)),
		Hazards:    append([]Point{}, prevState.Hazards...),
		GameState:  make(map[string]string, len(prevState.GameState)),
		PointState: make(map[Point]int, len(prevState.PointState)),
	}
	for key, value := range prevState.GameState {
		nextState.GameState[key] = value
	}
	for key, value := range prevState.PointState {
		nextState.PointState[key] = value
	}
	for i := 0; i < len(prevState.Snakes); i++ {
		nextState.Snakes[i].ID = prevState.Snakes[i].ID
		nextState.Snakes[i].Health = prevState.Snakes[i].Health
		nextState.Snakes[i].Body = append([]Point{}, prevState.Snakes[i].Body...)
		nextState.Snakes[i].EliminatedCause = prevState.Snakes[i].EliminatedCause
		nextState.Snakes[i].EliminatedOnTurn = prevState.Snakes[i].EliminatedOnTurn
		nextState.Snakes[i].EliminatedBy = prevState.Snakes[i].EliminatedBy
	}
	return nextState
}

// Builder method to set Turn and return the modified BoardState.
func (state *BoardState) WithTurn(turn int) *BoardState {
	state.Turn = turn
	return state
}

// Builder method to set Food and return the modified BoardState.
func (state *BoardState) WithFood(food []Point) *BoardState {
	state.Food = food
	return state
}

// Builder method to set Hazards and return the modified BoardState.
func (state *BoardState) WithHazards(hazards []Point) *BoardState {
	state.Hazards = hazards
	return state
}

// Builder method to set Snakes and return the modified BoardState.
func (state *BoardState) WithSnakes(snakes []Snake) *BoardState {
	state.Snakes = snakes
	return state
}

// Builder method to set State and return the modified BoardState.
func (state *BoardState) WithGameState(gameState map[string]string) *BoardState {
	state.GameState = gameState
	return state
}

// Builder method to set PointState and return the modified BoardState.
func (state *BoardState) WithPointState(pointState map[Point]int) *BoardState {
	state.PointState = pointState
	return state
}

// CreateDefaultBoardState is a convenience function for fully initializing a
// "default" board state with snakes and food.
// In a real game, the engine may generate the board without calling this
// function, or customize the results based on game-specific settings.
func CreateDefaultBoardState(rand Rand, width int, height int, snakeIDs []string) (*BoardState, error) {
	initialBoardState := NewBoardState(width, height)

	err := PlaceSnakesAutomatically(rand, initialBoardState, snakeIDs)
	if err != nil {
		return nil, err
	}

	err = PlaceFoodAutomatically(rand, initialBoardState)
	if err != nil {
		return nil, err
	}

	return initialBoardState, nil
}

// PlaceSnakesAutomatically initializes the array of snakes based on the provided snake IDs and the size of the board.
func PlaceSnakesAutomatically(rand Rand, b *BoardState, snakeIDs []string) error {

	if isSquareBoard(b) {
		// we don't allow > 8 snakes on very small boards
		if len(snakeIDs) > 8 && b.Width < BoardSizeSmall {
			return ErrorTooManySnakes
		}

		// we can do fixed placement for up to 8 snakes on minimum sized boards
		if len(snakeIDs) <= 8 && b.Width >= BoardSizeSmall {
			return PlaceSnakesFixed(rand, b, snakeIDs)
		}

		// for > 8 snakes, we can do distributed placement
		if b.Width >= BoardSizeMedium {
			return PlaceManySnakesDistributed(rand, b, snakeIDs)
		}
	}

	// last resort for unexpected board sizes we'll just randomly place snakes
	return PlaceSnakesRandomly(rand, b, snakeIDs)
}

func PlaceSnakesFixed(rand Rand, b *BoardState, snakeIDs []string) error {
	b.Snakes = make([]Snake, len(snakeIDs))

	for i := 0; i < len(snakeIDs); i++ {
		b.Snakes[i] = Snake{
			ID:     snakeIDs[i],
			Health: SnakeMaxHealth,
		}
	}

	// Create start 8 points
	mn, md, mx := 1, (b.Width-1)/2, b.Width-2
	cornerPoints := []Point{
		{X: mn, Y: mn},
		{X: mn, Y: mx},
		{X: mx, Y: mn},
		{X: mx, Y: mx},
	}
	cardinalPoints := []Point{
		{X: mn, Y: md},
		{X: md, Y: mn},
		{X: md, Y: mx},
		{X: mx, Y: md},
	}

	// Sanity check
	if len(b.Snakes) > (len(cornerPoints) + len(cardinalPoints)) {
		return ErrorTooManySnakes
	}

	// Randomly order them
	rand.Shuffle(len(cornerPoints), func(i int, j int) {
		cornerPoints[i], cornerPoints[j] = cornerPoints[j], cornerPoints[i]
	})
	rand.Shuffle(len(cardinalPoints), func(i int, j int) {
		cardinalPoints[i], cardinalPoints[j] = cardinalPoints[j], cardinalPoints[i]
	})

	var startPoints []Point
	if rand.Intn(2) == 0 {
		startPoints = append(startPoints, cornerPoints...)
		startPoints = append(startPoints, cardinalPoints...)
	} else {
		startPoints = append(startPoints, cardinalPoints...)
		startPoints = append(startPoints, cornerPoints...)
	}

	// Assign to snakes in order given
	for i := 0; i < len(b.Snakes); i++ {
		for j := 0; j < SnakeStartSize; j++ {
			b.Snakes[i].Body = append(b.Snakes[i].Body, startPoints[i])
		}

	}

	return nil
}

// PlaceManySnakesDistributed is a placement algorithm that works for up to 16 snakes
// It is intended for use on large boards and distributes snakes relatively evenly,
// and randomly, across quadrants.
func PlaceManySnakesDistributed(rand Rand, b *BoardState, snakeIDs []string) error {
	// this placement algorithm supports up to 16 snakes
	if len(snakeIDs) > 16 {
		return ErrorTooManySnakes
	}

	b.Snakes = make([]Snake, len(snakeIDs))

	for i := 0; i < len(snakeIDs); i++ {
		b.Snakes[i] = Snake{
			ID:     snakeIDs[i],
			Health: SnakeMaxHealth,
		}
	}

	quadHSpace := b.Width / 2
	quadVSpace := b.Height / 2

	hOffset := quadHSpace / 3
	vOffset := quadVSpace / 3

	quads := make([]RandomPositionBucket, 4)

	// quad 1
	quads[0] = RandomPositionBucket{}
	quads[0].Fill(
		Point{X: hOffset, Y: vOffset},
		Point{X: quadHSpace - hOffset, Y: vOffset},
		Point{X: hOffset, Y: quadVSpace - vOffset},
		Point{X: quadHSpace - hOffset, Y: quadVSpace - vOffset},
	)

	// quad 2
	quads[1] = RandomPositionBucket{}
	for _, p := range quads[0].positions {
		quads[1].Fill(Point{X: b.Width - p.X - 1, Y: p.Y})
	}

	// quad 3
	quads[2] = RandomPositionBucket{}
	for _, p := range quads[0].positions {
		quads[2].Fill(Point{X: p.X, Y: b.Height - p.Y - 1})
	}

	// quad 4
	quads[3] = RandomPositionBucket{}
	for _, p := range quads[0].positions {
		quads[3].Fill(Point{X: b.Width - p.X - 1, Y: b.Height - p.Y - 1})
	}

	currentQuad := rand.Intn(4) // randomly pick a quadrant to start from
	// evenly distribute snakes across quadrants, randomly, by rotating through the quadrants
	for i := 0; i < len(b.Snakes); i++ {
		p, err := quads[currentQuad].Take(rand)
		if err != nil {
			return err
		}
		for j := 0; j < SnakeStartSize; j++ {
			b.Snakes[i].Body = append(b.Snakes[i].Body, p)
		}

		currentQuad = (currentQuad + 1) % 4
	}

	return nil
}

type RandomPositionBucket struct {
	positions []Point
}

func (rpb *RandomPositionBucket) Fill(p ...Point) {
	rpb.positions = append(rpb.positions, p...)
}

func (rpb *RandomPositionBucket) Take(rand Rand) (Point, error) {
	if len(rpb.positions) == 0 {
		return Point{}, RulesetError("no more positions available")
	}

	// randomly pick the next position
	idx := rand.Intn(len(rpb.positions))
	p := rpb.positions[idx]

	// remove that position from the list using the fast slice removal method
	rpb.positions[idx] = rpb.positions[len(rpb.positions)-1]
	rpb.positions = rpb.positions[:len(rpb.positions)-1]

	return p, nil
}

func PlaceSnakesRandomly(rand Rand, b *BoardState, snakeIDs []string) error {
	b.Snakes = make([]Snake, len(snakeIDs))

	for i := 0; i < len(snakeIDs); i++ {
		b.Snakes[i] = Snake{
			ID:     snakeIDs[i],
			Health: SnakeMaxHealth,
		}
	}

	for i := 0; i < len(b.Snakes); i++ {
		unoccupiedPoints := removeCenterCoord(b, GetEvenUnoccupiedPoints(b))
		if len(unoccupiedPoints) <= 0 {
			return ErrorNoRoomForSnake
		}
		p := unoccupiedPoints[rand.Intn(len(unoccupiedPoints))]
		for j := 0; j < SnakeStartSize; j++ {
			b.Snakes[i].Body = append(b.Snakes[i].Body, p)
		}
	}
	return nil
}

// Adds all snakes without body coordinates to the board.
// This allows GameMaps to access the list of snakes and perform initial placement.
func InitializeSnakes(b *BoardState, snakeIDs []string) {
	b.Snakes = make([]Snake, len(snakeIDs))

	for i := 0; i < len(snakeIDs); i++ {
		b.Snakes[i] = Snake{
			ID:     snakeIDs[i],
			Health: SnakeMaxHealth,
			Body:   []Point{},
		}
	}
}

// PlaceSnake adds a snake to the board with the given ID and body coordinates.
func PlaceSnake(b *BoardState, snakeID string, body []Point) error {
	// Update an existing snake that already has a body
	for index, snake := range b.Snakes {
		if snake.ID == snakeID {
			b.Snakes[index].Body = body
			return nil
		}
	}
	// Add a new snake
	b.Snakes = append(b.Snakes, Snake{
		ID:     snakeID,
		Health: SnakeMaxHealth,
		Body:   body,
	})
	return nil
}

// PlaceFoodAutomatically initializes the array of food based on the size of the board and the number of snakes.
func PlaceFoodAutomatically(rand Rand, b *BoardState) error {
	if isSquareBoard(b) && b.Width >= BoardSizeSmall {
		return PlaceFoodFixed(rand, b)
	}

	return PlaceFoodRandomly(rand, b, len(b.Snakes))
}

// Deprecated: will be replaced by maps.PlaceFoodFixed
func PlaceFoodFixed(rand Rand, b *BoardState) error {
	centerCoord := Point{X: (b.Width - 1) / 2, Y: (b.Height - 1) / 2}

	isSmallBoard := b.Width*b.Height < BoardSizeMedium*BoardSizeMedium
	// Up to 4 snakes can be placed such that food is nearby on small boards.
	// Otherwise, we skip this and only try to place food in the center.
	if len(b.Snakes) <= 4 || !isSmallBoard {
		// Place 1 food within exactly 2 moves of each snake, but never towards the center or in a corner
		for i := 0; i < len(b.Snakes); i++ {
			snakeHead := b.Snakes[i].Body[0]
			possibleFoodLocations := []Point{
				{X: snakeHead.X - 1, Y: snakeHead.Y - 1},
				{X: snakeHead.X - 1, Y: snakeHead.Y + 1},
				{X: snakeHead.X + 1, Y: snakeHead.Y - 1},
				{X: snakeHead.X + 1, Y: snakeHead.Y + 1},
			}

			// Remove any invalid/unwanted positions
			availableFoodLocations := []Point{}
			for _, p := range possibleFoodLocations {

				// Don't place in the center
				if centerCoord == p {
					continue
				}

				// Ignore points already occupied by food
				isOccupiedAlready := false
				for _, food := range b.Food {
					if food.X == p.X && food.Y == p.Y {
						isOccupiedAlready = true
						break
					}
				}
				if isOccupiedAlready {
					continue
				}

				// Food must be further than snake from center on at least one axis
				isAwayFromCenter := false
				if p.X < snakeHead.X && snakeHead.X < centerCoord.X {
					isAwayFromCenter = true
				} else if centerCoord.X < snakeHead.X && snakeHead.X < p.X {
					isAwayFromCenter = true
				} else if p.Y < snakeHead.Y && snakeHead.Y < centerCoord.Y {
					isAwayFromCenter = true
				} else if centerCoord.Y < snakeHead.Y && snakeHead.Y < p.Y {
					isAwayFromCenter = true
				}
				if !isAwayFromCenter {
					continue
				}

				// Don't spawn food in corners
				if (p.X == 0 || p.X == (b.Width-1)) && (p.Y == 0 || p.Y == (b.Height-1)) {
					continue
				}

				availableFoodLocations = append(availableFoodLocations, p)
			}

			if len(availableFoodLocations) <= 0 {
				return ErrorNoRoomForFood
			}

			// Select randomly from available locations
			placedFood := availableFoodLocations[rand.Intn(len(availableFoodLocations))]
			b.Food = append(b.Food, placedFood)
		}
	}

	// Finally, always place 1 food in center of board for dramatic purposes
	isCenterOccupied := true
	unoccupiedPoints := GetUnoccupiedPoints(b, true, false)
	for _, point := range unoccupiedPoints {
		if point == centerCoord {
			isCenterOccupied = false
			break
		}
	}
	if isCenterOccupied {
		return ErrorNoRoomForFood
	}
	b.Food = append(b.Food, centerCoord)

	return nil
}

// PlaceFoodRandomly adds up to n new food to the board in random unoccupied squares
func PlaceFoodRandomly(rand Rand, b *BoardState, n int) error {
	for i := 0; i < n; i++ {
		unoccupiedPoints := GetUnoccupiedPoints(b, false, false)
		if len(unoccupiedPoints) > 0 {
			newFood := unoccupiedPoints[rand.Intn(len(unoccupiedPoints))]
			b.Food = append(b.Food, newFood)
		}
	}
	return nil
}

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func GetEvenUnoccupiedPoints(b *BoardState) []Point {
	// Start by getting unoccupied points
	unoccupiedPoints := GetUnoccupiedPoints(b, true, false)

	// Create a new array to hold points that are  even
	evenUnoccupiedPoints := []Point{}

	for _, point := range unoccupiedPoints {
		if ((point.X + point.Y) % 2) == 0 {
			evenUnoccupiedPoints = append(evenUnoccupiedPoints, point)
		}
	}
	return evenUnoccupiedPoints
}

// removeCenterCoord filters out the board's center point from a list of points.
func removeCenterCoord(b *BoardState, points []Point) []Point {
	centerCoord := Point{X: (b.Width - 1) / 2, Y: (b.Height - 1) / 2}
	var noCenterPoints []Point
	for _, p := range points {
		if p != centerCoord {
			noCenterPoints = append(noCenterPoints, p)
		}
	}

	return noCenterPoints
}

func GetUnoccupiedPoints(b *BoardState, includePossibleMoves bool, includeHazards bool) []Point {
	pointIsOccupied := map[int]map[int]bool{}
	for _, p := range b.Food {
		if _, xExists := pointIsOccupied[p.X]; !xExists {
			pointIsOccupied[p.X] = map[int]bool{}
		}
		pointIsOccupied[p.X][p.Y] = true
	}

	for _, snake := range b.Snakes {
		if snake.EliminatedCause != NotEliminated {
			continue
		}
		for i, p := range snake.Body {
			if _, xExists := pointIsOccupied[p.X]; !xExists {
				pointIsOccupied[p.X] = map[int]bool{}
			}
			pointIsOccupied[p.X][p.Y] = true

			if i == 0 && !includePossibleMoves {
				nextMovePoints := []Point{
					{X: p.X - 1, Y: p.Y},
					{X: p.X + 1, Y: p.Y},
					{X: p.X, Y: p.Y - 1},
					{X: p.X, Y: p.Y + 1},
				}
				for _, nextP := range nextMovePoints {
					if _, xExists := pointIsOccupied[nextP.X]; !xExists {
						pointIsOccupied[nextP.X] = map[int]bool{}
					}
					pointIsOccupied[nextP.X][nextP.Y] = true
				}
			}
		}
	}

	if includeHazards {
		for _, p := range b.Hazards {
			if _, xExists := pointIsOccupied[p.X]; !xExists {
				pointIsOccupied[p.X] = map[int]bool{}
			}
			pointIsOccupied[p.X][p.Y] = true
		}
	}

	unoccupiedPoints := []Point{}
	for x := 0; x < b.Width; x++ {
		for y := 0; y < b.Height; y++ {
			if _, xExists := pointIsOccupied[x]; xExists {
				if isOccupied, yExists := pointIsOccupied[x][y]; yExists {
					if isOccupied {
						continue
					}
				}
			}
			unoccupiedPoints = append(unoccupiedPoints, Point{X: x, Y: y})
		}
	}
	return unoccupiedPoints
}

func getDistanceBetweenPoints(a, b Point) int {
	return absInt(a.X-b.X) + absInt(a.Y-b.Y)
}

func isSquareBoard(b *BoardState) bool {
	return b.Width == b.Height
}

// EliminateSnake updates a snake's state to reflect that it was eliminated.
// - "cause" identifies what type of event caused the snake to be eliminated
// - "by" identifies which snake (if any, use empty string "" if none) eliminated the snake.
// - "turn" is the turn number that this snake was eliminated on.
func EliminateSnake(s *Snake, cause, by string, turn int) {
	s.EliminatedCause = cause
	s.EliminatedBy = by
	s.EliminatedOnTurn = turn
}
