package maps

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/BattlesnakeOfficial/rules"
)

type LimitInfoMap struct{}

func init() {
	globalRegistry.RegisterMap("limitInfo", LimitInfoMap{})
}

func (m LimitInfoMap) ID() string {
	return "limitInfo"
}

func (m LimitInfoMap) Meta() Metadata {
	return Metadata{
		Name:        "Standard with per snake view range of 5 cells",
		Description: "Standard snake placement and food spawning but limited vision/information",
		Author:      "Kien Nguyen & Yannik Mahlau",
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  16,
		BoardSizes:  OddSizes(rules.BoardSizeSmall, rules.BoardSizeXXLarge),
		Tags:        []string{},
	}
}

func (m LimitInfoMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	rand := settings.GetRand(0)

	if len(initialBoardState.Snakes) > int(m.Meta().MaxPlayers) {
		return rules.ErrorTooManySnakes
	}

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

func (m LimitInfoMap) PreUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return nil
}

func (m LimitInfoMap) PostUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	rand := settings.GetRand(lastBoardState.Turn)

	foodNeeded := checkFoodNeedingPlacement(rand, settings, lastBoardState)
	if foodNeeded > 0 {
		placeFoodRandomlySaveSpawn(rand, lastBoardState, editor, foodNeeded)
	}

	for k, v := range editor.GameState() {
		if strings.HasPrefix(k, "food_spawn_") {
			spawnTurn, _ := strconv.Atoi(v)
			if spawnTurn < lastBoardState.Turn {
				delete(editor.GameState(), k)
			}
		}
	}

	return nil
}

func placeFoodRandomlySaveSpawn(rand rules.Rand, b *rules.BoardState, editor Editor, n int) {
	unoccupiedPoints := rules.GetUnoccupiedPoints(b, false, false)
	placeFoodRandomlyAtPositionsSaveSpawn(rand, b, editor, n, unoccupiedPoints)
}

func placeFoodRandomlyAtPositionsSaveSpawn(rand rules.Rand, b *rules.BoardState, editor Editor, n int, positions []rules.Point) {
	if len(positions) < n {
		n = len(positions)
	}

	rand.Shuffle(len(positions), func(i int, j int) {
		positions[i], positions[j] = positions[j], positions[i]
	})

	for i := 0; i < n; i++ {
		editor.AddFood(positions[i])
		key := fmt.Sprintf("food_spawn_%d_%d", positions[i].X, positions[i].Y)
		editor.GameState()[key] = fmt.Sprintf("%d", b.Turn+1)
	}
}
