package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConstrictorRulesetInterface(t *testing.T) {
	var _ Ruleset = (*ConstrictorRuleset)(nil)
}

func TestConstrictorModifyInitialBoardState(t *testing.T) {
	tests := []struct {
		Height int32
		Width  int32
		IDs    []string
	}{
		{1, 1, []string{}},
		{1, 1, []string{"one"}},
		{2, 2, []string{"one"}},
		{2, 2, []string{"one", "two"}},
		{11, 1, []string{"one", "two"}},
		{11, 11, []string{}},
		{11, 11, []string{"one", "two", "three", "four", "five"}},
	}

	r := ConstrictorRuleset{}
	for testNum, test := range tests {
		state, err := CreateDefaultBoardState(test.Width, test.Height, test.IDs)
		require.NoError(t, err)
		require.NotNil(t, state)
		state, err = r.ModifyInitialBoardState(state)
		require.NoError(t, err)
		require.NotNil(t, state)
		require.Equal(t, test.Width, state.Width)
		require.Equal(t, test.Height, state.Height)
		require.Len(t, state.Food, 0, testNum)
		// Verify snakes
		require.Equal(t, len(test.IDs), len(state.Snakes))
		for i, id := range test.IDs {
			require.Equal(t, id, state.Snakes[i].ID)
			require.Equal(t, state.Snakes[i].Body[2], state.Snakes[i].Body[1])
		}
	}
}

func TestConstrictorCreateNextBoardState(t *testing.T) {
	tests := []struct {
		prevState     *BoardState
		moves         []SnakeMove
		expectedState *BoardState
	}{
		{
			&BoardState{
				Width:  3,
				Height: 3,
				Snakes: []Snake{
					{
						ID:     "one",
						Body:   []Point{{0, 0}, {0, 0}, {0, 0}},
						Health: 100,
					},
					{
						ID:     "two",
						Body:   []Point{{2, 2}, {2, 2}, {2, 2}},
						Health: 100,
					},
				},
				Food: []Point{},
			},
			[]SnakeMove{
				{ID: "one", Move: MoveUp},
				{ID: "two", Move: MoveDown},
			},
			&BoardState{
				Width:  3,
				Height: 3,
				Snakes: []Snake{
					{
						ID:     "one",
						Body:   []Point{{0, 1}, {0, 0}, {0, 0}},
						Health: 100,
					},
					{
						ID:     "two",
						Body:   []Point{{2, 1}, {2, 2}, {2, 2}},
						Health: 100,
					},
				},
				Food: []Point{},
			},
		},
		// Ensure snakes keep growing and are fed
		{
			&BoardState{
				Width:  3,
				Height: 3,
				Snakes: []Snake{
					{
						ID:     "one",
						Body:   []Point{{2, 0}, {1, 0}, {0, 0}, {0, 0}},
						Health: 75,
					},
				},
				Food: []Point{},
			},
			[]SnakeMove{
				{ID: "one", Move: MoveUp},
			},
			&BoardState{
				Width:  3,
				Height: 3,
				Snakes: []Snake{
					{
						ID:     "one",
						Body:   []Point{{2, 1}, {2, 0}, {1, 0}, {0, 0}, {0, 0}},
						Health: 100,
					},
				},
				Food: []Point{},
			},
		},
	}

	r := ConstrictorRuleset{}
	for _, test := range tests {
		nextState, err := r.CreateNextBoardState(test.prevState, test.moves)
		require.NoError(t, err)
		require.Equal(t, test.expectedState.Food, nextState.Food)
		require.Equal(t, test.expectedState.Snakes, nextState.Snakes)
	}
}
