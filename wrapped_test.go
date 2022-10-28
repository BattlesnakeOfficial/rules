package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getWrappedRuleset(settings Settings) Ruleset {
	return NewRulesetBuilder().WithSettings(settings).NamedRuleset(GameTypeWrapped)
}

func TestLeft(t *testing.T) {
	boardState := &BoardState{
		Width:  11,
		Height: 11,
		Snakes: []Snake{
			{ID: "bottomLeft", Health: 10, Body: []Point{{X: 0, Y: 0}}},
			{ID: "bottomRight", Health: 10, Body: []Point{{X: 10, Y: 0}}},
			{ID: "topLeft", Health: 10, Body: []Point{{X: 0, Y: 10}}},
			{ID: "topRight", Health: 10, Body: []Point{{X: 10, Y: 10}}},
		},
	}

	snakeMoves := []SnakeMove{
		{ID: "bottomLeft", Move: "left"},
		{ID: "bottomRight", Move: "left"},
		{ID: "topLeft", Move: "left"},
		{ID: "topRight", Move: "left"},
	}

	r := getWrappedRuleset(Settings{})

	gameOver, nextBoardState, err := r.Execute(boardState, snakeMoves)
	require.NoError(t, err)
	require.False(t, gameOver)
	require.Equal(t, len(boardState.Snakes), len(nextBoardState.Snakes))

	expectedSnakes := []Snake{
		{ID: "bottomLeft", Health: 10, Body: []Point{{X: 10, Y: 0}}},
		{ID: "bottomRight", Health: 10, Body: []Point{{X: 9, Y: 0}}},
		{ID: "topLeft", Health: 10, Body: []Point{{X: 10, Y: 10}}},
		{ID: "topRight", Health: 10, Body: []Point{{X: 9, Y: 10}}},
	}
	for i, snake := range nextBoardState.Snakes {
		require.Equal(t, expectedSnakes[i].ID, snake.ID, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedCause, snake.EliminatedCause, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedBy, snake.EliminatedBy, snake.ID)
		require.Equal(t, expectedSnakes[i].Body, snake.Body, snake.ID)
	}
}

func TestRight(t *testing.T) {
	boardState := &BoardState{
		Width:  11,
		Height: 11,
		Snakes: []Snake{
			{ID: "bottomLeft", Health: 10, Body: []Point{{X: 0, Y: 0}}},
			{ID: "bottomRight", Health: 10, Body: []Point{{X: 10, Y: 0}}},
			{ID: "topLeft", Health: 10, Body: []Point{{X: 0, Y: 10}}},
			{ID: "topRight", Health: 10, Body: []Point{{X: 10, Y: 10}}},
		},
	}

	snakeMoves := []SnakeMove{
		{ID: "bottomLeft", Move: "right"},
		{ID: "bottomRight", Move: "right"},
		{ID: "topLeft", Move: "right"},
		{ID: "topRight", Move: "right"},
	}

	r := getWrappedRuleset(Settings{})

	gameOver, nextBoardState, err := r.Execute(boardState, snakeMoves)
	require.NoError(t, err)
	require.False(t, gameOver)
	require.Equal(t, len(boardState.Snakes), len(nextBoardState.Snakes))

	expectedSnakes := []Snake{
		{ID: "bottomLeft", Health: 10, Body: []Point{{X: 1, Y: 0}}},
		{ID: "bottomRight", Health: 10, Body: []Point{{X: 0, Y: 0}}},
		{ID: "topLeft", Health: 10, Body: []Point{{X: 1, Y: 10}}},
		{ID: "topRight", Health: 10, Body: []Point{{X: 0, Y: 10}}},
	}
	for i, snake := range nextBoardState.Snakes {
		require.Equal(t, expectedSnakes[i].ID, snake.ID, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedCause, snake.EliminatedCause, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedBy, snake.EliminatedBy, snake.ID)
		require.Equal(t, expectedSnakes[i].Body, snake.Body, snake.ID)
	}
}

