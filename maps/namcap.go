package maps

import (
	"github.com/BattlesnakeOfficial/rules"
)

type NamcapMap struct{}

func init() {
	globalRegistry.RegisterMap("namcap", NamcapMap{})
}

func (m NamcapMap) ID() string {
	return "namcap"
}

func (m NamcapMap) Meta() Metadata {
	return Metadata{
		Name:        "Namcap",
		Description: "Generic dot eating game with supernatural enemies",
		Author:      "Battlesnake",
	}
}

func (m NamcapMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	rand := settings.GetRand(0)

	if initialBoardState.Width != 19 && initialBoardState.Height != 21 {
		return rules.RulesetError("This map can only be played on a 19X21 board")
	}

	snakePositions := []rules.Point{
		{X: 1, Y: 1},
		{X: 1, Y: 19},
		{X: 17, Y: 1},
		{X: 17, Y: 19},
	}

	if len(initialBoardState.Snakes) > len(snakePositions) {
		return rules.ErrorTooManySnakes
	}

	rand.Shuffle(len(snakePositions), func(i int, j int) {
		snakePositions[i], snakePositions[j] = snakePositions[j], snakePositions[i]
	})

	for index, snake := range initialBoardState.Snakes {
		head := snakePositions[index]
		editor.PlaceSnake(snake.ID, []rules.Point{head, head, head}, snake.Health)
	}

	for _, hazard := range hazards {
		editor.AddHazard(hazard)
	}

	return nil
}

func (m NamcapMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return nil
}

