package maps

import (
	"github.com/BattlesnakeOfficial/rules"
)

// SetupBoard is a shortcut for looking up a map by ID and initializing a new board state with it.
func SetupBoard(mapID string, settings rules.Settings, width, height int, snakeIDs []string) (*rules.BoardState, error) {
	boardState := rules.NewBoardState(width, height)

	rules.InitializeSnakes(boardState, snakeIDs)

	gameMap, err := GetMap(mapID)
	if err != nil {
		return nil, err
	}

	editor := NewBoardStateEditor(boardState)

	err = gameMap.SetupBoard(boardState, settings, editor)
	if err != nil {
		return nil, err
	}

	return boardState, nil
}

// PreUpdateBoard updates a board state with a map.
func PreUpdateBoard(gameMap GameMap, previousBoardState *rules.BoardState, settings rules.Settings) (*rules.BoardState, error) {
	nextBoardState := previousBoardState.Clone()
	editor := NewBoardStateEditor(nextBoardState)

	err := gameMap.PreUpdateBoard(previousBoardState, settings, editor)
	if err != nil {
		return nil, err
	}

	return nextBoardState, nil
}

func PostUpdateBoard(gameMap GameMap, previousBoardState *rules.BoardState, settings rules.Settings) (*rules.BoardState, error) {
	nextBoardState := previousBoardState.Clone()
	editor := NewBoardStateEditor(nextBoardState)

	err := gameMap.PostUpdateBoard(previousBoardState, settings, editor)
	if err != nil {
		return nil, err
	}

	return nextBoardState, nil
}

// An implementation of GameMap that just does predetermined placements, for testing.
type StubMap struct {
	Id             string
	SnakePositions map[string]rules.Point
	Food           []rules.Point
	Hazards        []rules.Point
	Error          error
}

func (m StubMap) ID() string {
	return m.Id
}

func (StubMap) Meta() Metadata {
	return Metadata{}
}

func (m StubMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	if m.Error != nil {
		return m.Error
	}
	for _, snake := range initialBoardState.Snakes {
		head := m.SnakePositions[snake.ID]
		editor.PlaceSnake(snake.ID, []rules.Point{head, head, head}, 100)
	}
	for _, food := range m.Food {
		editor.AddFood(food)
	}
	for _, hazard := range m.Hazards {
		editor.AddHazard(hazard)
	}
	return nil
}

func (m StubMap) PreUpdateBoard(previousBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return nil
}

func (m StubMap) PostUpdateBoard(previousBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	if m.Error != nil {
		return m.Error
	}
	for _, food := range m.Food {
		editor.AddFood(food)
	}
	for _, hazard := range m.Hazards {
		editor.AddHazard(hazard)
	}
	return nil
}

// drawRing draws a ring of hazard points offset from the outer edge of the board
func drawRing(bw, bh, hOffset, vOffset int) ([]rules.Point, error) {
	if bw < 1 {
		return nil, rules.RulesetError("board width too small")
	}

	if bh < 1 {
		return nil, rules.RulesetError("board height too small")
	}

	if hOffset >= bw-1 {
		return nil, rules.RulesetError("horizontal offset too large")
	}

	if vOffset >= bh-1 {
		return nil, rules.RulesetError("vertical offset too large")
	}

	if hOffset < 1 {
		return nil, rules.RulesetError("horizontal offset too small")
	}

	if vOffset < 1 {
		return nil, rules.RulesetError("vertical offset too small")
	}

	// calculate the start/end point of the horizontal borders
	xStart := hOffset - 1
	xEnd := bw - hOffset

	// calculate start/end point of the vertical borders
	yStart := vOffset - 1
	yEnd := bh - vOffset

	// we can pre-determine how many points will be in the ring and allocate a slice of exactly that size
	numPoints := 2 * (xEnd - xStart + 1) // horizontal hazard points

	// Add vertical walls, if there are any.
	// Sometimes there are no vertical walls when the ring height is only 2.
	// In that case, the vertical walls are handled by the horizontal walls
	if yEnd >= yStart {
		numPoints += 2*(yEnd-yStart+1) - 4
	}

	hazards := make([]rules.Point, 0, numPoints)

	// draw horizontal walls
	for x := xStart; x <= xEnd; x++ {
		hazards = append(hazards,
			rules.Point{X: x, Y: yStart},
			rules.Point{X: x, Y: yEnd},
		)
	}

	// draw vertical walls, but don't include corners that the horizontal walls already included
	for y := yStart + 1; y <= yEnd-1; y++ {
		hazards = append(hazards,
			rules.Point{X: xStart, Y: y},
			rules.Point{X: xEnd, Y: y},
		)
	}

	return hazards, nil
}

