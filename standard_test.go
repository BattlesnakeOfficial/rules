package rules

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func getStandardRuleset(settings Settings) Ruleset {
	return NewRulesetBuilder().WithSettings(settings).NamedRuleset(GameTypeStandard)
}

func TestSanity(t *testing.T) {
	r := getStandardRuleset(Settings{})

	state, err := CreateDefaultBoardState(MaxRand, 0, 0, []string{})
	require.NoError(t, err)
	require.NotNil(t, state)

	gameOver, state, err := r.Execute(state, []SnakeMove{})
	require.NoError(t, err)
	require.True(t, gameOver)
	require.NotNil(t, state)
	require.Equal(t, 0, state.Width)
	require.Equal(t, 0, state.Height)
	require.Len(t, state.Food, 0)
	require.Len(t, state.Snakes, 0)
}

func TestStandardName(t *testing.T) {
	r := getStandardRuleset(Settings{})
	require.Equal(t, "standard", r.Name())
}

// Checks that the error for a snake missing a move is returned
var standardCaseErrNoMoveFound = gameTestCase{
	"Standard Case Error No Move Found",
	&BoardState{
		Width:  10,
		Height: 10,
		Snakes: []Snake{
			{
				ID:     "one",
				Body:   []Point{{X: 1, Y: 1}, {X: 1, Y: 2}},
				Health: 100,
			},
			{
				ID:     "two",
				Body:   []Point{{X: 3, Y: 4}, {X: 3, Y: 3}},
				Health: 100,
			},
		},
		Food:    []Point{{X: 0, Y: 0}, {X: 1, Y: 0}},
		Hazards: []Point{},
	},
	[]SnakeMove{
		{ID: "one", Move: MoveUp},
	},
	ErrorNoMoveFound,
	nil,
}

// Checks that the error for a snake with no points is returned
var standardCaseErrZeroLengthSnake = gameTestCase{
	"Standard Case Error Zero Length Snake",
	&BoardState{
		Width:  10,
		Height: 10,
		Snakes: []Snake{
			{
				ID:     "one",
				Body:   []Point{{X: 1, Y: 1}, {X: 1, Y: 2}},
				Health: 100,
			},
			{
				ID:     "two",
				Body:   []Point{},
				Health: 100,
			},
		},
		Food:    []Point{{X: 0, Y: 0}, {X: 1, Y: 0}},
		Hazards: []Point{},
	},
	[]SnakeMove{
		{ID: "one", Move: MoveUp},
		{ID: "two", Move: MoveDown},
	},
	ErrorZeroLengthSnake,
	nil,
}

// Checks a basic state where a snake moves, eats and grows
var standardCaseMoveEatAndGrow = gameTestCase{
	"Standard Case Move Eat and Grow",
	&BoardState{
		Width:  10,
		Height: 10,
		Snakes: []Snake{
			{
				ID:     "one",
				Body:   []Point{{X: 1, Y: 1}, {X: 1, Y: 2}},
				Health: 100,
			},
			{
				ID:     "two",
				Body:   []Point{{X: 3, Y: 4}, {X: 3, Y: 3}},
				Health: 100,
			},
			{
				ID:              "three",
				Body:            []Point{},
				Health:          100,
				EliminatedCause: EliminatedByOutOfBounds,
			},
		},
		Food:    []Point{{X: 0, Y: 0}, {X: 1, Y: 0}},
		Hazards: []Point{},
	},
	[]SnakeMove{
		{ID: "one", Move: MoveDown},
		{ID: "two", Move: MoveUp},
		{ID: "three", Move: MoveLeft}, // Should be ignored
	},
	nil,
	&BoardState{
		Width:  10,
		Height: 10,
		Snakes: []Snake{
			{
				ID:     "one",
				Body:   []Point{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 1}},
				Health: 100,
			},
			{
				ID:     "two",
				Body:   []Point{{X: 3, Y: 5}, {X: 3, Y: 4}},
				Health: 99,
			},
			{
				ID:              "three",
				Body:            []Point{},
				Health:          100,
				EliminatedCause: EliminatedByOutOfBounds,
			},
		},
		Food:    []Point{{X: 0, Y: 0}},
		Hazards: []Point{},
	},
}

// Checks a basic state where two snakes of equal sizes collide, and both should
// be eliminated as a result.
var standardMoveAndCollideMAD = gameTestCase{
	"Standard Case Move and Collide",
	&BoardState{
		Turn:   0,
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
		Food:    []Point{},
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
				Body:             []Point{{X: 1, Y: 2}, {X: 1, Y: 1}},
				Health:           98,
				EliminatedCause:  EliminatedByCollision,
				EliminatedBy:     "two",
				EliminatedOnTurn: 1,
			},
			{
				ID:               "two",
				Body:             []Point{{X: 1, Y: 1}, {X: 1, Y: 2}},
				Health:           98,
				EliminatedCause:  EliminatedByCollision,
				EliminatedBy:     "one",
				EliminatedOnTurn: 1,
			},
		},
		Food:    []Point{},
		Hazards: []Point{},
	},
}

func TestStandardCreateNextBoardState(t *testing.T) {
	cases := []gameTestCase{
		standardCaseErrNoMoveFound,
		standardCaseErrZeroLengthSnake,
		standardCaseMoveEatAndGrow,
		standardMoveAndCollideMAD,
	}
	r := getStandardRuleset(Settings{})
	for _, gc := range cases {
		// test a RulesBuilder constructed instance
		gc.requireValidNextState(t, r)
		// also test a pipeline with the same settings
		gc.requireValidNextState(t, NewRulesetBuilder().PipelineRuleset(GameTypeStandard, NewPipeline(standardRulesetStages...)))
	}
}

func TestEatingOnLastMove(t *testing.T) {
	// We want to specifically ensure that snakes eating food on their last turn
	// survive. It used to be that this wasn't the case, and snakes were eliminated
	// if they moved onto food with their final move. This behaviour wasn't "wrong" or incorrect,
	// it just was less fun to watch. So let's ensure we're always giving snakes every possible
	// changes to reach food before eliminating them.
	tests := []struct {
		prevState     *BoardState
		moves         []SnakeMove
		expectedError error
		expectedState *BoardState
	}{
		{
			&BoardState{
				Width:  10,
				Height: 10,
				Snakes: []Snake{
					{
						ID:     "one",
						Body:   []Point{{X: 0, Y: 2}, {X: 0, Y: 1}, {X: 0, Y: 0}},
						Health: 1,
					},
					{
						ID:     "two",
						Body:   []Point{{X: 3, Y: 2}, {X: 3, Y: 3}, {X: 3, Y: 4}},
						Health: 1,
					},
				},
				Food: []Point{{X: 0, Y: 3}, {X: 9, Y: 9}},
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
						ID:     "one",
						Body:   []Point{{X: 0, Y: 3}, {X: 0, Y: 2}, {X: 0, Y: 1}, {X: 0, Y: 1}},
						Health: 100,
					},
					{
						ID:               "two",
						Body:             []Point{{X: 3, Y: 1}, {X: 3, Y: 2}, {X: 3, Y: 3}},
						Health:           0,
						EliminatedCause:  EliminatedByOutOfHealth,
						EliminatedOnTurn: 1,
					},
				},
				Food: []Point{{X: 9, Y: 9}},
			},
		},
	}

	rand.Seed(0) // Seed with a value that will reliably not spawn food
	r := getStandardRuleset(Settings{})
	for _, test := range tests {
		_, nextState, err := r.Execute(test.prevState, test.moves)
		require.Equal(t, err, test.expectedError)
		if test.expectedState != nil {
			require.Equal(t, test.expectedState.Width, nextState.Width)
			require.Equal(t, test.expectedState.Height, nextState.Height)
			require.Equal(t, test.expectedState.Food, nextState.Food)
			require.Equal(t, test.expectedState.Snakes, nextState.Snakes)
		}
	}
}

