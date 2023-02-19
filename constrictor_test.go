package rules

import (
	"testing"
)

// Test that two equal snakes collide and both get eliminated
// also checks:
//   - food removed
//   - health back to max
var constrictorMoveAndCollideMAD = gameTestCase{
	"Constrictor Case Move and Collide",
	&BoardState{
		Turn:   41,
		Width:  10,
		Height: 10,
		Snakes: []Snake{
			{
				ID:     "one",
				Body:   []Point{{X: 1, Y: 1}, {X: 2, Y: 1}},
				Health: 99,
			},
			{
				ID:     "two",
				Body:   []Point{{X: 1, Y: 2}, {X: 2, Y: 2}},
				Health: 99,
			},
		},
		Food:    []Point{{X: 10, Y: 10}, {X: 9, Y: 9}, {X: 8, Y: 8}},
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
				Body:             []Point{{X: 1, Y: 2}, {X: 1, Y: 1}, {X: 1, Y: 1}},
				Health:           100,
				EliminatedCause:  EliminatedByCollision,
				EliminatedBy:     "two",
				EliminatedOnTurn: 42,
			},
			{
				ID:               "two",
				Body:             []Point{{X: 1, Y: 1}, {X: 1, Y: 2}, {X: 1, Y: 2}},
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
