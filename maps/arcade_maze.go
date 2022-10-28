package maps

import (
	"github.com/BattlesnakeOfficial/rules"
)

type ArcadeMazeMap struct{}

func init() {
	globalRegistry.RegisterMap("arcade_maze", ArcadeMazeMap{})
}

func (m ArcadeMazeMap) ID() string {
	return "arcade_maze"
}

func (m ArcadeMazeMap) Meta() Metadata {
	return Metadata{
		Name:        "Arcade Maze",
		Description: "Generic arcade maze map with deadly hazard walls.",
		Author:      "Battlesnake",
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  6,
		BoardSizes:  FixedSizes(Dimensions{19, 21}),
		Tags:        []string{TAG_FOOD_PLACEMENT, TAG_HAZARD_PLACEMENT, TAG_SNAKE_PLACEMENT},
	}
}

func (m ArcadeMazeMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	rand := settings.GetRand(0)

	if initialBoardState.Width != 19 || initialBoardState.Height != 21 {
		return rules.RulesetError("This map can only be played on a 19X21 board")
	}

	// Shuffle the first four starting locations
	snakePositions := []rules.Point{
		{X: 4, Y: 7},
		{X: 14, Y: 7},
		{X: 4, Y: 17},
		{X: 14, Y: 17},
	}
	rand.Shuffle(len(snakePositions), func(i int, j int) {
		snakePositions[i], snakePositions[j] = snakePositions[j], snakePositions[i]
	})

	// Add a fifth and sixth starting location that are always placed last
	snakePositions = append(snakePositions, rules.Point{X: 9, Y: 9})
	snakePositions = append(snakePositions, rules.Point{X: 9, Y: 13})

	// Place snakes
	if len(initialBoardState.Snakes) > len(snakePositions) {
		return rules.ErrorTooManySnakes
	}
	for index, snake := range initialBoardState.Snakes {
		head := snakePositions[index]
		editor.PlaceSnake(snake.ID, []rules.Point{head, head, head}, snake.Health)
	}

	// Place static hazards
	for _, hazard := range ArcadeMazeHazards {
		editor.AddHazard(hazard)
	}

	if settings.Int(rules.ParamMinimumFood, 0) > 0 {
		// Add food in center
		editor.AddFood(rules.Point{X: 9, Y: 11})
	}

	return nil
}

func (m ArcadeMazeMap) PreUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return nil
}

func (m ArcadeMazeMap) PostUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	rand := settings.GetRand(lastBoardState.Turn)

	// Respect FoodSpawnChance setting
	foodSpawnChance := settings.Int(rules.ParamFoodSpawnChance, 0)
	if foodSpawnChance == 0 || rand.Intn(100) > foodSpawnChance {
		return nil
	}

	foodPositions := []rules.Point{
		{X: 1, Y: 1},
		{X: 3, Y: 11},
		{X: 4, Y: 7},
		{X: 4, Y: 17},
		{X: 9, Y: 1},
		{X: 9, Y: 5},
		{X: 9, Y: 11},
		{X: 9, Y: 17},
		{X: 14, Y: 7},
		{X: 14, Y: 17},
		{X: 15, Y: 11},
		{X: 17, Y: 1},
	}

	rand.Shuffle(len(foodPositions), func(i int, j int) {
		foodPositions[i], foodPositions[j] = foodPositions[j], foodPositions[i]
	})

	for _, food := range foodPositions {
		tileIsOccupied := false

	snakeLoop:
		for _, snake := range lastBoardState.Snakes {
			for _, point := range snake.Body {
				if point.X == food.X && point.Y == food.Y {
					tileIsOccupied = true
					break snakeLoop
				}
			}
		}

		for _, existingFood := range lastBoardState.Food {
			if existingFood.X == food.X && existingFood.Y == food.Y {
				tileIsOccupied = true
				break
			}
		}

		if !tileIsOccupied {
			editor.AddFood(food)
			break
		}
	}

	return nil
}

var ArcadeMazeHazards []rules.Point = []rules.Point{
	{X: 0, Y: 20},
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
	{X: 3, Y: 14},
	{X: 5, Y: 14},
	{X: 6, Y: 14},
	{X: 7, Y: 14},
	{X: 9, Y: 14},
	{X: 11, Y: 14},
	{X: 12, Y: 14},
	{X: 13, Y: 14},
	{X: 15, Y: 14},
	{X: 18, Y: 14},
	{X: 0, Y: 13},
	{X: 3, Y: 13},
	{X: 5, Y: 13},
	{X: 13, Y: 13},
	{X: 15, Y: 13},
	{X: 18, Y: 13},
	{X: 0, Y: 12},
	{X: 1, Y: 12},
	{X: 2, Y: 12},
	{X: 3, Y: 12},
	{X: 5, Y: 12},
	{X: 7, Y: 12},
	{X: 9, Y: 12},
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
	{X: 9, Y: 10},
	{X: 11, Y: 10},
	{X: 13, Y: 10},
	{X: 15, Y: 10},
	{X: 16, Y: 10},
	{X: 17, Y: 10},
	{X: 18, Y: 10},
	{X: 0, Y: 9},
	{X: 3, Y: 9},
	{X: 5, Y: 9},
	{X: 13, Y: 9},
	{X: 15, Y: 9},
	{X: 18, Y: 9},
	{X: 0, Y: 8},
	{X: 3, Y: 8},
	{X: 5, Y: 8},
	{X: 7, Y: 8},
	{X: 8, Y: 8},
	{X: 9, Y: 8},
	{X: 10, Y: 8},
	{X: 11, Y: 8},
	{X: 13, Y: 8},
	{X: 15, Y: 8},
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
	{X: 18, Y: 0},
}