func TestHeadToHeadOnFood(t *testing.T) {
	// We want to specifically ensure that snakes that collide head-to-head
	// on top of food successfully remove the food - that's the core behaviour this test
	// is enforcing. There's a known side effect of this though, in that both snakes will
	// have eaten prior to being evaluated on the head-to-head (+1 length, full health).
	// We're okay with that since it does not impact the result of the head-to-head,
	// however that behaviour could change in the future and this test could be updated.
	tests := []struct {
		prevState     *BoardState
		moves         []SnakeMove
		expectedError error
		expectedState *BoardState
	}{
		{
			&BoardState{
				Turn:   41,
				Width:  10,
				Height: 10,
				Snakes: []Snake{
					{
						ID:     "one",
						Body:   []Point{{X: 0, Y: 2}, {X: 0, Y: 1}, {X: 0, Y: 0}},
						Health: 10,
					},
					{
						ID:     "two",
						Body:   []Point{{X: 0, Y: 4}, {X: 0, Y: 5}, {X: 0, Y: 6}},
						Health: 10,
					},
				},
				Food: []Point{{X: 0, Y: 3}, {X: 9, Y: 9}},
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
						Body:             []Point{{X: 0, Y: 3}, {X: 0, Y: 2}, {X: 0, Y: 1}, {X: 0, Y: 1}},
						Health:           100,
						EliminatedCause:  EliminatedByHeadToHeadCollision,
						EliminatedBy:     "two",
						EliminatedOnTurn: 42,
					},
					{
						ID:               "two",
						Body:             []Point{{X: 0, Y: 3}, {X: 0, Y: 4}, {X: 0, Y: 5}, {X: 0, Y: 5}},
						Health:           100,
						EliminatedCause:  EliminatedByHeadToHeadCollision,
						EliminatedBy:     "one",
						EliminatedOnTurn: 42,
					},
				},
				Food: []Point{{X: 9, Y: 9}},
			},
		},
		{
			&BoardState{
				Turn:   41,
				Width:  10,
				Height: 10,
				Snakes: []Snake{
					{
						ID:     "one",
						Body:   []Point{{X: 0, Y: 2}, {X: 0, Y: 1}, {X: 0, Y: 0}},
						Health: 10,
					},
					{
						ID:     "two",
						Body:   []Point{{X: 0, Y: 4}, {X: 0, Y: 5}, {X: 0, Y: 6}, {X: 0, Y: 7}},
						Health: 10,
					},
				},
				Food: []Point{{X: 0, Y: 3}, {X: 9, Y: 9}},
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
						Body:             []Point{{X: 0, Y: 3}, {X: 0, Y: 2}, {X: 0, Y: 1}, {X: 0, Y: 1}},
						Health:           100,
						EliminatedCause:  EliminatedByHeadToHeadCollision,
						EliminatedBy:     "two",
						EliminatedOnTurn: 42,
					},
					{
						ID:     "two",
						Body:   []Point{{X: 0, Y: 3}, {X: 0, Y: 4}, {X: 0, Y: 5}, {X: 0, Y: 6}, {X: 0, Y: 6}},
						Health: 100,
					},
				},
				Food: []Point{{X: 9, Y: 9}},
			},
		},
	}

	rand.Seed(0) // Seed with a value that will reliably not spawn food
	r := getStandardRuleset(Settings{})
	for _, test := range tests {
		_, nextState, err := r.Execute(test.prevState, test.moves)
		require.Equal(t, test.expectedError, err)
		if test.expectedState != nil {
			require.Equal(t, test.expectedState.Width, nextState.Width)
			require.Equal(t, test.expectedState.Height, nextState.Height)
			require.Equal(t, test.expectedState.Food, nextState.Food)
			require.Equal(t, test.expectedState.Snakes, nextState.Snakes)
		}
	}
}

func TestRegressionIssue19(t *testing.T) {
	// Eliminated snakes passed to CreateNextBoardState should not impact next game state
	tests := []struct {
		prevState     *BoardState
		moves         []SnakeMove
		expectedError error
		expectedState *BoardState
	}{
		{
			&BoardState{
				Width:  10,
				Height: 10,
				Snakes: []Snake{
					{
						ID:     "one",
						Body:   []Point{{X: 0, Y: 2}, {X: 0, Y: 1}, {X: 0, Y: 0}},
						Health: 100,
					},
					{
						ID:     "two",
						Body:   []Point{{X: 0, Y: 5}, {X: 0, Y: 6}, {X: 0, Y: 7}},
						Health: 100,
					},
					{
						ID:              "eliminated",
						Body:            []Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}, {X: 0, Y: 4}, {X: 0, Y: 5}, {X: 0, Y: 6}},
						Health:          0,
						EliminatedCause: EliminatedByOutOfHealth,
					},
				},
				Food: []Point{{X: 9, Y: 9}},
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
						ID:     "one",
						Body:   []Point{{X: 0, Y: 3}, {X: 0, Y: 2}, {X: 0, Y: 1}},
						Health: 99,
					},
					{
						ID:     "two",
						Body:   []Point{{X: 0, Y: 4}, {X: 0, Y: 5}, {X: 0, Y: 6}},
						Health: 99,
					},
					{
						ID:              "eliminated",
						Body:            []Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}, {X: 0, Y: 4}, {X: 0, Y: 5}, {X: 0, Y: 6}},
						Health:          0,
						EliminatedCause: EliminatedByOutOfHealth,
					},
				},
				Food: []Point{{X: 9, Y: 9}},
			},
		},
	}

	rand.Seed(0) // Seed with a value that will reliably not spawn food
	r := getStandardRuleset(Settings{})
	for _, test := range tests {
		_, nextState, err := r.Execute(test.prevState, test.moves)
		require.Equal(t, err, test.expectedError)
		if test.expectedState != nil {
			require.Equal(t, test.expectedState.Width, nextState.Width)
			require.Equal(t, test.expectedState.Height, nextState.Height)
			require.Equal(t, test.expectedState.Food, nextState.Food)
			require.Equal(t, test.expectedState.Snakes, nextState.Snakes)
		}
	}

}

