package rules

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoyaleRulesetInterface(t *testing.T) {
	var _ Ruleset = (*RoyaleRuleset)(nil)
}

func TestRoyaleDefaultSanity(t *testing.T) {
	boardState := &BoardState{}
	r := RoyaleRuleset{}
	_, err := r.CreateNextBoardState(boardState, []SnakeMove{})
	require.Error(t, err)
	require.Equal(t, errors.New("royale game must shrink at least every turn"), err)

	r = RoyaleRuleset{ShrinkEveryNTurns: 1, DamagePerTurn: 1}
	_, err = r.CreateNextBoardState(boardState, []SnakeMove{})
	require.NoError(t, err)
}

func TestRoyaleOutOfBounds(t *testing.T) {
	tests := []struct {
		Width               int32
		Height              int32
		Turn                int32
		ShrinkEveryNTurns   int32
		Error               error
		ExpectedOutOfBounds []Point
	}{
		{Error: errors.New("royale game must shrink at least every turn")},
		{ShrinkEveryNTurns: 1, ExpectedOutOfBounds: []Point{}},
		{Turn: 1, ShrinkEveryNTurns: 1, ExpectedOutOfBounds: []Point{}},
		{Width: 3, Height: 3, Turn: 1, ShrinkEveryNTurns: 10, ExpectedOutOfBounds: []Point{}},
		{Width: 3, Height: 3, Turn: 9, ShrinkEveryNTurns: 10, ExpectedOutOfBounds: []Point{}},
		{
			Width: 3, Height: 3, Turn: 10, ShrinkEveryNTurns: 10,
			ExpectedOutOfBounds: []Point{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {1, 2}, {2, 0}, {2, 1}, {2, 2}},
		},
		{
			Width: 3, Height: 3, Turn: 11, ShrinkEveryNTurns: 10,
			ExpectedOutOfBounds: []Point{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {1, 2}, {2, 0}, {2, 1}, {2, 2}},
		},
		{
			Width: 3, Height: 3, Turn: 19, ShrinkEveryNTurns: 10,
			ExpectedOutOfBounds: []Point{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {1, 2}, {2, 0}, {2, 1}, {2, 2}},
		},
		{
			Width: 3, Height: 3, Turn: 20, ShrinkEveryNTurns: 10,
			ExpectedOutOfBounds: []Point{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {1, 1}, {1, 2}, {2, 0}, {2, 1}, {2, 2}},
		},
	}

	for _, test := range tests {
		b := &BoardState{Width: test.Width, Height: test.Height}
		r := RoyaleRuleset{
			Turn:              test.Turn,
			ShrinkEveryNTurns: test.ShrinkEveryNTurns,
		}

		err := r.populateOutOfBounds(b, test.Turn)
		require.Equal(t, test.Error, err)
		if err == nil {
			// Obstacles should match
			require.Equal(t, test.ExpectedOutOfBounds, r.OutOfBounds)
			for _, expectedP := range test.ExpectedOutOfBounds {
				wasFound := false
				for _, actualP := range r.OutOfBounds {
					if expectedP == actualP {
						wasFound = true
						break
					}
				}
				require.True(t, wasFound)
			}
		}
	}
}

func TestRoyaleDamageOutOfBounds(t *testing.T) {
	tests := []struct {
		Snakes                   []Snake
		OutOfBounds              []Point
		ExpectedEliminatedCauses []string
		ExpectedEliminatedByIDs  []string
	}{
		{},
		{
			Snakes:                   []Snake{{Body: []Point{{0, 0}}}},
			OutOfBounds:              []Point{},
			ExpectedEliminatedCauses: []string{NotEliminated},
			ExpectedEliminatedByIDs:  []string{""},
		},
		{
			Snakes:                   []Snake{{Body: []Point{{0, 0}}}},
			OutOfBounds:              []Point{{0, 0}},
			ExpectedEliminatedCauses: []string{EliminatedByStarvation},
			ExpectedEliminatedByIDs:  []string{""},
		},
		{
			Snakes:                   []Snake{{Body: []Point{{0, 0}, {1, 0}, {2, 0}}}},
			OutOfBounds:              []Point{{1, 0}, {2, 0}},
			ExpectedEliminatedCauses: []string{NotEliminated},
			ExpectedEliminatedByIDs:  []string{""},
		},
		{
			Snakes: []Snake{
				{Body: []Point{{0, 0}, {1, 0}, {2, 0}}},
				{Body: []Point{{3, 3}, {3, 4}, {3, 5}, {3, 6}}},
			},
			OutOfBounds:              []Point{{1, 0}, {2, 0}, {3, 4}, {3, 5}, {3, 6}},
			ExpectedEliminatedCauses: []string{NotEliminated, NotEliminated},
			ExpectedEliminatedByIDs:  []string{"", ""},
		},
		{
			Snakes: []Snake{
				{Body: []Point{{0, 0}, {1, 0}, {2, 0}}},
				{Body: []Point{{3, 3}, {3, 4}, {3, 5}, {3, 6}}},
			},
			OutOfBounds:              []Point{{3, 3}},
			ExpectedEliminatedCauses: []string{NotEliminated, EliminatedByStarvation},
			ExpectedEliminatedByIDs:  []string{"", ""},
		},
	}

	for _, test := range tests {
		b := &BoardState{Snakes: test.Snakes}
		r := RoyaleRuleset{OutOfBounds: test.OutOfBounds, DamagePerTurn: 100}
		err := r.damageOutOfBounds(b)
		require.NoError(t, err)

		for i, snake := range b.Snakes {
			require.Equal(t, test.ExpectedEliminatedCauses[i], snake.EliminatedCause)
		}

	}
}

