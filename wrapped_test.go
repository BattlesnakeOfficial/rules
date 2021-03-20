package rules

import (
	"testing"

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