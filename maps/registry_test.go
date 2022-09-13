package maps

import (
	"fmt"
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
			meta := gameMap.Meta()
			require.True(t, meta.Version > 0, fmt.Sprintf("registered maps must have a valid version (>= 1) - '%d' is invalid", meta.Version))
			require.NotZero(t, meta.MaxPlayers, "registered maps must have maximum players declared")
			require.LessOrEqual(t, meta.MaxPlayers, meta.MaxPlayers, "max players should always be >= min players")
			require.NotEmpty(t, meta.BoardSizes, "registered maps must have at least one supported size declared")
			require.NotNil(t, meta.Tags)
			var setupBoardState *rules.BoardState

			// "fuzz test" supported players
			mapSize := pickSize(meta)
			for i := meta.MinPlayers; i < meta.MaxPlayers; i++ {
				t.Run(fmt.Sprintf("%d players", i), func(t *testing.T) {
					initialBoardState := rules.NewBoardState(int(mapSize.Width), int(mapSize.Height))
					for j := 0; j < i; j++ {
						initialBoardState.Snakes = append(initialBoardState.Snakes, rules.Snake{ID: fmt.Sprint(j), Body: []rules.Point{}})
					}
					err := gameMap.SetupBoard(initialBoardState, testSettings, NewBoardStateEditor(initialBoardState))
					require.NoError(t, err, fmt.Sprintf("%d players should be supported by this map", i))
				})
			}

			// "fuzz test" supported map sizes
			if !meta.BoardSizes.IsUnlimited() {
				for _, mapSize := range meta.BoardSizes {
					t.Run(fmt.Sprintf("%dx%d map size", mapSize.Width, mapSize.Height), func(t *testing.T) {
						initialBoardState := rules.NewBoardState(int(mapSize.Width), int(mapSize.Height))
						for i := 0; i < meta.MaxPlayers; i++ {
							initialBoardState.Snakes = append(initialBoardState.Snakes, rules.Snake{ID: fmt.Sprint(i), Body: []rules.Point{}})
						}
						err := gameMap.SetupBoard(initialBoardState, testSettings, NewBoardStateEditor(initialBoardState))
						require.NoError(t, err, "error setting up map")
					})
				}
			}

			// Check that at least one map size can be setup without error
			for width := 0; width <= maxBoardWidth; width++ {
				for height := 0; height <= maxBoardHeight; height++ {
					initialBoardState := rules.NewBoardState(width, height)
					initialBoardState.Snakes = append(initialBoardState.Snakes, rules.Snake{ID: "1", Body: []rules.Point{}})
					if meta.MaxPlayers > 1 {
						initialBoardState.Snakes = append(initialBoardState.Snakes, rules.Snake{ID: "2", Body: []rules.Point{}})
					}
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

func pickSize(meta Metadata) Dimensions {
	// For unlimited, we can pick any size
	if meta.BoardSizes.IsUnlimited() {
		return Dimensions{Width: 11, Height: 11}
	}

	// For fixed, just pick the first supported size
	return meta.BoardSizes[0]
}

func TestListRegisteredMaps(t *testing.T) {
	keys := globalRegistry.List()
	mapCount := 0
	for k := range globalRegistry {
		// every registry key should exist in List results
		require.Contains(t, keys, k)
		mapCount++
	}
	// List should equal number of maps in the global registry
	require.Equal(t, len(keys), mapCount)
}
