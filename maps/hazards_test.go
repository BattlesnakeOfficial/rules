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

func TestConcentricRingsHazardsMap(t *testing.T) {

	tests := []struct {
		boardSize       int
		expectedHazards int
	}{
		{11, 48},
	}

	for _, tc := range tests {

		t.Run(fmt.Sprintf("%dx%d", tc.boardSize, tc.boardSize), func(t *testing.T) {
			m := maps.ConcentricRingsHazardsMap{}
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

func TestColumnsHazardsMap(t *testing.T) {
	m := maps.ColumnsHazardsMap{}
	state := rules.NewBoardState(11, 11)
	settings := rules.Settings{}

	editor := maps.NewBoardStateEditor(state)
	require.Empty(t, state.Hazards)
	err := m.SetupBoard(state, settings, editor)
	require.NoError(t, err)
	require.NotEmpty(t, state.Hazards)
	require.Len(t, state.Hazards, 25)

	// a few spot checks
	require.Contains(t, state.Hazards, rules.Point{X: 1, Y: 1})
	require.Contains(t, state.Hazards, rules.Point{X: 1, Y: 5})
	require.Contains(t, state.Hazards, rules.Point{X: 9, Y: 1})
	require.Contains(t, state.Hazards, rules.Point{X: 9, Y: 9})
	require.NotContains(t, state.Hazards, rules.Point{X: 0, Y: 1})
	require.NotContains(t, state.Hazards, rules.Point{X: 8, Y: 4})
	require.NotContains(t, state.Hazards, rules.Point{X: 2, Y: 2})
	require.NotContains(t, state.Hazards, rules.Point{X: 4, Y: 9})
	require.NotContains(t, state.Hazards, rules.Point{X: 1, Y: 0})
}

func TestRiversAndBridgetsHazardsMap(t *testing.T) {
	// check error handling
	m := maps.RiverAndBridgesHazardsMap{}
	settings := rules.Settings{}

	// check error for unsupported board sizes
	state := rules.NewBoardState(9, 9)
	editor := maps.NewBoardStateEditor(state)
	err := m.SetupBoard(state, settings, editor)
	require.Error(t, err)

	// check all the supported sizes
	for _, size := range []int{11, 19, 25} {
		state = rules.NewBoardState(size, size)
		editor = maps.NewBoardStateEditor(state)
		require.Empty(t, state.Hazards)
		err = m.SetupBoard(state, settings, editor)
		require.NoError(t, err)
		require.NotEmpty(t, state.Hazards)
	}

}

func TestSpiralHazardsMap(t *testing.T) {
	// check error handling
	m := maps.SpiralHazardsMap{}
	settings := rules.Settings{}
	settings = settings.WithSeed(10)

	state := rules.NewBoardState(11, 11)
	editor := maps.NewBoardStateEditor(state)
	err := m.SetupBoard(state, settings, editor)
	require.NoError(t, err)

	for i := 1; i < 1000; i++ {
		err = m.UpdateBoard(state, settings, editor)
		require.NoError(t, err)
		state.Turn = i
	}
	require.NotEmpty(t, state.Hazards)
	require.Equal(t, 11*11, len(state.Hazards), "hazards should eventually fille the entire map")
}