func TestUp(t *testing.T) {
	boardState := &BoardState{
		Width:  11,
		Height: 11,
		Snakes: []Snake{
			{ID: "bottomLeft", Health: 10, Body: []Point{{X: 0, Y: 0}}},
			{ID: "bottomRight", Health: 10, Body: []Point{{X: 10, Y: 0}}},
			{ID: "topLeft", Health: 10, Body: []Point{{X: 0, Y: 10}}},
			{ID: "topRight", Health: 10, Body: []Point{{X: 10, Y: 10}}},
		},
	}

	snakeMoves := []SnakeMove{
		{ID: "bottomLeft", Move: "up"},
		{ID: "bottomRight", Move: "up"},
		{ID: "topLeft", Move: "up"},
		{ID: "topRight", Move: "up"},
	}

	r := getWrappedRuleset(Settings{})

	gameOver, nextBoardState, err := r.Execute(boardState, snakeMoves)
	require.NoError(t, err)
	require.False(t, gameOver)
	require.Equal(t, len(boardState.Snakes), len(nextBoardState.Snakes))

	expectedSnakes := []Snake{
		{ID: "bottomLeft", Health: 10, Body: []Point{{X: 0, Y: 1}}},
		{ID: "bottomRight", Health: 10, Body: []Point{{X: 10, Y: 1}}},
		{ID: "topLeft", Health: 10, Body: []Point{{X: 0, Y: 0}}},
		{ID: "topRight", Health: 10, Body: []Point{{X: 10, Y: 0}}},
	}
	for i, snake := range nextBoardState.Snakes {
		require.Equal(t, expectedSnakes[i].ID, snake.ID, snake.ID)
		require.Equal(t, expectedSnakes[i].Body, snake.Body, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedCause, snake.EliminatedCause, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedBy, snake.EliminatedBy, snake.ID)
	}
}

func TestDown(t *testing.T) {
	boardState := &BoardState{
		Width:  11,
		Height: 11,
		Snakes: []Snake{
			{ID: "bottomLeft", Health: 10, Body: []Point{{X: 0, Y: 0}}},
			{ID: "bottomRight", Health: 10, Body: []Point{{X: 10, Y: 0}}},
			{ID: "topLeft", Health: 10, Body: []Point{{X: 0, Y: 10}}},
			{ID: "topRight", Health: 10, Body: []Point{{X: 10, Y: 10}}},
		},
	}

	snakeMoves := []SnakeMove{
		{ID: "bottomLeft", Move: "down"},
		{ID: "bottomRight", Move: "down"},
		{ID: "topLeft", Move: "down"},
		{ID: "topRight", Move: "down"},
	}

	r := getWrappedRuleset(Settings{})

	gameOver, nextBoardState, err := r.Execute(boardState, snakeMoves)
	require.NoError(t, err)
	require.False(t, gameOver)
	require.Equal(t, len(boardState.Snakes), len(nextBoardState.Snakes))

	expectedSnakes := []Snake{
		{ID: "bottomLeft", Health: 10, Body: []Point{{X: 0, Y: 10}}},
		{ID: "bottomRight", Health: 10, Body: []Point{{X: 10, Y: 10}}},
		{ID: "topLeft", Health: 10, Body: []Point{{X: 0, Y: 9}}},
		{ID: "topRight", Health: 10, Body: []Point{{X: 10, Y: 9}}},
	}
	for i, snake := range nextBoardState.Snakes {
		require.Equal(t, expectedSnakes[i].ID, snake.ID, snake.ID)
		require.Equal(t, expectedSnakes[i].Body, snake.Body, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedCause, snake.EliminatedCause, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedBy, snake.EliminatedBy, snake.ID)
	}
}

func TestEdgeCrossingCollision(t *testing.T) {
	boardState := &BoardState{
		Width:  11,
		Height: 11,
		Snakes: []Snake{
			{ID: "left", Health: 10, Body: []Point{{X: 0, Y: 5}}},
			{ID: "rightEdge", Health: 10, Body: []Point{
				{X: 10, Y: 1},
				{X: 10, Y: 2},
				{X: 10, Y: 3},
				{X: 10, Y: 4},
				{X: 10, Y: 5},
				{X: 10, Y: 6},
			}},
		},
	}

	snakeMoves := []SnakeMove{
		{ID: "left", Move: "left"},
		{ID: "rightEdge", Move: "down"},
	}

	r := getWrappedRuleset(Settings{})

	gameOver, nextBoardState, err := r.Execute(boardState, snakeMoves)
	require.NoError(t, err)
	require.False(t, gameOver)
	require.Equal(t, len(boardState.Snakes), len(nextBoardState.Snakes))

	expectedSnakes := []Snake{
		{ID: "left", Health: 0, Body: []Point{{X: 10, Y: 5}}, EliminatedCause: EliminatedByCollision, EliminatedBy: "rightEdge"},
		{ID: "rightEdge", Health: 10, Body: []Point{
			{X: 10, Y: 0},
			{X: 10, Y: 1},
			{X: 10, Y: 2},
			{X: 10, Y: 3},
			{X: 10, Y: 4},
			{X: 10, Y: 5},
		}},
	}
	for i, snake := range nextBoardState.Snakes {
		require.Equal(t, expectedSnakes[i].ID, snake.ID, snake.ID)
		require.Equal(t, expectedSnakes[i].Body, snake.Body, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedCause, snake.EliminatedCause, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedBy, snake.EliminatedBy, snake.ID)
	}
}

