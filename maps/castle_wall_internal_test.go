package maps

import (
	"fmt"
	"testing"

	"github.com/BattlesnakeOfficial/rules"

	"github.com/stretchr/testify/require"
)

func TestCastleWallMaps(t *testing.T) {

	tests := []struct {
		gameMap        GameMap
		startPositions [][]rules.Point
	}{
		{
			gameMap:        CastleWallMediumHazardsMap{},
			startPositions: castleWallMediumStartPositions,
		},
		{
			gameMap:        CastleWallLargeHazardsMap{},
			startPositions: castleWallLargeStartPositions,
		},
		{
			gameMap:        CastleWallExtraLargeHazardsMap{},
			startPositions: castleWallExtraLargeStartPositions,
		},
	}

	for _, test := range tests {
		t.Run(test.gameMap.ID(), func(t *testing.T) {
			m := test.gameMap
			settings := rules.Settings{}
			sizes := test.gameMap.Meta().BoardSizes
			for _, s := range sizes {
				initialState := rules.NewBoardState(int(s.Width), int(s.Height))
				startPositions := test.startPositions
				for i := 0; i < int(test.gameMap.Meta().MaxPlayers); i++ {
					initialState.Snakes = append(initialState.Snakes, rules.Snake{ID: fmt.Sprint(i), Body: []rules.Point{}})
				}
				nextState := rules.NewBoardState(int(s.Width), int(s.Height))
				editor := NewBoardStateEditor(nextState)
				err := m.SetupBoard(initialState, settings, editor)
				require.NoError(t, err)
				for _, s := range nextState.Snakes {
					require.Len(t, s.Body, rules.SnakeStartSize, "Placed snakes should have the right length")
					require.Equal(t, s.Health, rules.SnakeMaxHealth, "Placed snakes should have the right health")
					require.NotEmpty(t, s.ID, "Snake ID shouldn't be empty (should get copied when placed)")

					// Check that the snake is placed at one of the specified start positions
					validStart := false
					for _, q := range startPositions {
						for i := 0; i < len(q); i++ {
							if q[i].X == s.Body[0].X && q[i].Y == s.Body[0].Y {
								validStart = true
								break
							}
						}
					}
					require.True(t, validStart, "Snake must be placed in one of the specified start positions")
				}
			}
		})
	}
}