func TestMoveSnakes(t *testing.T) {
	b := &BoardState{
		Snakes: []Snake{
			{
				ID:     "one",
				Body:   []Point{{X: 10, Y: 110}, {X: 11, Y: 110}},
				Health: 111111,
			},
			{
				ID:     "two",
				Body:   []Point{{X: 23, Y: 220}, {X: 22, Y: 220}, {X: 21, Y: 220}, {X: 20, Y: 220}},
				Health: 222222,
			},
			{
				ID:              "three",
				Body:            []Point{{X: 0, Y: 0}},
				Health:          1,
				EliminatedCause: EliminatedByOutOfBounds,
			},
		},
	}

	tests := []struct {
		MoveOne       string
		ExpectedOne   []Point
		MoveTwo       string
		ExpectedTwo   []Point
		MoveThree     string
		ExpectedThree []Point
	}{
		{
			MoveDown, []Point{{X: 10, Y: 109}, {X: 10, Y: 110}},
			MoveUp, []Point{{X: 23, Y: 221}, {X: 23, Y: 220}, {X: 22, Y: 220}, {X: 21, Y: 220}},
			MoveDown, []Point{{X: 0, Y: 0}},
		},
		{
			MoveRight, []Point{{X: 11, Y: 109}, {X: 10, Y: 109}},
			MoveLeft, []Point{{X: 22, Y: 221}, {X: 23, Y: 221}, {X: 23, Y: 220}, {X: 22, Y: 220}},
			MoveDown, []Point{{X: 0, Y: 0}},
		},
		{
			MoveRight, []Point{{X: 12, Y: 109}, {X: 11, Y: 109}},
			MoveLeft, []Point{{X: 21, Y: 221}, {X: 22, Y: 221}, {X: 23, Y: 221}, {X: 23, Y: 220}},
			MoveDown, []Point{{X: 0, Y: 0}},
		},
		{
			MoveRight, []Point{{X: 13, Y: 109}, {X: 12, Y: 109}},
			MoveLeft, []Point{{X: 20, Y: 221}, {X: 21, Y: 221}, {X: 22, Y: 221}, {X: 23, Y: 221}},
			MoveDown, []Point{{X: 0, Y: 0}},
		},
		{
			MoveDown, []Point{{X: 13, Y: 108}, {X: 13, Y: 109}},
			MoveUp, []Point{{X: 20, Y: 222}, {X: 20, Y: 221}, {X: 21, Y: 221}, {X: 22, Y: 221}},
			MoveDown, []Point{{X: 0, Y: 0}},
		},
	}

	r := getStandardRuleset(Settings{})
	for _, test := range tests {
		moves := []SnakeMove{
			{ID: "one", Move: test.MoveOne},
			{ID: "two", Move: test.MoveTwo},
			{ID: "three", Move: test.MoveThree},
		}
		_, err := MoveSnakesStandard(b, r.Settings(), moves)

		require.NoError(t, err)
		require.Len(t, b.Snakes, 3)

		require.Equal(t, 111111, b.Snakes[0].Health)
		require.Equal(t, 222222, b.Snakes[1].Health)
		require.Equal(t, 1, b.Snakes[2].Health)

		require.Len(t, b.Snakes[0].Body, 2)
		require.Len(t, b.Snakes[1].Body, 4)
		require.Len(t, b.Snakes[2].Body, 1)

		require.Equal(t, len(b.Snakes[0].Body), len(test.ExpectedOne))
		for i, e := range test.ExpectedOne {
			require.Equal(t, e, b.Snakes[0].Body[i])
		}
		require.Equal(t, len(b.Snakes[1].Body), len(test.ExpectedTwo))
		for i, e := range test.ExpectedTwo {
			require.Equal(t, e, b.Snakes[1].Body[i])
		}
		require.Equal(t, len(b.Snakes[2].Body), len(test.ExpectedThree))
		for i, e := range test.ExpectedThree {
			require.Equal(t, e, b.Snakes[2].Body[i])
		}
	}
}

func TestMoveSnakesWrongID(t *testing.T) {
	b := &BoardState{
		Snakes: []Snake{
			{
				ID:   "one",
				Body: []Point{{X: 1, Y: 1}},
			},
		},
	}
	moves := []SnakeMove{
		{
			ID:   "not found",
			Move: MoveUp,
		},
	}

	r := getStandardRuleset(Settings{})
	_, err := MoveSnakesStandard(b, r.Settings(), moves)
	require.Equal(t, ErrorNoMoveFound, err)
}

func TestMoveSnakesNotEnoughMoves(t *testing.T) {
	b := &BoardState{
		Snakes: []Snake{
			{
				ID:   "one",
				Body: []Point{{X: 1, Y: 1}},
			},
			{
				ID:   "two",
				Body: []Point{{X: 2, Y: 2}},
			},
		},
	}
	moves := []SnakeMove{
		{
			ID:   "two",
			Move: MoveUp,
		},
	}

	r := getStandardRuleset(Settings{})
	_, err := MoveSnakesStandard(b, r.Settings(), moves)
	require.Equal(t, ErrorNoMoveFound, err)
}

func TestMoveSnakesExtraMovesIgnored(t *testing.T) {
	b := &BoardState{
		Snakes: []Snake{
			{
				ID:   "one",
				Body: []Point{{X: 1, Y: 1}},
			},
		},
	}
	moves := []SnakeMove{
		{
			ID:   "one",
			Move: MoveDown,
		},
		{
			ID:   "two",
			Move: MoveLeft,
		},
	}

	r := getStandardRuleset(Settings{})
	_, err := MoveSnakesStandard(b, r.Settings(), moves)
	require.NoError(t, err)
	require.Equal(t, []Point{{X: 1, Y: 0}}, b.Snakes[0].Body)
}

func TestMoveSnakesDefault(t *testing.T) {
	tests := []struct {
		Body     []Point
		Move     string
		Expected []Point
	}{
		{
			Body:     []Point{{X: 0, Y: 0}},
			Move:     "invalid",
			Expected: []Point{{X: 0, Y: 1}},
		},
		{
			Body:     []Point{{X: 5, Y: 5}, {X: 5, Y: 5}},
			Move:     "",
			Expected: []Point{{X: 5, Y: 6}, {X: 5, Y: 5}},
		},
		{
			Body:     []Point{{X: 5, Y: 5}, {X: 5, Y: 4}},
			Expected: []Point{{X: 5, Y: 6}, {X: 5, Y: 5}},
		},
		{
			Body:     []Point{{X: 5, Y: 4}, {X: 5, Y: 5}},
			Expected: []Point{{X: 5, Y: 3}, {X: 5, Y: 4}},
		},
		{
			Body:     []Point{{X: 5, Y: 4}, {X: 5, Y: 5}},
			Expected: []Point{{X: 5, Y: 3}, {X: 5, Y: 4}},
		},
		{
			Body:     []Point{{X: 4, Y: 5}, {X: 5, Y: 5}},
			Expected: []Point{{X: 3, Y: 5}, {X: 4, Y: 5}},
		},
		{
			Body:     []Point{{X: 5, Y: 5}, {X: 4, Y: 5}},
			Expected: []Point{{X: 6, Y: 5}, {X: 5, Y: 5}},
		},
	}

	r := getStandardRuleset(Settings{})
	for _, test := range tests {
		b := &BoardState{
			Snakes: []Snake{
				{ID: "one", Body: test.Body},
			},
		}
		moves := []SnakeMove{{ID: "one", Move: test.Move}}

		_, err := MoveSnakesStandard(b, r.Settings(), moves)
		require.NoError(t, err)
		require.Len(t, b.Snakes, 1)
		require.Equal(t, len(test.Body), len(b.Snakes[0].Body))
		require.Equal(t, len(test.Expected), len(b.Snakes[0].Body))
		for i, e := range test.Expected {
			require.Equal(t, e, b.Snakes[0].Body[i])
		}
	}
}

