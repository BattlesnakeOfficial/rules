package maps

import (
	"github.com/BattlesnakeOfficial/rules"
)

type StandardMap struct{}

func init() {
	globalRegistry.RegisterMap(StandardMap{})
}

func (m StandardMap) ID() string {
	return "standard"
}

func (m StandardMap) Meta() Metadata {
	return Metadata{
		Name:        "Standard",
		Description: "Standard snake placement and food spawning",
		Author:      "Battlesnake",
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
	minFood := int(settings.MinimumFood)
	foodSpawnChance := int(settings.FoodSpawnChance)
	numCurrentFood := len(lastBoardState.Food)

	if numCurrentFood < minFood {
		placeFoodRandomly(rand, lastBoardState, editor, minFood-numCurrentFood)
		return nil
	}
	if foodSpawnChance > 0 && (100-rand.Intn(100)) < foodSpawnChance {
		placeFoodRandomly(rand, lastBoardState, editor, 1)
		return nil
	}

	return nil
}

func placeFoodRandomly(rand rules.Rand, b *rules.BoardState, editor Editor, n int) {
	unoccupiedPoints := rules.GetUnoccupiedPoints(b, false)

	if len(unoccupiedPoints) < n {
		n = len(unoccupiedPoints)
	}

	rand.Shuffle(len(unoccupiedPoints), func(i int, j int) {
		unoccupiedPoints[i], unoccupiedPoints[j] = unoccupiedPoints[j], unoccupiedPoints[i]
	})

	for i := 0; i < n; i++ {
		editor.AddFood(unoccupiedPoints[i])
	}
}
