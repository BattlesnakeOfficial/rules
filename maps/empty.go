package maps

import (
	"github.com/BattlesnakeOfficial/rules"
)

type EmptyMap struct{}

func init() {
	globalRegistry.RegisterMap("empty", EmptyMap{})
}

func (m EmptyMap) ID() string {
	return "empty"
}

func (m EmptyMap) Meta() Metadata {
	return Metadata{
		Name:        "Empty",
		Description: "Default snake placement with no food",
		Author:      "Battlesnake",
		Version:     2,
		MinPlayers:  1,
		MaxPlayers:  16,
		BoardSizes:  OddSizes(rules.BoardSizeSmall, rules.BoardSizeXXLarge),
		Tags:        []string{},
	}
}

func (m EmptyMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	rand := settings.GetRand(0)

	if len(initialBoardState.Snakes) > int(m.Meta().MaxPlayers) {
		return rules.ErrorTooManySnakes
	}

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

func (m EmptyMap) PreUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return nil
}

func (m EmptyMap) PostUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return nil
}
