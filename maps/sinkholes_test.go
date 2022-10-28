package maps_test

import (
	"fmt"
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/BattlesnakeOfficial/rules/maps"
	"github.com/stretchr/testify/require"
)

func TestSinkholesMap(t *testing.T) {

	tests := []struct {
		boardSize               int
		expectedHazards         int
		expectedHazardsInCenter int
	}{
		{7, 27, 3},
		{11, 149, 5},
		{19, 431, 7},
	}

	for _, tc := range tests {

		t.Run(fmt.Sprintf("%dx%d", tc.boardSize, tc.boardSize), func(t *testing.T) {
			m := maps.SinkholesMap{}
			state := rules.NewBoardState(tc.boardSize, tc.boardSize)
			settings := rules.Settings{}

			// ensure the ring of hazards is added to the board at setup
			editor := maps.NewBoardStateEditor(state)
			require.Empty(t, state.Hazards)
			err := m.SetupBoard(state, settings, editor)
			require.NoError(t, err)
			require.Empty(t, state.Hazards)

			totalTurns := 100
			for i := 0; i < totalTurns; i++ {
				state.Turn = i
				err = m.PostUpdateBoard(state, settings, editor)
				require.NoError(t, err)
			}
			require.NotEmpty(t, state.Hazards)
			require.Len(t, state.Hazards, tc.expectedHazards)

			centerPoint := rules.Point{X: tc.boardSize / 2, Y: tc.boardSize / 2}
			numCenterHazards := 0
			for _, p := range state.Hazards {
				if p == centerPoint {
					numCenterHazards += 1
				}
			}
			require.Equal(t, numCenterHazards, tc.expectedHazardsInCenter)
		})
	}
}
