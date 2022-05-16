package maps

import "github.com/BattlesnakeOfficial/rules"

// SetupBoard is a shortcut for looking up a map by ID and initializing a new board state with it.
func SetupBoard(mapID string, settings rules.Settings, width, height int, snakeIDs []string) (*rules.BoardState, error) {
	boardState := rules.NewBoardState(int32(width), int32(height))

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
