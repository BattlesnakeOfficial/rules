package maps_test

import (
	"fmt"
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/BattlesnakeOfficial/rules/maps"
	"github.com/stretchr/testify/require"
)

func TestInnerBorderHazardsMap(t *testing.T) {

	tests := []struct {
		boardSize       int
		expectedHazards int
	}{
		{11, 32},
		{19, 64},
		{25, 88},
	}

	for _, tc := range tests {

		t.Run(fmt.Sprintf("%dx%d", tc.boardSize, tc.boardSize), func(t *testing.T) {
			m := maps.InnerBorderHazardsMap{}
			state := rules.NewBoardState(tc.boardSize, tc.boardSize)
			settings := rules.Settings{}

			// ensure the ring of hazards is added to the board at setup
			editor := maps.NewBoardStateEditor(state)
			require.Empty(t, state.Hazards)
			err := m.SetupBoard(state, settings, editor)
			require.NoError(t, err)
			require.NotEmpty(t, state.Hazards)
			require.Len(t, state.Hazards, tc.expectedHazards)
		})
	}
}
