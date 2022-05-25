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
		Height int
		Width  int
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
		state, err := CreateDefaultBoardState(MaxRand, test.Width, test.Height, test.IDs)
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

// Test that two equal snakes collide and both get eliminated
// also checks:
//	- food removed
//  - health back to max
var constrictorMoveAndCollideMAD = gameTestCase{
	"Constrictor Case Move and Collide",
	&BoardState{
		Width:  10,
		Height: 10,
		Snakes: []Snake{
			{
				ID:     "one",
				Body:   []Point{{1, 1}, {2, 1}},
				Health: 99,
			},
			{
				ID:     "two",
				Body:   []Point{{1, 2}, {2, 2}},
				Health: 99,
			},
		},
		Food:    []Point{{10, 10}, {9, 9}, {8, 8}},
		Hazards: []Point{},
	},
	[]SnakeMove{
		{ID: "one", Move: MoveUp},
		{ID: "two", Move: MoveDown},
	},
	nil,
	&BoardState{
		Width:  10,
		Height: 10,
		Snakes: []Snake{
			{
				ID:              "one",
				Body:            []Point{{1, 2}, {1, 1}, {1, 1}},
				Health:          100,
				EliminatedCause: EliminatedByCollision,
				EliminatedBy:    "two",
			},
			{
				ID:              "two",
				Body:            []Point{{1, 1}, {1, 2}, {1, 2}},
				Health:          100,
				EliminatedCause: EliminatedByCollision,
				EliminatedBy:    "one",
			},
		},
		Food:    []Point{},
		Hazards: []Point{},
	},
}

func TestConstrictorCreateNextBoardState(t *testing.T) {
	cases := []gameTestCase{
		standardCaseErrNoMoveFound,
		standardCaseErrZeroLengthSnake,
		constrictorMoveAndCollideMAD,
	}
	rb := NewRulesetBuilder().WithParams(map[string]string{
		ParamGameType: GameTypeConstrictor,
	})
	r := ConstrictorRuleset{}
	for _, gc := range cases {
		gc.requireValidNextState(t, &r)
		// also test a RulesBuilder constructed instance
		gc.requireValidNextState(t, rb.Ruleset())
		// also test a pipeline with the same settings
		gc.requireValidNextState(t, rb.PipelineRuleset(GameTypeConstrictor, NewPipeline(constrictorRulesetStages...)))
	}
}
