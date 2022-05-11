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

	editor := NewBoardStateEditor(boardState, settings.Rand())

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
	editor := NewBoardStateEditor(nextBoardState, settings.Rand())

	err = gameMap.SetupBoard(previousBoardState, settings, editor)
	if err != nil {
		return nil, err
	}

	return nextBoardState, nil
}
