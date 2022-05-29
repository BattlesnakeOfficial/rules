package maps

import (
	"github.com/BattlesnakeOfficial/rules"
)

type EmptyMap struct{}

func init() {
	globalRegistry.RegisterMap(EmptyMap{})
}

func (m EmptyMap) ID() string {
	return "empty"
}

func (m EmptyMap) Meta() Metadata {
	return Metadata{
		Name:        "Empty",
		Description: "Default snake placement with no food",
		Author:      "Battlesnake",
	}
}

func (m EmptyMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	rand := settings.GetRand(0)

	snakeIDs := make([]string, 0, len(initialBoardState.Snakes))
	for _, snake := range initialBoardState.Snakes {
		snakeIDs = append(snakeIDs, snake.ID)
	}

	tempBoardState := rules.NewBoardState(initialBoardState.Width, initialBoardState.Height)
	err := rules.PlaceSnakesAutomatically(rand, tempBoardState, snakeIDs)
	if err != nil {
		return err
	}

	// Copy snakes from temp board state
	for _, snake := range tempBoardState.Snakes {
		editor.PlaceSnake(snake.ID, snake.Body, snake.Health)
	}

	return nil
}

func (m EmptyMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return nil
}
