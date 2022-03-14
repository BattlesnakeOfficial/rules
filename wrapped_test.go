package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLeft(t *testing.T) {
	boardState := &BoardState{
		Width:  11,
		Height: 11,
		Snakes: []Snake{
			{ID: "bottomLeft", Health: 10, Body: []Point{{0, 0}}},
			{ID: "bottomRight", Health: 10, Body: []Point{{10, 0}}},
			{ID: "topLeft", Health: 10, Body: []Point{{0, 10}}},
			{ID: "topRight", Health: 10, Body: []Point{{10, 10}}},
		},
	}

	snakeMoves := []SnakeMove{
		{ID: "bottomLeft", Move: "left"},
		{ID: "bottomRight", Move: "left"},
		{ID: "topLeft", Move: "left"},
		{ID: "topRight", Move: "left"},
	}

	r := WrappedRuleset{}

	nextBoardState, err := r.CreateNextBoardState(boardState, snakeMoves)
	require.NoError(t, err)
	require.Equal(t, len(boardState.Snakes), len(nextBoardState.Snakes))

	expectedSnakes := []Snake{
		{ID: "bottomLeft", Health: 10, Body: []Point{{10, 0}}},
		{ID: "bottomRight", Health: 10, Body: []Point{{9, 0}}},
		{ID: "topLeft", Health: 10, Body: []Point{{10, 10}}},
		{ID: "topRight", Health: 10, Body: []Point{{9, 10}}},
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
			{ID: "bottomLeft", Health: 10, Body: []Point{{0, 0}}},
			{ID: "bottomRight", Health: 10, Body: []Point{{10, 0}}},
			{ID: "topLeft", Health: 10, Body: []Point{{0, 10}}},
			{ID: "topRight", Health: 10, Body: []Point{{10, 10}}},
		},
	}

	snakeMoves := []SnakeMove{
		{ID: "bottomLeft", Move: "right"},
		{ID: "bottomRight", Move: "right"},
		{ID: "topLeft", Move: "right"},
		{ID: "topRight", Move: "right"},
	}

	r := WrappedRuleset{}

	nextBoardState, err := r.CreateNextBoardState(boardState, snakeMoves)
	require.NoError(t, err)
	require.Equal(t, len(boardState.Snakes), len(nextBoardState.Snakes))

	expectedSnakes := []Snake{
		{ID: "bottomLeft", Health: 10, Body: []Point{{1, 0}}},
		{ID: "bottomRight", Health: 10, Body: []Point{{0, 0}}},
		{ID: "topLeft", Health: 10, Body: []Point{{1, 10}}},
		{ID: "topRight", Health: 10, Body: []Point{{0, 10}}},
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
			{ID: "bottomLeft", Health: 10, Body: []Point{{0, 0}}},
			{ID: "bottomRight", Health: 10, Body: []Point{{10, 0}}},
			{ID: "topLeft", Health: 10, Body: []Point{{0, 10}}},
			{ID: "topRight", Health: 10, Body: []Point{{10, 10}}},
		},
	}

	snakeMoves := []SnakeMove{
		{ID: "bottomLeft", Move: "up"},
		{ID: "bottomRight", Move: "up"},
		{ID: "topLeft", Move: "up"},
		{ID: "topRight", Move: "up"},
	}

	r := WrappedRuleset{}

	nextBoardState, err := r.CreateNextBoardState(boardState, snakeMoves)
	require.NoError(t, err)
	require.Equal(t, len(boardState.Snakes), len(nextBoardState.Snakes))

	expectedSnakes := []Snake{
		{ID: "bottomLeft", Health: 10, Body: []Point{{0, 1}}},
		{ID: "bottomRight", Health: 10, Body: []Point{{10, 1}}},
		{ID: "topLeft", Health: 10, Body: []Point{{0, 0}}},
		{ID: "topRight", Health: 10, Body: []Point{{10, 0}}},
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
			{ID: "bottomLeft", Health: 10, Body: []Point{{0, 0}}},
			{ID: "bottomRight", Health: 10, Body: []Point{{10, 0}}},
			{ID: "topLeft", Health: 10, Body: []Point{{0, 10}}},
			{ID: "topRight", Health: 10, Body: []Point{{10, 10}}},
		},
	}

	snakeMoves := []SnakeMove{
		{ID: "bottomLeft", Move: "down"},
		{ID: "bottomRight", Move: "down"},
		{ID: "topLeft", Move: "down"},
		{ID: "topRight", Move: "down"},
	}

	r := WrappedRuleset{}

	nextBoardState, err := r.CreateNextBoardState(boardState, snakeMoves)
	require.NoError(t, err)
	require.Equal(t, len(boardState.Snakes), len(nextBoardState.Snakes))

	expectedSnakes := []Snake{
		{ID: "bottomLeft", Health: 10, Body: []Point{{0, 10}}},
		{ID: "bottomRight", Health: 10, Body: []Point{{10, 10}}},
		{ID: "topLeft", Health: 10, Body: []Point{{0, 9}}},
		{ID: "topRight", Health: 10, Body: []Point{{10, 9}}},
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
			{ID: "left", Health: 10, Body: []Point{{0, 5}}},
			{ID: "rightEdge", Health: 10, Body: []Point{
				{10, 1},
				{10, 2},
				{10, 3},
				{10, 4},
				{10, 5},
				{10, 6},
			}},
		},
	}

	snakeMoves := []SnakeMove{
		{ID: "left", Move: "left"},
		{ID: "rightEdge", Move: "down"},
	}

	r := WrappedRuleset{}

	nextBoardState, err := r.CreateNextBoardState(boardState, snakeMoves)
	require.NoError(t, err)
	require.Equal(t, len(boardState.Snakes), len(nextBoardState.Snakes))

	expectedSnakes := []Snake{
		{ID: "left", Health: 0, Body: []Point{{10, 5}}, EliminatedCause: EliminatedByCollision, EliminatedBy: "rightEdge"},
		{ID: "rightEdge", Health: 10, Body: []Point{
			{10, 0},
			{10, 1},
			{10, 2},
			{10, 3},
			{10, 4},
			{10, 5},
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
			{ID: "left", Health: 10, Body: []Point{{0, 5}, {1, 5}}},
			{ID: "other", Health: 10, Body: []Point{{5, 5}}},
		},
		Food: []Point{
			{10, 5},
		},
	}

	snakeMoves := []SnakeMove{
		{ID: "left", Move: "left"},
		{ID: "other", Move: "left"},
	}

	r := WrappedRuleset{}

	nextBoardState, err := r.CreateNextBoardState(boardState, snakeMoves)
	require.NoError(t, err)
	require.Equal(t, len(boardState.Snakes), len(nextBoardState.Snakes))

	expectedSnakes := []Snake{
		{ID: "left", Health: 100, Body: []Point{{10, 5}, {0, 5}, {0, 5}}},
		{ID: "other", Health: 9, Body: []Point{{4, 5}}},
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
	assert.Equal(t, int32(0), wrap(0, 0, 0))
	assert.Equal(t, int32(0), wrap(0, 1, 0))
	assert.Equal(t, int32(0), wrap(0, 0, 1))
	assert.Equal(t, int32(1), wrap(1, 0, 1))

	// wrap to min
	assert.Equal(t, int32(0), wrap(2, 0, 1))

	// wrap to max
	assert.Equal(t, int32(1), wrap(-1, 0, 1))
}
