package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBodyPassthrough(t *testing.T) {
	r := TeamRuleset{
		BodyPassthrough: true,
		teams: map[string][]string{
			"A": {"1", "2"},
			"B": {"3"},
		},
	}
	initialState := &BoardState{
		Height: 10,
		Width:  10,
		Snakes: []Snake{
			{ID: "1", Health: 100, Body: []Point{{X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}}},
			{ID: "2", Health: 100, Body: []Point{{X: 2, Y: 2}, {X: 2, Y: 3}, {X: 2, Y: 4}}},
			{ID: "3", Health: 100, Body: []Point{{X: 3, Y: 2}}},
		},
	}
	moves := []SnakeMove{
		{ID: "1", Move: "left"},
		{ID: "2", Move: "up"},
		{ID: "3", Move: "left"},
	}
	newState, err := r.ResolveMoves(initialState, moves)
	require.NoError(t, err)
	require.Equal(t, EliminatedByCollision, newState.Snakes[2].EliminatedCause)
	require.Equal(t, "2", newState.Snakes[2].EliminatedBy)

	require.Empty(t, newState.Snakes[1].EliminatedCause)
	require.Empty(t, newState.Snakes[1].EliminatedBy)

	r.BodyPassthrough = false

	newState, err = r.ResolveMoves(initialState, moves)
	require.NoError(t, err)

	require.Equal(t, EliminatedByCollision, newState.Snakes[1].EliminatedCause)
	require.Equal(t, "1", newState.Snakes[1].EliminatedBy)
}

func TestSharedStats(t *testing.T) {
	r := TeamRuleset{
		SharedStats: true,
		teams: map[string][]string{
			"A": {"1", "2"},
			"B": {"3"},
		},
	}

	initialState := &BoardState{
		Height: 10,
		Width:  10,
		Snakes: []Snake{
			{ID: "1", Health: 90, Body: []Point{{X: 1, Y: 1}}},
			{ID: "2", Health: 90, Body: []Point{{X: 2, Y: 2}}},
			{ID: "3", Health: 90, Body: []Point{{X: 3, Y: 3}}},
		},
		Food: []Point{{X: 0, Y: 1}},
	}

	moves := []SnakeMove{
		{ID: "1", Move: "left"},
		{ID: "2", Move: "left"},
		{ID: "3", Move: "left"},
	}

	newState, err := r.ResolveMoves(initialState, moves)

	require.NoError(t, err)
	require.Len(t, newState.Snakes[0].Body, 2)
	require.Equal(t, int32(100), newState.Snakes[0].Health)
	require.Len(t, newState.Snakes[1].Body, 2)
	require.Equal(t, int32(100), newState.Snakes[1].Health)
	require.Len(t, newState.Snakes[2].Body, 1)
	require.Equal(t, int32(89), newState.Snakes[2].Health)

	r.SharedStats = false

	newState, err = r.ResolveMoves(initialState, moves)

	require.NoError(t, err)
	require.Len(t, newState.Snakes[0].Body, 2)
	require.Equal(t, int32(100), newState.Snakes[0].Health)
	require.Len(t, newState.Snakes[1].Body, 1)
	require.Equal(t, int32(89), newState.Snakes[1].Health)
}

func TestSharedDeath(t *testing.T) {
	r := TeamRuleset{
		SharedDeath: true,
		teams: map[string][]string{
			"A": {"1", "2"},
		},
	}

	initialState := &BoardState{
		Height: 10,
		Width:  10,
		Snakes: []Snake{
			{ID: "1", Health: 90, Body: []Point{{X: 1, Y: 1}}},
			{ID: "2", Health: 1, Body: []Point{{X: 2, Y: 2}}},
		},
	}

	moves := []SnakeMove{
		{ID: "1", Move: "left"},
		{ID: "2", Move: "left"},
	}

	newState, err := r.ResolveMoves(initialState, moves)

	require.NoError(t, err)

	require.Equal(t, EliminatedByTeamMemberDied, newState.Snakes[0].EliminatedCause)
	require.Equal(t, "2", newState.Snakes[0].EliminatedBy)
	require.Equal(t, EliminatedByStarvation, newState.Snakes[1].EliminatedCause)

	r.SharedDeath = false

	newState, err = r.ResolveMoves(initialState, moves)

	require.NoError(t, err)
	require.Equal(t, NotEliminated, newState.Snakes[0].EliminatedCause)
	require.Equal(t, EliminatedByStarvation, newState.Snakes[1].EliminatedCause)
}
