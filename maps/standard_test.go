package maps

import (
	"fmt"
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/stretchr/testify/require"
)

func TestStandardMapInterface(t *testing.T) {
	var _ GameMap = StandardMap{}
}

func TestStandardMapSetupBoard(t *testing.T) {
	m := StandardMap{}
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
				Food:    []rules.Point{{X: 3, Y: 3}},
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
			rules.ErrorNoRoomForSnake,
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
					{ID: "2", Body: []rules.Point{{X: 1, Y: 5}, {X: 1, Y: 5}, {X: 1, Y: 5}}, Health: 100},
					{ID: "3", Body: []rules.Point{{X: 1, Y: 9}, {X: 1, Y: 9}, {X: 1, Y: 9}}, Health: 100},
					{ID: "4", Body: []rules.Point{{X: 5, Y: 1}, {X: 5, Y: 1}, {X: 5, Y: 1}}, Health: 100},
					{ID: "5", Body: []rules.Point{{X: 5, Y: 9}, {X: 5, Y: 9}, {X: 5, Y: 9}}, Health: 100},
					{ID: "6", Body: []rules.Point{{X: 9, Y: 1}, {X: 9, Y: 1}, {X: 9, Y: 1}}, Health: 100},
					{ID: "7", Body: []rules.Point{{X: 9, Y: 5}, {X: 9, Y: 5}, {X: 9, Y: 5}}, Health: 100},
					{ID: "8", Body: []rules.Point{{X: 9, Y: 9}, {X: 9, Y: 9}, {X: 9, Y: 9}}, Health: 100},
				},
				Food: []rules.Point{
					{X: 0, Y: 2},
					{X: 0, Y: 4},
					{X: 0, Y: 8},
					{X: 4, Y: 0},
					{X: 4, Y: 10},
					{X: 8, Y: 0},
					{X: 10, Y: 4},
					{X: 8, Y: 10},
					{X: 5, Y: 5},
				},
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
					{ID: "1", Body: []rules.Point{{X: 1, Y: 5}, {X: 1, Y: 5}, {X: 1, Y: 5}}, Health: 100},
					{ID: "2", Body: []rules.Point{{X: 1, Y: 9}, {X: 1, Y: 9}, {X: 1, Y: 9}}, Health: 100},
					{ID: "3", Body: []rules.Point{{X: 5, Y: 1}, {X: 5, Y: 1}, {X: 5, Y: 1}}, Health: 100},
					{ID: "4", Body: []rules.Point{{X: 5, Y: 9}, {X: 5, Y: 9}, {X: 5, Y: 9}}, Health: 100},
					{ID: "5", Body: []rules.Point{{X: 9, Y: 1}, {X: 9, Y: 1}, {X: 9, Y: 1}}, Health: 100},
					{ID: "6", Body: []rules.Point{{X: 9, Y: 5}, {X: 9, Y: 5}, {X: 9, Y: 5}}, Health: 100},
					{ID: "7", Body: []rules.Point{{X: 9, Y: 9}, {X: 9, Y: 9}, {X: 9, Y: 9}}, Health: 100},
					{ID: "8", Body: []rules.Point{{X: 1, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 1}}, Health: 100},
				},
				Food: []rules.Point{
					{X: 0, Y: 6},
					{X: 2, Y: 10},
					{X: 6, Y: 0},
					{X: 6, Y: 10},
					{X: 10, Y: 2},
					{X: 10, Y: 6},
					{X: 10, Y: 8},
					{X: 2, Y: 0},
					{X: 5, Y: 5},
				},
				Hazards: []rules.Point{},
			},
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nextBoardState := rules.NewBoardState(test.initialBoardState.Width, test.initialBoardState.Height)
			editor := NewBoardStateEditor(nextBoardState, test.rand)

			err := m.SetupBoard(*test.initialBoardState, settings, editor)

			if test.err != nil {
				require.Equal(t, test.err, err)
			} else {
				require.Equal(t, test.expected, nextBoardState)
			}
		})
	}
}

