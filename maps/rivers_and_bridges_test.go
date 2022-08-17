package maps_test

import (
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/BattlesnakeOfficial/rules/maps"
	"github.com/stretchr/testify/require"
)

func TestRiversAndBridgetsHazardsMap(t *testing.T) {
	// check error handling
	m := maps.RiverAndBridgesMediumHazardsMap{}
	settings := rules.Settings{}

	// check error for unsupported board sizes
	state := rules.NewBoardState(9, 9)
	editor := maps.NewBoardStateEditor(state)
	err := m.SetupBoard(state, settings, editor)
	require.Error(t, err)

	// check all the supported sizes
	for _, size := range []int{11, 19, 25} {
		state = rules.NewBoardState(size, size)
		state.Snakes = append(state.Snakes, rules.Snake{ID: "1", Body: []rules.Point{}})
		editor = maps.NewBoardStateEditor(state)
		require.Empty(t, state.Hazards)
		err = m.SetupBoard(state, settings, editor)
		require.NoError(t, err)
		require.NotEmpty(t, state.Hazards)
	}
}