func TestGetDefaultMove(t *testing.T) {
	tests := []struct {
		SnakeBody    []Point
		ExpectedMove string
	}{
		// Default is always up
		{
			SnakeBody:    []Point{},
			ExpectedMove: MoveUp,
		},
		{
			SnakeBody:    []Point{{X: 0, Y: 0}},
			ExpectedMove: MoveUp,
		},
		{
			SnakeBody:    []Point{{X: -1, Y: -1}},
			ExpectedMove: MoveUp,
		},
		// Stacked (fallback to default)
		{
			SnakeBody:    []Point{{X: 2, Y: 2}, {X: 2, Y: 2}},
			ExpectedMove: MoveUp,
		},
		// Neck next to head
		{
			SnakeBody:    []Point{{X: 2, Y: 2}, {X: 2, Y: 1}},
			ExpectedMove: MoveUp,
		},
		{
			SnakeBody:    []Point{{X: 2, Y: 2}, {X: 2, Y: 3}},
			ExpectedMove: MoveDown,
		},
		{
			SnakeBody:    []Point{{X: 2, Y: 2}, {X: 1, Y: 2}},
			ExpectedMove: MoveRight,
		},
		{
			SnakeBody:    []Point{{X: 2, Y: 2}, {X: 3, Y: 2}},
			ExpectedMove: MoveLeft,
		},
		// Board wrap cases
		{
			SnakeBody:    []Point{{X: 0, Y: 0}, {X: 0, Y: 2}},
			ExpectedMove: MoveUp,
		},
		{
			SnakeBody:    []Point{{X: 0, Y: 0}, {X: 2, Y: 0}},
			ExpectedMove: MoveRight,
		},
		{
			SnakeBody:    []Point{{X: 0, Y: 2}, {X: 0, Y: 0}},
			ExpectedMove: MoveDown,
		},
		{
			SnakeBody:    []Point{{X: 2, Y: 0}, {X: 0, Y: 0}},
			ExpectedMove: MoveLeft,
		},
	}

	for _, test := range tests {
		actualMove := getDefaultMove(test.SnakeBody)
		require.Equal(t, test.ExpectedMove, actualMove)
	}
}

func TestReduceSnakeHealth(t *testing.T) {
	b := &BoardState{
		Snakes: []Snake{
			{
				Body:   []Point{{X: 0, Y: 0}, {X: 0, Y: 1}},
				Health: 99,
			},
			{
				Body:   []Point{{X: 5, Y: 8}, {X: 6, Y: 8}, {X: 7, Y: 8}},
				Health: 2,
			},
			{
				Body:            []Point{{X: 0, Y: 0}, {X: 0, Y: 1}},
				Health:          50,
				EliminatedCause: EliminatedByCollision,
			},
		},
	}

	r := getStandardRuleset(Settings{})
	_, err := ReduceSnakeHealthStandard(b, r.Settings(), mockSnakeMoves())
	require.NoError(t, err)
	require.Equal(t, b.Snakes[0].Health, 98)
	require.Equal(t, b.Snakes[1].Health, 1)
	require.Equal(t, b.Snakes[2].Health, 50)

	_, err = ReduceSnakeHealthStandard(b, r.Settings(), mockSnakeMoves())
	require.NoError(t, err)
	require.Equal(t, b.Snakes[0].Health, 97)
	require.Equal(t, b.Snakes[1].Health, 0)
	require.Equal(t, b.Snakes[2].Health, 50)

	_, err = ReduceSnakeHealthStandard(b, r.Settings(), mockSnakeMoves())
	require.NoError(t, err)
	require.Equal(t, b.Snakes[0].Health, 96)
	require.Equal(t, b.Snakes[1].Health, -1)
	require.Equal(t, b.Snakes[2].Health, 50)

	_, err = ReduceSnakeHealthStandard(b, r.Settings(), mockSnakeMoves())
	require.NoError(t, err)
	require.Equal(t, b.Snakes[0].Health, 95)
	require.Equal(t, b.Snakes[1].Health, -2)
	require.Equal(t, b.Snakes[2].Health, 50)
}

func TestSnakeIsOutOfHealth(t *testing.T) {
	tests := []struct {
		Health   int
		Expected bool
	}{
		{Health: math.MinInt, Expected: true},
		{Health: -10, Expected: true},
		{Health: -2, Expected: true},
		{Health: -1, Expected: true},
		{Health: 0, Expected: true},
		{Health: 1, Expected: false},
		{Health: 2, Expected: false},
		{Health: 10, Expected: false},
		{Health: math.MaxInt, Expected: false},
	}

	for _, test := range tests {
		s := &Snake{Health: test.Health}
		require.Equal(t, test.Expected, snakeIsOutOfHealth(s), "Health: %+v", test.Health)
	}
}

func TestSnakeIsOutOfBounds(t *testing.T) {
	boardWidth := 10
	boardHeight := 100

	tests := []struct {
		Point    Point
		Expected bool
	}{
		{Point{X: math.MinInt, Y: math.MinInt}, true},
		{Point{X: math.MinInt, Y: 0}, true},
		{Point{X: 0, Y: math.MinInt}, true},
		{Point{X: -1, Y: -1}, true},
		{Point{X: -1, Y: 0}, true},
		{Point{X: 0, Y: -1}, true},
		{Point{X: 0, Y: 0}, false},
		{Point{X: 1, Y: 0}, false},
		{Point{X: 0, Y: 1}, false},
		{Point{X: 1, Y: 1}, false},
		{Point{X: 9, Y: 9}, false},
		{Point{X: 9, Y: 10}, false},
		{Point{X: 9, Y: 11}, false},
		{Point{X: 10, Y: 9}, true},
		{Point{X: 10, Y: 10}, true},
		{Point{X: 10, Y: 11}, true},
		{Point{X: 11, Y: 9}, true},
		{Point{X: 11, Y: 10}, true},
		{Point{X: 11, Y: 11}, true},
		{Point{X: math.MaxInt, Y: 11}, true},
		{Point{X: 9, Y: 99}, false},
		{Point{X: 9, Y: 100}, true},
		{Point{X: 9, Y: 101}, true},
		{Point{X: 9, Y: math.MaxInt}, true},
		{Point{X: math.MaxInt, Y: math.MaxInt}, true},
	}

	for _, test := range tests {
		// Test with point as head
		s := Snake{Body: []Point{test.Point}}
		require.Equal(t, test.Expected, snakeIsOutOfBounds(&s, boardWidth, boardHeight), "Head%+v", test.Point)
		// Test with point as body
		s = Snake{Body: []Point{{X: 0, Y: 0}, {X: 0, Y: 0}, test.Point}}
		require.Equal(t, test.Expected, snakeIsOutOfBounds(&s, boardWidth, boardHeight), "Body%+v", test.Point)
	}
}