func TestStandardMapUpdateBoard(t *testing.T) {
	m := StandardMap{}

	tests := []struct {
		name              string
		initialBoardState *rules.BoardState
		settings          rules.Settings
		rand              rules.Rand

		expected *rules.BoardState
	}{
		{
			"empty no food",
			rules.NewBoardState(2, 2),
			rules.Settings{
				FoodSpawnChance: 0,
				MinimumFood:     0,
			},
			rules.MinRand,
			&rules.BoardState{
				Width:   2,
				Height:  2,
				Snakes:  []rules.Snake{},
				Food:    []rules.Point{},
				Hazards: []rules.Point{},
			},
		},
		{
			"empty MinimumFood",
			rules.NewBoardState(2, 2),
			rules.Settings{
				FoodSpawnChance: 0,
				MinimumFood:     2,
			},
			rules.MinRand,
			&rules.BoardState{
				Width:   2,
				Height:  2,
				Snakes:  []rules.Snake{},
				Food:    []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}},
				Hazards: []rules.Point{},
			},
		},
		{
			"not empty MinimumFood",
			&rules.BoardState{
				Width:   2,
				Height:  2,
				Snakes:  []rules.Snake{},
				Food:    []rules.Point{{X: 0, Y: 1}},
				Hazards: []rules.Point{},
			},
			rules.Settings{
				FoodSpawnChance: 0,
				MinimumFood:     2,
			},
			rules.MinRand,
			&rules.BoardState{
				Width:   2,
				Height:  2,
				Snakes:  []rules.Snake{},
				Food:    []rules.Point{{X: 0, Y: 1}, {X: 0, Y: 0}},
				Hazards: []rules.Point{},
			},
		},
		{
			"empty FoodSpawnChance inactive",
			rules.NewBoardState(2, 2),
			rules.Settings{
				FoodSpawnChance: 50,
				MinimumFood:     0,
			},
			rules.MinRand,
			&rules.BoardState{
				Width:   2,
				Height:  2,
				Snakes:  []rules.Snake{},
				Food:    []rules.Point{},
				Hazards: []rules.Point{},
			},
		},
		{
			"empty FoodSpawnChance active",
			rules.NewBoardState(2, 2),
			rules.Settings{
				FoodSpawnChance: 50,
				MinimumFood:     0,
			},
			rules.MaxRand,
			&rules.BoardState{
				Width:   2,
				Height:  2,
				Snakes:  []rules.Snake{},
				Food:    []rules.Point{{X: 0, Y: 1}},
				Hazards: []rules.Point{},
			},
		},
		{
			"not empty FoodSpawnChance active",
			&rules.BoardState{
				Width:   2,
				Height:  2,
				Snakes:  []rules.Snake{},
				Food:    []rules.Point{{X: 0, Y: 0}},
				Hazards: []rules.Point{},
			},
			rules.Settings{
				FoodSpawnChance: 50,
				MinimumFood:     0,
			},
			rules.MaxRand,
			&rules.BoardState{
				Width:   2,
				Height:  2,
				Snakes:  []rules.Snake{},
				Food:    []rules.Point{{X: 0, Y: 0}, {X: 1, Y: 0}},
				Hazards: []rules.Point{},
			},
		},
		{
			"not empty FoodSpawnChance no room",
			&rules.BoardState{
				Width:   2,
				Height:  2,
				Snakes:  []rules.Snake{},
				Food:    []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 0}, {X: 1, Y: 1}},
				Hazards: []rules.Point{},
			},
			rules.Settings{
				FoodSpawnChance: 50,
				MinimumFood:     0,
			},
			rules.MaxRand,
			&rules.BoardState{
				Width:   2,
				Height:  2,
				Snakes:  []rules.Snake{},
				Food:    []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 0}, {X: 1, Y: 1}},
				Hazards: []rules.Point{},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			nextBoardState := test.initialBoardState.Clone()
			editor := NewBoardStateEditor(nextBoardState, test.rand)

			err := m.UpdateBoard(*test.initialBoardState.Clone(), test.settings, editor)

			require.NoError(t, err)
			require.Equal(t, test.expected, nextBoardState)
		})
	}
}

func generateSnakes(n int) []rules.Snake {
	var snakes []rules.Snake
	for i := 0; i < n; i++ {
		snakes = append(snakes, rules.Snake{
			ID: fmt.Sprint(i + 1),
		})
	}
	return snakes
}
