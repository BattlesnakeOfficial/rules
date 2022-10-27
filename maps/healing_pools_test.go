package maps_test

import (
	"fmt"
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/BattlesnakeOfficial/rules/maps"
	"github.com/stretchr/testify/require"
)

func TestHealingPoolsMap(t *testing.T) {

	tests := []struct {
		boardSize        int
		expectedHazards  int
		allowableHazards []rules.Point
	}{
		{7, 1, []rules.Point{{X: 3, Y: 3}}},
		{11, 2, []rules.Point{{X: 3, Y: 3},
			{X: 7, Y: 7},
			{X: 3, Y: 7},
			{X: 7, Y: 3},
			{X: 3, Y: 5},
			{X: 7, Y: 5},
			{X: 5, Y: 7},
			{X: 5, Y: 3}}},
		{19, 4, []rules.Point{{X: 5, Y: 5},
			{X: 13, Y: 13},
			{X: 5, Y: 13},
			{X: 13, Y: 5},
			{X: 5, Y: 10},
			{X: 13, Y: 10},
			{X: 10, Y: 13},
			{X: 10, Y: 5}}},
	}

	for _, tc := range tests {

		t.Run(fmt.Sprintf("%dx%d", tc.boardSize, tc.boardSize), func(t *testing.T) {
			m := maps.HealingPoolsMap{}
			state := rules.NewBoardState(tc.boardSize, tc.boardSize)
			shrinkEveryNTurns := 10
			settings := rules.NewSettingsWithParams(rules.ParamShrinkEveryNTurns, fmt.Sprint(shrinkEveryNTurns))

			// ensure the hazards are added to the board at setup
			editor := maps.NewBoardStateEditor(state)
			require.Empty(t, state.Hazards)
			err := m.SetupBoard(state, settings, editor)
			require.NoError(t, err)
			require.NotEmpty(t, state.Hazards)
			require.Len(t, state.Hazards, tc.expectedHazards)

			for _, p := range state.Hazards {
				require.Contains(t, tc.allowableHazards, p)
			}

			// ensure the hazards are removed
			totalTurns := shrinkEveryNTurns*tc.expectedHazards + 1
			for i := 0; i < totalTurns; i++ {
				state.Turn = i
				err = m.PostUpdateBoard(state, settings, editor)
				require.NoError(t, err)
			}

			require.Equal(t, 0, len(state.Hazards))
		})
	}
}