var hazards []rules.Point = []rules.Point{
	{X: 0, Y: 20},
	{X: 1, Y: 20},
	{X: 2, Y: 20},
	{X: 3, Y: 20},
	{X: 4, Y: 20},
	{X: 5, Y: 20},
	{X: 6, Y: 20},
	{X: 7, Y: 20},
	{X: 8, Y: 20},
	{X: 9, Y: 20},
	{X: 10, Y: 20},
	{X: 11, Y: 20},
	{X: 12, Y: 20},
	{X: 13, Y: 20},
	{X: 14, Y: 20},
	{X: 15, Y: 20},
	{X: 16, Y: 20},
	{X: 17, Y: 20},
	{X: 18, Y: 20},
	{X: 0, Y: 19},
	{X: 9, Y: 19},
	{X: 18, Y: 19},
	{X: 0, Y: 18},
	{X: 2, Y: 18},
	{X: 3, Y: 18},
	{X: 5, Y: 18},
	{X: 6, Y: 18},
	{X: 7, Y: 18},
	{X: 9, Y: 18},
	{X: 11, Y: 18},
	{X: 12, Y: 18},
	{X: 13, Y: 18},
	{X: 15, Y: 18},
	{X: 16, Y: 18},
	{X: 18, Y: 18},
	{X: 0, Y: 17},
	{X: 18, Y: 17},
	{X: 0, Y: 16},
	{X: 2, Y: 16},
	{X: 3, Y: 16},
	{X: 5, Y: 16},
	{X: 7, Y: 16},
	{X: 8, Y: 16},
	{X: 9, Y: 16},
	{X: 10, Y: 16},
	{X: 11, Y: 16},
	{X: 13, Y: 16},
	{X: 15, Y: 16},
	{X: 16, Y: 16},
	{X: 18, Y: 16},
	{X: 0, Y: 15},
	{X: 5, Y: 15},
	{X: 9, Y: 15},
	{X: 13, Y: 15},
	{X: 18, Y: 15},
	{X: 0, Y: 14},
	{X: 1, Y: 14},
	{X: 2, Y: 14},
	{X: 3, Y: 14},
	{X: 5, Y: 14},
	{X: 6, Y: 14},
	{X: 7, Y: 14},
	{X: 9, Y: 14},
	{X: 11, Y: 14},
	{X: 12, Y: 14},
	{X: 13, Y: 14},
	{X: 15, Y: 14},
	{X: 16, Y: 14},
	{X: 17, Y: 14},
	{X: 18, Y: 14},
	{X: 3, Y: 13},
	{X: 5, Y: 13},
	{X: 13, Y: 13},
	{X: 15, Y: 13},
	{X: 0, Y: 12},
	{X: 1, Y: 12},
	{X: 2, Y: 12},
	{X: 3, Y: 12},
	{X: 5, Y: 12},
	{X: 7, Y: 12},
	{X: 8, Y: 12},
	{X: 10, Y: 12},
	{X: 11, Y: 12},
	{X: 13, Y: 12},
	{X: 15, Y: 12},
	{X: 16, Y: 12},
	{X: 17, Y: 12},
	{X: 18, Y: 12},
	{X: 7, Y: 11},
	{X: 11, Y: 11},
	{X: 0, Y: 10},
	{X: 1, Y: 10},
	{X: 2, Y: 10},
	{X: 3, Y: 10},
	{X: 5, Y: 10},
	{X: 7, Y: 10},
	{X: 8, Y: 10},
	{X: 9, Y: 10},
	{X: 10, Y: 10},
	{X: 11, Y: 10},
	{X: 13, Y: 10},
	{X: 15, Y: 10},
	{X: 16, Y: 10},
	{X: 17, Y: 10},
	{X: 18, Y: 10},
	{X: 3, Y: 9},
	{X: 5, Y: 9},
	{X: 13, Y: 9},
	{X: 15, Y: 9},
	{X: 0, Y: 8},
	{X: 1, Y: 8},
	{X: 2, Y: 8},
	{X: 3, Y: 8},
	{X: 5, Y: 8},
	{X: 7, Y: 8},
	{X: 8, Y: 8},
	{X: 9, Y: 8},
	{X: 10, Y: 8},
	{X: 11, Y: 8},
	{X: 13, Y: 8},
	{X: 15, Y: 8},
	{X: 16, Y: 8},
	{X: 17, Y: 8},
	{X: 18, Y: 8},
	{X: 0, Y: 7},
	{X: 9, Y: 7},
	{X: 18, Y: 7},
	{X: 0, Y: 6},
	{X: 2, Y: 6},
	{X: 3, Y: 6},
	{X: 5, Y: 6},
	{X: 6, Y: 6},
	{X: 7, Y: 6},
	{X: 9, Y: 6},
	{X: 11, Y: 6},
	{X: 12, Y: 6},
	{X: 13, Y: 6},
	{X: 15, Y: 6},
	{X: 16, Y: 6},
	{X: 18, Y: 6},
	{X: 0, Y: 5},
	{X: 3, Y: 5},
	{X: 15, Y: 5},
	{X: 18, Y: 5},
	{X: 0, Y: 4},
	{X: 1, Y: 4},
	{X: 3, Y: 4},
	{X: 5, Y: 4},
	{X: 7, Y: 4},
	{X: 8, Y: 4},
	{X: 9, Y: 4},
	{X: 10, Y: 4},
	{X: 11, Y: 4},
	{X: 13, Y: 4},
	{X: 15, Y: 4},
	{X: 17, Y: 4},
	{X: 18, Y: 4},
	{X: 0, Y: 3},
	{X: 5, Y: 3},
	{X: 9, Y: 3},
	{X: 13, Y: 3},
	{X: 18, Y: 3},
	{X: 0, Y: 2},
	{X: 2, Y: 2},
	{X: 3, Y: 2},
	{X: 4, Y: 2},
	{X: 5, Y: 2},
	{X: 6, Y: 2},
	{X: 7, Y: 2},
	{X: 9, Y: 2},
	{X: 11, Y: 2},
	{X: 12, Y: 2},
	{X: 13, Y: 2},
	{X: 14, Y: 2},
	{X: 15, Y: 2},
	{X: 16, Y: 2},
	{X: 18, Y: 2},
	{X: 0, Y: 1},
	{X: 18, Y: 1},
	{X: 0, Y: 0},
	{X: 1, Y: 0},
	{X: 2, Y: 0},
	{X: 3, Y: 0},
	{X: 4, Y: 0},
	{X: 5, Y: 0},
	{X: 6, Y: 0},
	{X: 7, Y: 0},
	{X: 8, Y: 0},
	{X: 9, Y: 0},
	{X: 10, Y: 0},
	{X: 11, Y: 0},
	{X: 12, Y: 0},
	{X: 13, Y: 0},
	{X: 14, Y: 0},
	{X: 15, Y: 0},
	{X: 16, Y: 0},
	{X: 17, Y: 0},
	{X: 18, Y: 0},
}
