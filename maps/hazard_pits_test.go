package maps_test

import (
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/BattlesnakeOfficial/rules/maps"
	"github.com/stretchr/testify/require"
)

func TestHazardPitsMap(t *testing.T) {
	// check error handling
	m := maps.HazardPitsMap{}
	settings := rules.Settings{}
	// check error for unsupported board sizes
	state := rules.NewBoardState(9,
		9)
	editor := maps.NewBoardStateEditor(state)
	err := m.SetupBoard(state, settings, editor)
	require.Error(t, err)

	// too big
	state = rules.NewBoardState(19, 19)
	editor = maps.NewBoardStateEditor(state)
	err = m.SetupBoard(state, settings, editor)
	require.Error(t, err)

	// Too many snakes
	state = rules.NewBoardState(19, 19)
	editor = maps.NewBoardStateEditor(state)
	state.Snakes = append(state.Snakes, rules.Snake{ID: "1", Body: []rules.Point{}})
	state.Snakes = append(state.Snakes, rules.Snake{ID: "2", Body: []rules.Point{}})
	state.Snakes = append(state.Snakes, rules.Snake{ID: "3", Body: []rules.Point{}})
	state.Snakes = append(state.Snakes, rules.Snake{ID: "4", Body: []rules.Point{}})
	state.Snakes = append(state.Snakes, rules.Snake{ID: "5", Body: []rules.Point{}})
	err = m.SetupBoard(state, settings, editor)
	require.Error(t, err)

	state = rules.NewBoardState(int(11), int(11))
	m = maps.HazardPitsMap{}
	settings = rules.NewSettingsWithParams(rules.ParamShrinkEveryNTurns, "1")
	editor = maps.NewBoardStateEditor(state)
	require.Empty(t, state.Hazards)
	err = m.SetupBoard(state, settings, editor)
	require.NoError(t, err)
	require.Empty(t, state.Hazards)
	// Verify the hazard progression through the turns
	for i := 0; i < 16; i++ {
		state.Turn = i
		err = m.PostUpdateBoard(state, settings, editor)
		require.NoError(t, err)
		if i == 1 {
			require.Len(t, state.Hazards, 21)
		} else if i == 2 {
			require.Len(t, state.Hazards, 42)
		} else if i == 3 {
			require.Len(t, state.Hazards, 63)
		} else if i == 4 {
			require.Len(t, state.Hazards, 84)
		} else if i == 5 {
			require.Len(t, state.Hazards, 84)
		} else if i == 6 {
			require.Len(t, state.Hazards, 84)
		} else if i == 7 {
			require.Len(t, state.Hazards, 0)
		} else if i == 8 {
			require.Len(t, state.Hazards, 21)
		}
	}
}
