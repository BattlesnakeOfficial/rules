package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFeastRulesetInterface(t *testing.T) {
	var _ Ruleset = (*FeastRuleset)(nil)
}

func TestFeastCreateInitialBoardState(t *testing.T) {
	tests := []struct {
		Height          int32
		Width           int32
		IDs             []string
		ExpectedNumFood int
		Err             error
	}{
		{1, 1, []string{}, 1, nil},
		{1, 1, []string{"one"}, 0, nil},
		{2, 2, []string{"one"}, 3, nil},
		{2, 2, []string{"one", "two"}, 2, nil},
		{11, 1, []string{"one", "two"}, 9, nil},
		{11, 11, []string{}, 121, nil},
		{11, 11, []string{"one", "two", "three", "four", "five"}, 116, nil},
	}

	r := FeastRuleset{}
	for testNum, test := range tests {
		state, err := r.CreateInitialBoardState(test.Width, test.Height, test.IDs)
		require.Equal(t, test.Err, err)
		if err != nil {
			require.Nil(t, state)
			continue
		}
		require.NotNil(t, state)
		require.Equal(t, test.Width, state.Width)
		require.Equal(t, test.Height, state.Height)
		require.Equal(t, len(test.IDs), len(state.Snakes))
		for i, id := range test.IDs {
			require.Equal(t, id, state.Snakes[i].ID)
		}
		require.Len(t, state.Food, test.ExpectedNumFood, testNum)
	}
}

func TestFeastCreateNextBoardState(t *testing.T) {
	tests := []struct {
		prevState     *BoardState
		moves         []SnakeMove
		expectedError error
		expectedState *BoardState
	}{
		{
			&BoardState{
				Width:  3,
				Height: 3,
				Snakes: []Snake{
					{
						ID:     "one",
						Body:   []Point{{0, 0}, {1, 0}, {2, 0}},
						Health: 100,
					},
					{
						ID:     "two",
						Body:   []Point{{2, 2}, {1, 2}, {0, 2}},
						Health: 100,
					},
				},
				Food: []Point{},
			},
			[]SnakeMove{
				{ID: "one", Move: MoveDown},
				{ID: "two", Move: MoveUp},
			},
			nil,
			&BoardState{
				Width:  3,
				Height: 3,
				Snakes: []Snake{
					{
						ID:     "one",
						Body:   []Point{{0, 1}, {0, 0}, {1, 0}},
						Health: 99,
					},
					{
						ID:     "two",
						Body:   []Point{{2, 1}, {2, 2}, {1, 2}},
						Health: 99,
					},
				},
				Food: []Point{{0, 2}, {1, 1}, {2, 0}},
			},
		},
	}

	r := FeastRuleset{}
	for _, test := range tests {
		nextState, err := r.CreateNextBoardState(test.prevState, test.moves)
		require.Equal(t, test.expectedError, err)
		require.Equal(t, test.expectedState, nextState)
	}
}
