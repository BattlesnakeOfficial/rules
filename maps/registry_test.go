package maps

import (
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/stretchr/testify/require"
)

const maxBoardWidth, maxBoardHeight = 25, 25

var testSettings rules.Settings = rules.Settings{
	FoodSpawnChance:     25,
	MinimumFood:         1,
	HazardDamagePerTurn: 14,
	RoyaleSettings: rules.RoyaleSettings{
		ShrinkEveryNTurns: 1,
	},
}

func TestRegisteredMaps(t *testing.T) {
	for mapName, gameMap := range globalRegistry {
		t.Run(mapName, func(t *testing.T) {
			require.Equalf(t, mapName, gameMap.ID(), "%#v game map doesn't return its own ID", mapName)

			var setupBoardState *rules.BoardState

			for width := 0; width < maxBoardWidth; width++ {
				for height := 0; height < maxBoardHeight; height++ {
					initialBoardState := rules.NewBoardState(width, height)
					initialBoardState.Snakes = append(initialBoardState.Snakes, rules.Snake{ID: "1", Body: []rules.Point{}})
					initialBoardState.Snakes = append(initialBoardState.Snakes, rules.Snake{ID: "2", Body: []rules.Point{}})
					passedBoardState := initialBoardState.Clone()
					tempBoardState := initialBoardState.Clone()
					err := gameMap.SetupBoard(passedBoardState, testSettings, NewBoardStateEditor(tempBoardState))
					if err == nil {
						setupBoardState = tempBoardState
						require.Equal(t, initialBoardState, passedBoardState, "BoardState should not be modified directly by GameMap.SetupBoard")
						break
					}
				}
			}
			require.NotNil(t, setupBoardState, "Map does not successfully setup the board at any supported combination of width and height")
			require.NotNil(t, setupBoardState.Food)
			require.NotNil(t, setupBoardState.Hazards)
			require.NotNil(t, setupBoardState.Snakes)
			for _, snake := range setupBoardState.Snakes {
				require.NotEmpty(t, snake.Body, "Map should place all snakes by initializing their body")
			}

			previousBoardState := rules.NewBoardState(rules.BoardSizeMedium, rules.BoardSizeMedium)
			previousBoardState.Food = append(previousBoardState.Food, []rules.Point{{X: 1, Y: 2}, {X: 3, Y: 4}}...)
			previousBoardState.Hazards = append(previousBoardState.Food, []rules.Point{{X: 4, Y: 3}, {X: 2, Y: 1}}...)
			previousBoardState.Snakes = append(previousBoardState.Snakes, rules.Snake{
				ID:     "1",
				Body:   []rules.Point{{X: 5, Y: 5}, {X: 5, Y: 4}, {X: 5, Y: 3}},
				Health: 100,
			})
			previousBoardState.Turn = 0

			passedBoardState := previousBoardState.Clone()
			tempBoardState := previousBoardState.Clone()
			err := gameMap.UpdateBoard(passedBoardState, testSettings, NewBoardStateEditor(tempBoardState))
			require.NoError(t, err, "GameMap.UpdateBoard returned an error")
			require.Equal(t, previousBoardState, passedBoardState, "BoardState should not be modified directly by GameMap.UpdateBoard")
		})
	}
}
