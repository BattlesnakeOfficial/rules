package maps_test

import (
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/BattlesnakeOfficial/rules/maps"
	"github.com/stretchr/testify/require"
)

func TestEmptyMapInterface(t *testing.T) {
	var _ maps.GameMap = maps.EmptyMap{}
}

func TestEmptyMapSetupBoard(t *testing.T) {
	m := maps.EmptyMap{}
	settings := rules.Settings{}

	tests := []struct {
		name              string
		initialBoardState *rules.BoardState
		rand              rules.Rand

		expected *rules.BoardState
		err      error
	}{
		{
			"empty 7x7",
			rules.NewBoardState(7, 7),
			rules.MinRand,
			rules.NewBoardState(7, 7),
			nil,
		},
		{
			"not enough room for snakes 7x7",
			rules.NewBoardState(7, 7).WithSnakes(generateSnakes(17)),
			rules.MinRand,
			nil,
			rules.ErrorTooManySnakes,
		},
		{
			"not enough room for snakes 5x5",
			rules.NewBoardState(5, 5).WithSnakes(generateSnakes(14)),
			rules.MinRand,
			nil,
			rules.ErrorTooManySnakes,
		},
		{
			"full 11x11 min",
			rules.NewBoardState(11, 11).WithSnakes(generateSnakes(8)),
			rules.MinRand,
			rules.NewBoardState(11, 11).WithSnakes([]rules.Snake{
				{ID: "1", Body: []rules.Point{{X: 1, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 1}}, Health: 100},
				{ID: "2", Body: []rules.Point{{X: 1, Y: 9}, {X: 1, Y: 9}, {X: 1, Y: 9}}, Health: 100},
				{ID: "3", Body: []rules.Point{{X: 9, Y: 1}, {X: 9, Y: 1}, {X: 9, Y: 1}}, Health: 100},
				{ID: "4", Body: []rules.Point{{X: 9, Y: 9}, {X: 9, Y: 9}, {X: 9, Y: 9}}, Health: 100},
				{ID: "5", Body: []rules.Point{{X: 1, Y: 5}, {X: 1, Y: 5}, {X: 1, Y: 5}}, Health: 100},
				{ID: "6", Body: []rules.Point{{X: 5, Y: 1}, {X: 5, Y: 1}, {X: 5, Y: 1}}, Health: 100},
				{ID: "7", Body: []rules.Point{{X: 5, Y: 9}, {X: 5, Y: 9}, {X: 5, Y: 9}}, Health: 100},
				{ID: "8", Body: []rules.Point{{X: 9, Y: 5}, {X: 9, Y: 5}, {X: 9, Y: 5}}, Health: 100},
			}),
			nil,
		},
		{
			"full 11x11 max",
			rules.NewBoardState(11, 11).WithSnakes(generateSnakes(8)),
			rules.MaxRand,
			rules.NewBoardState(11, 11).WithSnakes([]rules.Snake{
				{ID: "1", Body: []rules.Point{{X: 5, Y: 1}, {X: 5, Y: 1}, {X: 5, Y: 1}}, Health: 100},
				{ID: "2", Body: []rules.Point{{X: 5, Y: 9}, {X: 5, Y: 9}, {X: 5, Y: 9}}, Health: 100},
				{ID: "3", Body: []rules.Point{{X: 9, Y: 5}, {X: 9, Y: 5}, {X: 9, Y: 5}}, Health: 100},
				{ID: "4", Body: []rules.Point{{X: 1, Y: 5}, {X: 1, Y: 5}, {X: 1, Y: 5}}, Health: 100},
				{ID: "5", Body: []rules.Point{{X: 1, Y: 9}, {X: 1, Y: 9}, {X: 1, Y: 9}}, Health: 100},
				{ID: "6", Body: []rules.Point{{X: 9, Y: 1}, {X: 9, Y: 1}, {X: 9, Y: 1}}, Health: 100},
				{ID: "7", Body: []rules.Point{{X: 9, Y: 9}, {X: 9, Y: 9}, {X: 9, Y: 9}}, Health: 100},
				{ID: "8", Body: []rules.Point{{X: 1, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 1}}, Health: 100},
			}),
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nextBoardState := rules.NewBoardState(test.initialBoardState.Width, test.initialBoardState.Height)
			editor := maps.NewBoardStateEditor(nextBoardState)
			settings := settings.WithRand(test.rand)

			err := m.SetupBoard(test.initialBoardState, settings, editor)

			if test.err != nil {
				require.Equal(t, test.err, err)
			} else {
				require.Equal(t, test.expected, nextBoardState)
			}
		})
	}
}

func TestEmptyMapUpdateBoard(t *testing.T) {
	m := maps.EmptyMap{}
	initialBoardState := rules.NewBoardState(2, 2).WithFood([]rules.Point{{X: 0, Y: 0}})
	settings := rules.NewSettingsWithParams(rules.ParamFoodSpawnChance, "50", rules.ParamMinimumFood, "2").WithRand(rules.MaxRand)
	nextBoardState := initialBoardState.Clone()

	err := m.PostUpdateBoard(initialBoardState.Clone(), settings, maps.NewBoardStateEditor(nextBoardState))

	require.NoError(t, err)
	expectedBoardState := rules.NewBoardState(2, 2).WithFood([]rules.Point{{X: 0, Y: 0}})
	require.Equal(t, expectedBoardState, nextBoardState)
}