func TestSnakeHasBodyCollidedSelf(t *testing.T) {
	tests := []struct {
		Body     []Point
		Expected bool
	}{
		{[]Point{{X: 1, Y: 1}}, false},
		// Self stacks should self collide
		// (we rely on snakes moving before we check self-collision on turn one)
		{[]Point{{X: 2, Y: 2}, {X: 2, Y: 2}}, true},
		{[]Point{{X: 3, Y: 3}, {X: 3, Y: 3}, {X: 3, Y: 3}}, true},
		{[]Point{{X: 5, Y: 5}, {X: 5, Y: 5}, {X: 5, Y: 5}, {X: 5, Y: 5}, {X: 5, Y: 5}}, true},
		// Non-collision cases
		{[]Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 0}}, false},
		{[]Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}}, false},
		{[]Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 0}}, false},
		{[]Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 0}, {X: 4, Y: 0}}, false},
		{[]Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}, false},
		{[]Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 2}, {X: 0, Y: 2}}, false},
		{[]Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}, {X: 0, Y: 4}}, false},
		// Collision cases
		{[]Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: 0}}, true},
		{[]Point{{X: 0, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0}}, true},
		{[]Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 1}, {X: 0, Y: 0}}, true},
		{[]Point{{X: 4, Y: 4}, {X: 3, Y: 4}, {X: 3, Y: 3}, {X: 4, Y: 4}, {X: 4, Y: 4}}, true},
		{[]Point{{X: 3, Y: 3}, {X: 3, Y: 4}, {X: 3, Y: 3}, {X: 4, Y: 4}, {X: 4, Y: 5}}, true},
	}

	for _, test := range tests {
		s := Snake{Body: test.Body}
		require.Equal(t, test.Expected, snakeHasBodyCollided(&s, &s), "Body%q", s.Body)
	}
}

func TestSnakeHasBodyCollidedOther(t *testing.T) {
	tests := []struct {
		SnakeBody []Point
		OtherBody []Point
		Expected  bool
	}{
		{
			// Just heads
			[]Point{{X: 0, Y: 0}},
			[]Point{{X: 1, Y: 1}},
			false,
		},
		{
			// Head-to-heads are not considered in body collisions
			[]Point{{X: 0, Y: 0}},
			[]Point{{X: 0, Y: 0}},
			false,
		},
		{
			// Stacked bodies
			[]Point{{X: 0, Y: 0}},
			[]Point{{X: 0, Y: 0}, {X: 0, Y: 0}},
			true,
		},
		{
			// Separate stacked bodies
			[]Point{{X: 0, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: 0}},
			[]Point{{X: 1, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 1}},
			false,
		},
		{
			// Stacked bodies, separated heads
			[]Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 0}},
			[]Point{{X: 2, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 0}},
			false,
		},
		{
			// Mid-snake collision
			[]Point{{X: 1, Y: 1}},
			[]Point{{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}},
			true,
		},
	}

	for _, test := range tests {
		s := &Snake{Body: test.SnakeBody}
		o := &Snake{Body: test.OtherBody}
		require.Equal(t, test.Expected, snakeHasBodyCollided(s, o), "Snake%q Other%q", s.Body, o.Body)
	}
}

func TestSnakeHasLostHeadToHead(t *testing.T) {
	tests := []struct {
		SnakeBody        []Point
		OtherBody        []Point
		Expected         bool
		ExpectedOpposite bool
	}{
		{
			// Just heads
			[]Point{{X: 0, Y: 0}},
			[]Point{{X: 1, Y: 1}},
			false, false,
		},
		{
			// Just heads colliding
			[]Point{{X: 0, Y: 0}},
			[]Point{{X: 0, Y: 0}},
			true, true,
		},
		{
			// One snake larger
			[]Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
			[]Point{{X: 0, Y: 0}},
			false, true,
		},
		{
			// Other snake equal
			[]Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
			[]Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}},
			true, true,
		},
		{
			// Other snake longer
			[]Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}},
			[]Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}},
			true, false,
		},
		{
			// Body collision
			[]Point{{X: 0, Y: 1}, {X: 1, Y: 1}, {X: 2, Y: 1}},
			[]Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}},
			false, false,
		},
		{
			// Separate stacked bodies, head collision
			[]Point{{X: 3, Y: 10}, {X: 2, Y: 10}, {X: 2, Y: 10}},
			[]Point{{X: 3, Y: 10}, {X: 4, Y: 10}, {X: 4, Y: 10}},
			true, true,
		},
		{
			// Separate stacked bodies, head collision
			[]Point{{X: 10, Y: 3}, {X: 10, Y: 2}, {X: 10, Y: 1}, {X: 10, Y: 0}},
			[]Point{{X: 10, Y: 3}, {X: 10, Y: 4}, {X: 10, Y: 5}},
			false, true,
		},
	}

	for _, test := range tests {
		s := Snake{Body: test.SnakeBody}
		o := Snake{Body: test.OtherBody}
		require.Equal(t, test.Expected, snakeHasLostHeadToHead(&s, &o), "Snake%q Other%q", s.Body, o.Body)
		require.Equal(t, test.ExpectedOpposite, snakeHasLostHeadToHead(&o, &s), "Snake%q Other%q", s.Body, o.Body)
	}

}