func TestEdgeCrossingEating(t *testing.T) {
	boardState := &BoardState{
		Width:  11,
		Height: 11,
		Snakes: []Snake{
			{ID: "left", Health: 10, Body: []Point{{X: 0, Y: 5}, {X: 1, Y: 5}}},
			{ID: "other", Health: 10, Body: []Point{{X: 5, Y: 5}}},
		},
		Food: []Point{
			{X: 10, Y: 5},
		},
	}

	snakeMoves := []SnakeMove{
		{ID: "left", Move: "left"},
		{ID: "other", Move: "left"},
	}

	r := getWrappedRuleset(Settings{})

	gameOver, nextBoardState, err := r.Execute(boardState, snakeMoves)
	require.NoError(t, err)
	require.False(t, gameOver)
	require.Equal(t, len(boardState.Snakes), len(nextBoardState.Snakes))

	expectedSnakes := []Snake{
		{ID: "left", Health: 100, Body: []Point{{X: 10, Y: 5}, {X: 0, Y: 5}, {X: 0, Y: 5}}},
		{ID: "other", Health: 9, Body: []Point{{X: 4, Y: 5}}},
	}
	for i, snake := range nextBoardState.Snakes {
		require.Equal(t, expectedSnakes[i].ID, snake.ID, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedCause, snake.EliminatedCause, snake.ID)
		require.Equal(t, expectedSnakes[i].EliminatedBy, snake.EliminatedBy, snake.ID)
		require.Equal(t, expectedSnakes[i].Body, snake.Body, snake.ID)
		require.Equal(t, expectedSnakes[i].Health, snake.Health, snake.ID)

	}
}

func TestWrap(t *testing.T) {
	// no wrap
	assert.Equal(t, 0, wrap(0, 0, 0))
	assert.Equal(t, 0, wrap(0, 1, 0))
	assert.Equal(t, 0, wrap(0, 0, 1))
	assert.Equal(t, 1, wrap(1, 0, 1))

	// wrap to min
	assert.Equal(t, 0, wrap(2, 0, 1))

	// wrap to max
	assert.Equal(t, 1, wrap(-1, 0, 1))
}

// Checks that snakes moving out of bounds get wrapped to the other side.
var wrappedCaseMoveAndWrap = gameTestCase{
	"Wrapped Case Move and Wrap",
	&BoardState{
		Width:  10,
		Height: 10,
		Snakes: []Snake{
			{
				ID:     "one",
				Body:   []Point{{X: 0, Y: 0}, {X: 1, Y: 0}},
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
				EliminatedCause: EliminatedBySelfCollision,
			},
		},
		Food:    []Point{},
		Hazards: []Point{},
	},
	[]SnakeMove{
		{ID: "one", Move: MoveLeft},
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
				Body:   []Point{{X: 9, Y: 0}, {X: 0, Y: 0}},
				Health: 99,
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
				EliminatedCause: EliminatedBySelfCollision,
			},
		},
		Food:    []Point{},
		Hazards: []Point{},
	},
}

func TestWrappedCreateNextBoardState(t *testing.T) {
	cases := []gameTestCase{
		// inherits these test cases from standard
		standardCaseErrNoMoveFound,
		standardCaseErrZeroLengthSnake,
		standardCaseMoveEatAndGrow,
		standardMoveAndCollideMAD,
		wrappedCaseMoveAndWrap,
	}
	r := getWrappedRuleset(Settings{})
	for _, gc := range cases {
		// test a RulesBuilder constructed instance
		gc.requireValidNextState(t, r)
		// also test a pipeline with the same settings
		gc.requireValidNextState(t, NewRulesetBuilder().PipelineRuleset(GameTypeWrapped, NewPipeline(wrappedRulesetStages...)))
	}
}
