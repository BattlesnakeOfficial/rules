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

//  UpdateBoard is a shortcut for looking up a map by ID and updating an existing board state with it.
func PreUpdateBoard(mapID string, previousBoardState *rules.BoardState, settings rules.Settings) (*rules.BoardState, error) {
	gameMap, err := GetMap(mapID)
	if err != nil {
		return nil, err
	}

	nextBoardState := previousBoardState.Clone()
	editor := NewBoardStateEditor(nextBoardState)

	err = gameMap.PreUpdateBoard(previousBoardState, settings, editor)
	if err != nil {
		return nil, err
	}

	return nextBoardState, nil
}

//  UpdateBoard is a shortcut for looking up a map by ID and updating an existing board state with it.
func UpdateBoard(mapID string, previousBoardState *rules.BoardState, settings rules.Settings) (*rules.BoardState, error) {
	gameMap, err := GetMap(mapID)
	if err != nil {
		return nil, err
	}

	nextBoardState := previousBoardState.Clone()
	editor := NewBoardStateEditor(nextBoardState)

	err = gameMap.UpdateBoard(previousBoardState, settings, editor)
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

func (m StubMap) UpdateBoard(previousBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
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