func TestMaybeEliminateSnakes(t *testing.T) {
	tests := []struct {
		Name                     string
		Snakes                   []Snake
		ExpectedEliminatedCauses []string
		ExpectedEliminatedBy     []string
		Err                      error
	}{
		{
			"Empty",
			[]Snake{},
			[]string{},
			[]string{},
			nil,
		},
		{
			"Zero Snake",
			[]Snake{
				{},
			},
			[]string{NotEliminated},
			[]string{""},
			ErrorZeroLengthSnake,
		},
		{
			"Single Starvation",
			[]Snake{
				{ID: "1", Body: []Point{{X: 1, Y: 1}}},
			},
			[]string{EliminatedByOutOfHealth},
			[]string{""},
			nil,
		},
		{
			"Not Eliminated",
			[]Snake{
				{ID: "1", Health: 1, Body: []Point{{X: 1, Y: 1}}},
			},
			[]string{NotEliminated},
			[]string{""},
			nil,
		},
		{
			"Out of Bounds",
			[]Snake{
				{ID: "1", Health: 1, Body: []Point{{X: -1, Y: 1}}},
			},
			[]string{EliminatedByOutOfBounds},
			[]string{""},
			nil,
		},
		{
			"Self Collision",
			[]Snake{
				{ID: "1", Health: 1, Body: []Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 0}}},
			},
			[]string{EliminatedBySelfCollision},
			[]string{"1"},
			nil,
		},
		{
			"Multiple Separate Deaths",
			[]Snake{
				{ID: "1", Health: 1, Body: []Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 0}}},
				{ID: "2", Health: 1, Body: []Point{{X: -1, Y: 1}}},
			},
			[]string{
				EliminatedBySelfCollision,
				EliminatedByOutOfBounds},
			[]string{"1", ""},
			nil,
		},
		{
			"Other Collision",
			[]Snake{
				{ID: "1", Health: 1, Body: []Point{{X: 0, Y: 2}, {X: 0, Y: 3}, {X: 0, Y: 4}}},
				{ID: "2", Health: 1, Body: []Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}},
			},
			[]string{
				EliminatedByCollision,
				NotEliminated},
			[]string{"2", ""},
			nil,
		},
		{
			"All Eliminated Head 2 Head",
			[]Snake{
				{ID: "1", Health: 1, Body: []Point{{X: 1, Y: 1}}},
				{ID: "2", Health: 1, Body: []Point{{X: 1, Y: 1}}},
				{ID: "3", Health: 1, Body: []Point{{X: 1, Y: 1}}},
			},
			[]string{
				EliminatedByHeadToHeadCollision,
				EliminatedByHeadToHeadCollision,
				EliminatedByHeadToHeadCollision,
			},
			[]string{"2", "1", "1"},
			nil,
		},
		{
			"One Snake wins Head 2 Head",
			[]Snake{
				{ID: "1", Health: 1, Body: []Point{{X: 1, Y: 1}, {X: 0, Y: 1}}},
				{ID: "2", Health: 1, Body: []Point{{X: 1, Y: 1}, {X: 1, Y: 2}, {X: 1, Y: 3}}},
				{ID: "3", Health: 1, Body: []Point{{X: 1, Y: 1}}},
			},
			[]string{
				EliminatedByHeadToHeadCollision,
				NotEliminated,
				EliminatedByHeadToHeadCollision,
			},
			[]string{"2", "", "2"},
			nil,
		},
		{
			"All Snakes Body Eliminated",
			[]Snake{
				{ID: "1", Health: 1, Body: []Point{{X: 4, Y: 4}, {X: 3, Y: 3}}},
				{ID: "2", Health: 1, Body: []Point{{X: 3, Y: 3}, {X: 2, Y: 2}}},
				{ID: "3", Health: 1, Body: []Point{{X: 2, Y: 2}, {X: 1, Y: 1}}},
				{ID: "4", Health: 1, Body: []Point{{X: 1, Y: 1}, {X: 4, Y: 4}}},
				{ID: "5", Health: 1, Body: []Point{{X: 4, Y: 4}}}, // Body collision takes priority
			},
			[]string{
				EliminatedByCollision,
				EliminatedByCollision,
				EliminatedByCollision,
				EliminatedByCollision,
				EliminatedByCollision,
			},
			[]string{"4", "1", "2", "3", "4"},
			nil,
		},
		{
			"All Snakes Eliminated Head 2 Head",
			[]Snake{
				{ID: "1", Health: 1, Body: []Point{{X: 4, Y: 4}, {X: 4, Y: 5}}},
				{ID: "2", Health: 1, Body: []Point{{X: 4, Y: 4}, {X: 4, Y: 3}}},
				{ID: "3", Health: 1, Body: []Point{{X: 4, Y: 4}, {X: 5, Y: 4}}},
				{ID: "4", Health: 1, Body: []Point{{X: 4, Y: 4}, {X: 3, Y: 4}}},
			},
			[]string{
				EliminatedByHeadToHeadCollision,
				EliminatedByHeadToHeadCollision,
				EliminatedByHeadToHeadCollision,
				EliminatedByHeadToHeadCollision,
			},
			[]string{"2", "1", "1", "1"},
			nil,
		},
		{
			"4 Snakes Head 2 Head",
			[]Snake{
				{ID: "1", Health: 1, Body: []Point{{X: 4, Y: 4}, {X: 4, Y: 5}}},
				{ID: "2", Health: 1, Body: []Point{{X: 4, Y: 4}, {X: 4, Y: 3}}},
				{ID: "3", Health: 1, Body: []Point{{X: 4, Y: 4}, {X: 5, Y: 4}, {X: 6, Y: 4}}},
				{ID: "4", Health: 1, Body: []Point{{X: 4, Y: 4}, {X: 3, Y: 4}}},
			},
			[]string{
				EliminatedByHeadToHeadCollision,
				EliminatedByHeadToHeadCollision,
				NotEliminated,
				EliminatedByHeadToHeadCollision,
			},
			[]string{"3", "3", "", "3"},
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			b := &BoardState{
				Width:  10,
				Height: 10,
				Snakes: test.Snakes,
			}
			_, err := EliminateSnakesStandard(b, Settings{}, mockSnakeMoves())
			require.Equal(t, test.Err, err)
			for i, snake := range b.Snakes {
				require.Equal(t, test.ExpectedEliminatedCauses[i], snake.EliminatedCause)
				require.Equal(t, test.ExpectedEliminatedBy[i], snake.EliminatedBy)
			}
		})
	}
}

func TestMaybeEliminateSnakesPriority(t *testing.T) {
	tests := []struct {
		Snakes                   []Snake
		ExpectedEliminatedCauses []string
		ExpectedEliminatedBy     []string
	}{
		{
			[]Snake{
				{ID: "1", Health: 0, Body: []Point{{X: -1, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0}}},
				{ID: "2", Health: 1, Body: []Point{{X: -1, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0}}},
				{ID: "3", Health: 1, Body: []Point{{X: 1, Y: 0}, {X: 0, Y: 0}, {X: 1, Y: 0}}},
				{ID: "4", Health: 1, Body: []Point{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 1, Y: 2}}},
				{ID: "5", Health: 1, Body: []Point{{X: 2, Y: 2}, {X: 2, Y: 1}, {X: 2, Y: 0}}},
				{ID: "6", Health: 1, Body: []Point{{X: 2, Y: 2}, {X: 2, Y: 3}, {X: 2, Y: 4}, {X: 2, Y: 5}}},
			},
			[]string{
				EliminatedByOutOfHealth,
				EliminatedByOutOfBounds,
				EliminatedBySelfCollision,
				EliminatedByCollision,
				EliminatedByHeadToHeadCollision,
				NotEliminated,
			},
			[]string{"", "", "3", "3", "6", ""},
		},
	}

	r := getStandardRuleset(Settings{})
	for _, test := range tests {
		b := &BoardState{Width: 10, Height: 10, Snakes: test.Snakes}
		_, err := EliminateSnakesStandard(b, r.Settings(), mockSnakeMoves())
		require.NoError(t, err)
		for i, snake := range b.Snakes {
			require.Equal(t, test.ExpectedEliminatedCauses[i], snake.EliminatedCause, snake.ID)
			require.Equal(t, test.ExpectedEliminatedBy[i], snake.EliminatedBy, snake.ID)
		}
	}
}

