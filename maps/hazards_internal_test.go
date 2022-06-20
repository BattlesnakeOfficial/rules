package maps

import (
	"fmt"
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/stretchr/testify/require"
)

func TestRiversAndBridgesSnakePlacement(t *testing.T) {
	m := RiverAndBridgesHazardsMap{}
	settings := rules.Settings{}

	// check all the supported sizes
	for _, size := range []int{11, 19, 25} {
		initialState := rules.NewBoardState(size, size)
		startPositions := riversAndBridgesStartPositions[rules.Point{X: size, Y: size}]
		maxSnakes := len(startPositions)
		for i := 0; i < maxSnakes; i++ {
			initialState.Snakes = append(initialState.Snakes, rules.Snake{ID: fmt.Sprint(i), Body: []rules.Point{}})
		}

		nextState := rules.NewBoardState(size, size)
		editor := NewBoardStateEditor(nextState)
		err := m.SetupBoard(initialState, settings, editor)
		require.NoError(t, err)
		for _, s := range nextState.Snakes {
			require.Len(t, s.Body, rules.SnakeStartSize, "Placed snakes should have the right length")
			require.Equal(t, s.Health, rules.SnakeMaxHealth, "Placed snakes should have the right health")
			require.NotEmpty(t, s.ID, "Snake ID shouldn't be empty (should get copied when placed)")

			// Check that the snake is placed at one of the specified start positions
			validStart := false
			for _, p := range startPositions {
				if p.X == s.Body[0].X && p.Y == s.Body[0].Y {
					validStart = true
					break
				}
			}
			require.True(t, validStart, "Snake must be placed in one of the specified start positions")
		}
	}
}