func maxInt(n1 int, n ...int) int {
	max := n1
	for _, v := range n {
		if v > max {
			max = v
		}
	}

	return max
}

func isOnBoard(w, h, x, y int) bool {
	if x >= w || x < 0 {
		return false
	}

	if y >= h || y < 0 {
		return false
	}

	return true
}

func PlaceSnakesInQuadrants(rand rules.Rand, editor Editor, snakes []rules.Snake, quadrants [][]rules.Point) error {
	if len(quadrants) != 4 {
		return rules.RulesetError("invalid start point configuration - not divided into quadrants")
	}

	// make sure all quadrants have the same number of positions
	for i := 1; i < 4; i++ {
		if len(quadrants[i]) != len(quadrants[0]) {
			return rules.RulesetError("invalid start point configuration - quadrants aren't even")
		}
	}

	quads := make([]rules.RandomPositionBucket, 4)
	for i := 0; i < 4; i++ {
		quads[i].Fill(quadrants[i]...)
	}

	currentQuad := rand.Intn(4) // randomly pick a quadrant to start from

	// evenly distribute snakes across quadrants, randomly, by rotating through the quadrants
	for _, snake := range snakes {
		p, err := quads[currentQuad].Take(rand)
		if err != nil {
			return err
		}

		editor.PlaceSnake(snake.ID, []rules.Point{p, p, p}, rules.SnakeMaxHealth)

		currentQuad = (currentQuad + 1) % 4
	}

	return nil
}

func PlaceFoodFixed(rand rules.Rand, initialBoardState *rules.BoardState, editor Editor) error {
	width, height := initialBoardState.Width, initialBoardState.Height
	centerCoord := rules.Point{X: (width - 1) / 2, Y: (height - 1) / 2}

	isSmallBoard := width*height < rules.BoardSizeMedium*rules.BoardSizeMedium

	// Up to 4 snakes can be placed such that food is nearby on small boards.
	// Otherwise, we skip this and only try to place food in the center.
	snakeBodies := editor.SnakeBodies()
	if len(snakeBodies) <= 4 || !isSmallBoard {
		// Place 1 food within exactly 2 moves of each snake, but never towards the center or in a corner
		for _, snakeBody := range snakeBodies {
			snakeHead := snakeBody[0]
			possibleFoodLocations := []rules.Point{
				{X: snakeHead.X - 1, Y: snakeHead.Y - 1},
				{X: snakeHead.X - 1, Y: snakeHead.Y + 1},
				{X: snakeHead.X + 1, Y: snakeHead.Y - 1},
				{X: snakeHead.X + 1, Y: snakeHead.Y + 1},
			}

			// Remove any invalid/unwanted positions
			availableFoodLocations := []rules.Point{}
			for _, p := range possibleFoodLocations {

				// Don't place in the center
				if centerCoord == p {
					continue
				}

				// Ignore points already occupied by food or hazards
				if editor.IsOccupied(p, true, true, true) {
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
				if (p.X == 0 || p.X == (width-1)) && (p.Y == 0 || p.Y == (height-1)) {
					continue
				}

				availableFoodLocations = append(availableFoodLocations, p)
			}

			if len(availableFoodLocations) <= 0 {
				return rules.ErrorNoRoomForFood
			}

			// Select randomly from available locations
			placedFood := availableFoodLocations[rand.Intn(len(availableFoodLocations))]
			editor.AddFood(placedFood)
		}
	}

	// Finally, try to place 1 food in center of board for dramatic purposes
	if !editor.IsOccupied(centerCoord, true, true, true) {
		editor.AddFood(centerCoord)
	}

	return nil
}