func TestMaybeDamageHazards(t *testing.T) {
	tests := []struct {
		Snakes                    []Snake
		Hazards                   []Point
		Food                      []Point
		ExpectedEliminatedCauses  []string
		ExpectedEliminatedByIDs   []string
		ExpectedEliminatedOnTurns []int
	}{
		{},
		{
			Snakes:                    []Snake{{Body: []Point{{X: 0, Y: 0}}}},
			Hazards:                   []Point{},
			ExpectedEliminatedCauses:  []string{NotEliminated},
			ExpectedEliminatedByIDs:   []string{""},
			ExpectedEliminatedOnTurns: []int{0},
		},
		{
			Snakes:                    []Snake{{Body: []Point{{X: 0, Y: 0}}}},
			Hazards:                   []Point{{X: 0, Y: 0}},
			ExpectedEliminatedCauses:  []string{EliminatedByHazard},
			ExpectedEliminatedByIDs:   []string{""},
			ExpectedEliminatedOnTurns: []int{42},
		},
		{
			Snakes:                    []Snake{{Body: []Point{{X: 0, Y: 0}}}},
			Hazards:                   []Point{{X: 0, Y: 0}},
			Food:                      []Point{{X: 0, Y: 0}},
			ExpectedEliminatedCauses:  []string{NotEliminated},
			ExpectedEliminatedByIDs:   []string{""},
			ExpectedEliminatedOnTurns: []int{0},
		},
		{
			Snakes:                    []Snake{{Body: []Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}}}},
			Hazards:                   []Point{{X: 1, Y: 0}, {X: 2, Y: 0}},
			ExpectedEliminatedCauses:  []string{NotEliminated},
			ExpectedEliminatedByIDs:   []string{""},
			ExpectedEliminatedOnTurns: []int{0},
		},
		{
			Snakes: []Snake{
				{Body: []Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}}},
				{Body: []Point{{X: 3, Y: 3}, {X: 3, Y: 4}, {X: 3, Y: 5}, {X: 3, Y: 6}}},
			},
			Hazards:                   []Point{{X: 1, Y: 0}, {X: 2, Y: 0}, {X: 3, Y: 4}, {X: 3, Y: 5}, {X: 3, Y: 6}},
			ExpectedEliminatedCauses:  []string{NotEliminated, NotEliminated},
			ExpectedEliminatedByIDs:   []string{"", ""},
			ExpectedEliminatedOnTurns: []int{0, 0},
		},
		{
			Snakes: []Snake{
				{Body: []Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}}},
				{Body: []Point{{X: 3, Y: 3}, {X: 3, Y: 4}, {X: 3, Y: 5}, {X: 3, Y: 6}}},
			},
			Hazards:                   []Point{{X: 3, Y: 3}},
			ExpectedEliminatedCauses:  []string{NotEliminated, EliminatedByHazard},
			ExpectedEliminatedByIDs:   []string{"", ""},
			ExpectedEliminatedOnTurns: []int{0, 42},
		},
	}

	for _, test := range tests {
		b := &BoardState{Turn: 41, Snakes: test.Snakes, Hazards: test.Hazards, Food: test.Food}
		r := getStandardRuleset(NewSettingsWithParams(ParamHazardDamagePerTurn, "100"))
		_, err := DamageHazardsStandard(b, r.Settings(), mockSnakeMoves())
		require.NoError(t, err)

		for i, snake := range b.Snakes {
			require.Equal(t, test.ExpectedEliminatedCauses[i], snake.EliminatedCause)
			require.Equal(t, test.ExpectedEliminatedByIDs[i], snake.EliminatedBy)
			require.Equal(t, test.ExpectedEliminatedOnTurns[i], snake.EliminatedOnTurn)
		}

	}
}

func TestHazardDamagePerTurn(t *testing.T) {
	tests := []struct {
		Health                   int
		HazardDamagePerTurn      int
		Food                     bool
		ExpectedHealth           int
		ExpectedEliminationCause string
		Error                    error
	}{
		{100, 1, false, 99, NotEliminated, nil},
		{100, 1, true, 100, NotEliminated, nil},
		{100, 99, false, 1, NotEliminated, nil},
		{100, 99, true, 100, NotEliminated, nil},
		{100, -1, false, 100, NotEliminated, nil},
		{99, -2, false, 100, NotEliminated, nil},
		{100, 100, false, 0, EliminatedByHazard, nil},
		{100, 101, false, 0, EliminatedByHazard, nil},
		{100, 999, false, 0, EliminatedByHazard, nil},
		{100, 100, true, 100, NotEliminated, nil},
		{2, 1, false, 1, NotEliminated, nil},
		{1, 1, false, 0, EliminatedByHazard, nil},
		{1, 999, false, 0, EliminatedByHazard, nil},
		{0, 1, false, 0, EliminatedByHazard, nil},
		{0, 999, false, 0, EliminatedByHazard, nil},
	}

	for _, test := range tests {
		b := &BoardState{Snakes: []Snake{{Health: test.Health, Body: []Point{{X: 0, Y: 0}}}}, Hazards: []Point{{X: 0, Y: 0}}}
		if test.Food {
			b.Food = []Point{{X: 0, Y: 0}}
		}
		r := getStandardRuleset(NewSettingsWithParams(ParamHazardDamagePerTurn, fmt.Sprint(test.HazardDamagePerTurn)))

		_, err := DamageHazardsStandard(b, r.Settings(), mockSnakeMoves())
		require.Equal(t, test.Error, err)
		require.Equal(t, test.ExpectedHealth, b.Snakes[0].Health)
		require.Equal(t, test.ExpectedEliminationCause, b.Snakes[0].EliminatedCause)
	}
}

