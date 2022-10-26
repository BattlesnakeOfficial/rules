package rules

import (
	"testing"
)

// Test that two equal snakes collide and both get eliminated
// also checks:
//	- food removed
//  - health back to max
var constrictorMoveAndCollideMAD = gameTestCase{
	"Constrictor Case Move and Collide",
	&BoardState{
		Turn:   41,
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
				ID:               "one",
				Body:             []Point{{1, 2}, {1, 1}, {1, 1}},
				Health:           100,
				EliminatedCause:  EliminatedByCollision,
				EliminatedBy:     "two",
				EliminatedOnTurn: 42,
			},
			{
				ID:               "two",
				Body:             []Point{{1, 1}, {1, 2}, {1, 2}},
				Health:           100,
				EliminatedCause:  EliminatedByCollision,
				EliminatedBy:     "one",
				EliminatedOnTurn: 42,
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
	r := NewRulesetBuilder().NamedRuleset(GameTypeConstrictor)
	for _, gc := range cases {
		// test a RulesBuilder constructed instance
		gc.requireValidNextState(t, r)
		// also test a pipeline with the same settings
		gc.requireValidNextState(t, NewRulesetBuilder().PipelineRuleset(GameTypeConstrictor, NewPipeline(constrictorRulesetStages...)))
	}
}
