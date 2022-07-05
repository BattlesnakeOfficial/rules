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
			&rules.BoardState{
				Width:   7,
				Height:  7,
				Snakes:  []rules.Snake{},
				Food:    []rules.Point{},
				Hazards: []rules.Point{},
			},
			nil,
		},
		{
			"not enough room for snakes 7x7",
			&rules.BoardState{
				Width:   7,
				Height:  7,
				Snakes:  generateSnakes(9),
				Food:    []rules.Point{},
				Hazards: []rules.Point{},
			},
			rules.MinRand,
			nil,
			rules.ErrorTooManySnakes,
		},
		{
			"not enough room for snakes 5x5",
			&rules.BoardState{
				Width:   5,
				Height:  5,
				Snakes:  generateSnakes(14),
				Food:    []rules.Point{},
				Hazards: []rules.Point{},
			},
			rules.MinRand,
			nil,
			rules.ErrorTooManySnakes,
		},
		{
			"full 11x11 min",
			&rules.BoardState{
				Width:   11,
				Height:  11,
				Snakes:  generateSnakes(8),
				Food:    []rules.Point{},
				Hazards: []rules.Point{},
			},
			rules.MinRand,
			&rules.BoardState{
				Width:  11,
				Height: 11,
				Snakes: []rules.Snake{
					{ID: "1", Body: []rules.Point{{X: 1, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 1}}, Health: 100},
					{ID: "2", Body: []rules.Point{{X: 1, Y: 9}, {X: 1, Y: 9}, {X: 1, Y: 9}}, Health: 100},
					{ID: "3", Body: []rules.Point{{X: 9, Y: 1}, {X: 9, Y: 1}, {X: 9, Y: 1}}, Health: 100},
					{ID: "4", Body: []rules.Point{{X: 9, Y: 9}, {X: 9, Y: 9}, {X: 9, Y: 9}}, Health: 100},
					{ID: "5", Body: []rules.Point{{X: 1, Y: 5}, {X: 1, Y: 5}, {X: 1, Y: 5}}, Health: 100},
					{ID: "6", Body: []rules.Point{{X: 5, Y: 1}, {X: 5, Y: 1}, {X: 5, Y: 1}}, Health: 100},
					{ID: "7", Body: []rules.Point{{X: 5, Y: 9}, {X: 5, Y: 9}, {X: 5, Y: 9}}, Health: 100},
					{ID: "8", Body: []rules.Point{{X: 9, Y: 5}, {X: 9, Y: 5}, {X: 9, Y: 5}}, Health: 100},
				},
				Food:    []rules.Point{},
				Hazards: []rules.Point{},
			},
			nil,
		},
		{
			"full 11x11 max",
			&rules.BoardState{
				Width:   11,
				Height:  11,
				Snakes:  generateSnakes(8),
				Food:    []rules.Point{},
				Hazards: []rules.Point{},
			},
			rules.MaxRand,
			&rules.BoardState{
				Width:  11,
				Height: 11,
				Snakes: []rules.Snake{
					{ID: "1", Body: []rules.Point{{X: 5, Y: 1}, {X: 5, Y: 1}, {X: 5, Y: 1}}, Health: 100},
					{ID: "2", Body: []rules.Point{{X: 5, Y: 9}, {X: 5, Y: 9}, {X: 5, Y: 9}}, Health: 100},
					{ID: "3", Body: []rules.Point{{X: 9, Y: 5}, {X: 9, Y: 5}, {X: 9, Y: 5}}, Health: 100},
					{ID: "4", Body: []rules.Point{{X: 1, Y: 5}, {X: 1, Y: 5}, {X: 1, Y: 5}}, Health: 100},
					{ID: "5", Body: []rules.Point{{X: 1, Y: 9}, {X: 1, Y: 9}, {X: 1, Y: 9}}, Health: 100},
					{ID: "6", Body: []rules.Point{{X: 9, Y: 1}, {X: 9, Y: 1}, {X: 9, Y: 1}}, Health: 100},
					{ID: "7", Body: []rules.Point{{X: 9, Y: 9}, {X: 9, Y: 9}, {X: 9, Y: 9}}, Health: 100},
					{ID: "8", Body: []rules.Point{{X: 1, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 1}}, Health: 100},
				},
				Food:    []rules.Point{},
				Hazards: []rules.Point{},
			},
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
	initialBoardState := &rules.BoardState{
		Width:   2,
		Height:  2,
		Snakes:  []rules.Snake{},
		Food:    []rules.Point{{X: 0, Y: 0}},
		Hazards: []rules.Point{},
	}
	settings := rules.Settings{
		FoodSpawnChance: 50,
		MinimumFood:     2,
	}.WithRand(rules.MaxRand)
	nextBoardState := initialBoardState.Clone()

	err := m.UpdateBoard(initialBoardState.Clone(), settings, maps.NewBoardStateEditor(nextBoardState))

	require.NoError(t, err)
	require.Equal(t, &rules.BoardState{
		Width:   2,
		Height:  2,
		Snakes:  []rules.Snake{},
		Food:    []rules.Point{{X: 0, Y: 0}},
		Hazards: []rules.Point{},
	}, nextBoardState)
}
