package maps

import (
	"github.com/BattlesnakeOfficial/rules"
)

type TestHealthMap struct{}

func init() {
	globalRegistry.RegisterMap("test_health", TestHealthMap{})
}

func (m TestHealthMap) ID() string {
	return "test_health"
}

func (m TestHealthMap) Meta() Metadata {
	return Metadata{
		Name:        "Empty",
		Description: "Default snake placement with no food",
		Author:      "Battlesnake",
		Version:     2,
		MinPlayers:  1,
		MaxPlayers:  16,
		BoardSizes:  OddSizes(rules.BoardSizeSmall, rules.BoardSizeXXLarge),
	}
}

func (m TestHealthMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
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

func (m TestHealthMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	if lastBoardState.Turn == 30 {
		snake := lastBoardState.Snakes[0]
		editor.PlaceSnake(snake.ID, snake.Body, 0)
	}

	return nil
}

func (m TestHealthMap) PreUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	if lastBoardState.Turn == 20 {
		snake := lastBoardState.Snakes[0]
		head := snake.Body[0]

		editor.AddHazard(rules.Point{X: head.X + 1, Y: head.Y})
		editor.AddHazard(rules.Point{X: head.X - 1, Y: head.Y})
		editor.AddHazard(rules.Point{X: head.X, Y: head.Y + 1})
		editor.AddHazard(rules.Point{X: head.X, Y: head.Y - 1})
	}

	return nil
}