func TestRoyaleDamagePerTurn(t *testing.T) {
	tests := []struct {
		Health                   int32
		DamagePerTurn            int32
		ExpectedHealth           int32
		ExpectedEliminationCause string
		Error                    error
	}{
		{100, 0, 100, NotEliminated, errors.New("royale damage per turn must be greater than zero")},
		{100, -100, 100, NotEliminated, errors.New("royale damage per turn must be greater than zero")},
		{100, 1, 99, NotEliminated, nil},
		{100, 99, 1, NotEliminated, nil},
		{100, 100, 0, EliminatedByStarvation, nil},
		{100, 101, 0, EliminatedByStarvation, nil},
		{100, 999, 0, EliminatedByStarvation, nil},
		{2, 1, 1, NotEliminated, nil},
		{1, 1, 0, EliminatedByStarvation, nil},
		{1, 999, 0, EliminatedByStarvation, nil},
		{0, 1, 0, EliminatedByStarvation, nil},
		{0, 999, 0, EliminatedByStarvation, nil},
	}

	for _, test := range tests {
		b := &BoardState{Snakes: []Snake{{Health: test.Health, Body: []Point{{0, 0}}}}}
		r := RoyaleRuleset{OutOfBounds: []Point{{0, 0}}, DamagePerTurn: test.DamagePerTurn}

		err := r.damageOutOfBounds(b)
		require.Equal(t, test.Error, err)
		require.Equal(t, test.ExpectedHealth, b.Snakes[0].Health)
		require.Equal(t, test.ExpectedEliminationCause, b.Snakes[0].EliminatedCause)
	}
}

func TestRoyalDamageNextTurn(t *testing.T) {
	b := &BoardState{Width: 10, Height: 10, Snakes: []Snake{{ID: "one", Health: 100, Body: []Point{{1, 1}}}}}
	r := RoyaleRuleset{Turn: 10, ShrinkEveryNTurns: 10, DamagePerTurn: 30}
	m := []SnakeMove{{ID: "one", Move: "right"}}

	n, err := r.CreateNextBoardState(b, m)
	require.NoError(t, err)
	require.Equal(t, NotEliminated, n.Snakes[0].EliminatedCause)
	require.Equal(t, int32(99), n.Snakes[0].Health)
	require.Equal(t, Point{2, 1}, n.Snakes[0].Body[0])
	require.Equal(t, 36, len(r.OutOfBounds))

	r.Turn = 20
	n, err = r.CreateNextBoardState(b, m)
	require.NoError(t, err)
	require.Equal(t, NotEliminated, n.Snakes[0].EliminatedCause)
	require.Equal(t, int32(99), n.Snakes[0].Health)
	require.Equal(t, Point{2, 1}, n.Snakes[0].Body[0])
	require.Equal(t, 64, len(r.OutOfBounds))

	r.Turn = 21
	n, err = r.CreateNextBoardState(b, m)
	require.NoError(t, err)
	require.Equal(t, NotEliminated, n.Snakes[0].EliminatedCause)
	require.Equal(t, int32(69), n.Snakes[0].Health)
	require.Equal(t, Point{2, 1}, n.Snakes[0].Body[0])
	require.Equal(t, 64, len(r.OutOfBounds))

	b.Snakes[0].Health = 15
	n, err = r.CreateNextBoardState(b, m)
	require.NoError(t, err)
	require.Equal(t, EliminatedByStarvation, n.Snakes[0].EliminatedCause)
	require.Equal(t, int32(0), n.Snakes[0].Health)
	require.Equal(t, Point{2, 1}, n.Snakes[0].Body[0])
	require.Equal(t, 64, len(r.OutOfBounds))
}
