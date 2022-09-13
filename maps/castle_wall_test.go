package maps_test

import (
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/BattlesnakeOfficial/rules/maps"
	"github.com/stretchr/testify/require"
)

func TestCastleWallHazardsMap(t *testing.T) {
	// check error handling
	m := maps.CastleWallMediumHazardsMap{}
	settings := rules.Settings{}

	// check error for unsupported board sizes
	state := rules.NewBoardState(7, 7)
	editor := maps.NewBoardStateEditor(state)
	err := m.SetupBoard(state, settings, editor)
	require.Error(t, err)

	tests := []struct {
		Map    maps.GameMap
		Width  int
		Height int
	}{
		{maps.CastleWallMediumHazardsMap{}, 11, 11},
		{maps.CastleWallLargeHazardsMap{}, 19, 19},
		{maps.CastleWallExtraLargeHazardsMap{}, 25, 25},
	}

	// check all the supported sizes
	for _, test := range tests {
		state = rules.NewBoardState(test.Width, test.Height)
		state.Snakes = append(state.Snakes, rules.Snake{ID: "1", Body: []rules.Point{}})
		editor = maps.NewBoardStateEditor(state)
		require.Empty(t, state.Hazards)
		err = test.Map.SetupBoard(state, settings, editor)
		require.NoError(t, err)
		require.NotEmpty(t, state.Hazards)
	}
}
