package maps

import (
	"github.com/BattlesnakeOfficial/rules"
)

type StandardMap struct{}

func init() {
	globalRegistry.RegisterMap("standard", StandardMap{})
}

func (m StandardMap) ID() string {
	return "standard"
}

func (m StandardMap) Meta() Metadata {
	return Metadata{
		Name:        "Standard",
		Description: "Standard snake placement and food spawning",
		Author:      "Battlesnake",
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  8,
		BoardSizes:  AnySize(),
	}
}

func (m StandardMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	rand := settings.GetRand(0)

	snakeIDs := make([]string, 0, len(initialBoardState.Snakes))
	for _, snake := range initialBoardState.Snakes {
		snakeIDs = append(snakeIDs, snake.ID)
	}

	tempBoardState, err := rules.CreateDefaultBoardState(rand, initialBoardState.Width, initialBoardState.Height, snakeIDs)
	if err != nil {
		return err
	}

	// Copy food from temp board state
	for _, food := range tempBoardState.Food {
		editor.AddFood(food)
	}

	// Copy snakes from temp board state
	for _, snake := range tempBoardState.Snakes {
		editor.PlaceSnake(snake.ID, snake.Body, snake.Health)
	}

	return nil
}

func (m StandardMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	rand := settings.GetRand(lastBoardState.Turn)

	foodNeeded := checkFoodNeedingPlacement(rand, settings, lastBoardState)
	if foodNeeded > 0 {
		placeFoodRandomly(rand, lastBoardState, editor, foodNeeded)
	}

	return nil
}

func checkFoodNeedingPlacement(rand rules.Rand, settings rules.Settings, state *rules.BoardState) int {
	minFood := int(settings.MinimumFood)
	foodSpawnChance := int(settings.FoodSpawnChance)
	numCurrentFood := len(state.Food)

	if numCurrentFood < minFood {
		return minFood - numCurrentFood
	}
	if foodSpawnChance > 0 && (100-rand.Intn(100)) < foodSpawnChance {
		return 1
	}

	return 0
}

func placeFoodRandomly(rand rules.Rand, b *rules.BoardState, editor Editor, n int) {
	unoccupiedPoints := rules.GetUnoccupiedPoints(b, false)
	placeFoodRandomlyAtPositions(rand, b, editor, n, unoccupiedPoints)
}

func placeFoodRandomlyAtPositions(rand rules.Rand, b *rules.BoardState, editor Editor, n int, positions []rules.Point) {
	if len(positions) < n {
		n = len(positions)
	}

	rand.Shuffle(len(positions), func(i int, j int) {
		positions[i], positions[j] = positions[j], positions[i]
	})

	for i := 0; i < n; i++ {
		editor.AddFood(positions[i])
	}
}