func TestMaybeFeedSnakes(t *testing.T) {
	tests := []struct {
		Name           string
		Snakes         []Snake
		Food           []Point
		ExpectedSnakes []Snake
		ExpectedFood   []Point
	}{
		{
			Name: "snake not on food",
			Snakes: []Snake{
				{Health: 5, Body: []Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}},
			},
			Food: []Point{{X: 3, Y: 3}},
			ExpectedSnakes: []Snake{
				{Health: 5, Body: []Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}},
			},
			ExpectedFood: []Point{{X: 3, Y: 3}},
		},
		{
			Name: "snake on food",
			Snakes: []Snake{
				{Health: SnakeMaxHealth - 1, Body: []Point{{X: 2, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 2}, {X: 2, Y: 2}}},
			},
			Food: []Point{{X: 2, Y: 1}},
			ExpectedSnakes: []Snake{
				{Health: SnakeMaxHealth, Body: []Point{{X: 2, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 2}, {X: 2, Y: 2}, {X: 2, Y: 2}}},
			},
			ExpectedFood: []Point{},
		},
		{
			Name: "food under body",
			Snakes: []Snake{
				{Body: []Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}},
			},
			Food: []Point{{X: 0, Y: 1}},
			ExpectedSnakes: []Snake{
				{Body: []Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}},
			},
			ExpectedFood: []Point{{X: 0, Y: 1}},
		},
		{
			Name: "snake on food but already eliminated",
			Snakes: []Snake{
				{Body: []Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}, EliminatedCause: "EliminatedByOutOfBounds"},
			},
			Food: []Point{{X: 0, Y: 0}},
			ExpectedSnakes: []Snake{
				{Body: []Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}},
			},
			ExpectedFood: []Point{{X: 0, Y: 0}},
		},
		{
			Name: "multiple snakes on same food",
			Snakes: []Snake{
				{Health: SnakeMaxHealth, Body: []Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}},
				{Health: SnakeMaxHealth - 9, Body: []Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}}},
			},
			Food: []Point{{X: 0, Y: 0}, {X: 4, Y: 4}},
			ExpectedSnakes: []Snake{
				{Health: SnakeMaxHealth, Body: []Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 2}}},
				{Health: SnakeMaxHealth, Body: []Point{{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 0}}},
			},
			ExpectedFood: []Point{{X: 4, Y: 4}},
		},
	}

	r := getStandardRuleset(Settings{})
	for _, test := range tests {
		b := &BoardState{
			Snakes: test.Snakes,
			Food:   test.Food,
		}
		_, err := FeedSnakesStandard(b, r.Settings(), nil)
		require.NoError(t, err, test.Name)
		require.Equal(t, len(test.ExpectedSnakes), len(b.Snakes), test.Name)
		for i := 0; i < len(b.Snakes); i++ {
			require.Equal(t, test.ExpectedSnakes[i].Health, b.Snakes[i].Health, test.Name)
			require.Equal(t, test.ExpectedSnakes[i].Body, b.Snakes[i].Body, test.Name)
		}
		require.Equal(t, test.ExpectedFood, b.Food, test.Name)
	}
}

func TestMaybeSpawnFoodMinimum(t *testing.T) {
	tests := []struct {
		MinimumFood  int
		Food         []Point
		ExpectedFood int
	}{
		// Use pre-tested seeds and results
		{0, []Point{}, 0},
		{1, []Point{}, 1},
		{9, []Point{}, 9},
		{7, []Point{{X: 4, Y: 5}, {X: 4, Y: 4}, {X: 4, Y: 1}}, 7},
	}

	for _, test := range tests {
		r := getStandardRuleset(NewSettingsWithParams(ParamMinimumFood, fmt.Sprint(test.MinimumFood)))
		b := &BoardState{
			Height: 11,
			Width:  11,
			Snakes: []Snake{
				{Body: []Point{{X: 1, Y: 0}, {X: 1, Y: 1}}},
				{Body: []Point{{X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}}},
			},
			Food: test.Food,
		}

		_, err := SpawnFoodStandard(b, r.Settings(), mockSnakeMoves())
		require.NoError(t, err)
		require.Equal(t, test.ExpectedFood, len(b.Food))
	}
}

func TestMaybeSpawnFoodZeroChance(t *testing.T) {
	r := getStandardRuleset(NewSettingsWithParams(ParamFoodSpawnChance, "0"))
	b := &BoardState{
		Height: 11,
		Width:  11,
		Snakes: []Snake{
			{Body: []Point{{X: 1, Y: 0}, {X: 1, Y: 1}}},
			{Body: []Point{{X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}}},
		},
		Food: []Point{},
	}
	for i := 0; i < 1000; i++ {
		_, err := SpawnFoodStandard(b, r.Settings(), nil)
		require.NoError(t, err)
		require.Equal(t, len(b.Food), 0)
	}
}

func TestMaybeSpawnFoodHundredChance(t *testing.T) {
	r := getStandardRuleset(NewSettingsWithParams(ParamFoodSpawnChance, "100"))
	b := &BoardState{
		Height: 11,
		Width:  11,
		Snakes: []Snake{
			{Body: []Point{{X: 1, Y: 0}, {X: 1, Y: 1}}},
			{Body: []Point{{X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}}},
		},
		Food: []Point{},
	}
	for i := 1; i <= 22; i++ {
		_, err := SpawnFoodStandard(b, r.Settings(), mockSnakeMoves())
		require.NoError(t, err)
		require.Equal(t, i, len(b.Food))
	}
}

func TestMaybeSpawnFoodHalfChance(t *testing.T) {
	tests := []struct {
		Seed         int64
		Food         []Point
		ExpectedFood int
	}{
		// Use pre-tested seeds and results
		{123, []Point{}, 1},
		{12345, []Point{}, 0},
		{456, []Point{{X: 4, Y: 4}}, 1},
		{789, []Point{{X: 4, Y: 4}}, 2},
		{511, []Point{{X: 4, Y: 4}}, 1},
		{165, []Point{{X: 4, Y: 4}}, 2},
	}

	r := getStandardRuleset(NewSettingsWithParams(ParamFoodSpawnChance, "50"))
	for _, test := range tests {
		b := &BoardState{
			Height: 4,
			Width:  5,
			Snakes: []Snake{
				{Body: []Point{{X: 1, Y: 0}, {X: 1, Y: 1}}},
				{Body: []Point{{X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}}},
			},
			Food: test.Food,
		}

		rand.Seed(test.Seed)
		_, err := SpawnFoodStandard(b, r.Settings(), mockSnakeMoves())
		require.NoError(t, err)
		require.Equal(t, test.ExpectedFood, len(b.Food), "Seed %d", test.Seed)
	}
}

func TestIsGameOver(t *testing.T) {
	tests := []struct {
		Snakes   []Snake
		Expected bool
	}{
		{[]Snake{}, true},
		{[]Snake{{}}, true},
		{[]Snake{{}, {}}, false},
		{[]Snake{{}, {}, {}, {}, {}}, false},
		{
			[]Snake{
				{EliminatedCause: EliminatedByCollision},
				{EliminatedCause: NotEliminated},
			},
			true,
		},
		{
			[]Snake{
				{EliminatedCause: NotEliminated},
				{EliminatedCause: EliminatedByCollision},
				{EliminatedCause: NotEliminated},
				{EliminatedCause: NotEliminated},
			},
			false,
		},
		{
			[]Snake{
				{EliminatedCause: EliminatedByOutOfBounds},
				{EliminatedCause: EliminatedByOutOfBounds},
				{EliminatedCause: EliminatedByOutOfBounds},
				{EliminatedCause: EliminatedByOutOfBounds},
			},
			true,
		},
		{
			[]Snake{
				{EliminatedCause: EliminatedByOutOfBounds},
				{EliminatedCause: EliminatedByOutOfBounds},
				{EliminatedCause: EliminatedByOutOfBounds},
				{EliminatedCause: NotEliminated},
			},
			true,
		},
		{
			[]Snake{
				{EliminatedCause: EliminatedByOutOfBounds},
				{EliminatedCause: EliminatedByOutOfBounds},
				{EliminatedCause: NotEliminated},
				{EliminatedCause: NotEliminated},
			},
			false,
		},
	}

	r := getStandardRuleset(Settings{})
	for _, test := range tests {
		b := &BoardState{
			Height: 11,
			Width:  11,
			Snakes: test.Snakes,
			Food:   []Point{},
		}

		actual, _, err := r.Execute(b, nil)
		require.NoError(t, err)
		require.Equal(t, test.Expected, actual)
	}
}
