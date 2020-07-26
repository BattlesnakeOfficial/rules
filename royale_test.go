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
	require.Equal(t, err, errors.New("royale game must shrink at least every 1 turn"))

	r = RoyaleRuleset{ShrinkEveryNTurns: 1}
	_, err = r.CreateNextBoardState(boardState, []SnakeMove{})
	require.NoError(t, err)
}

func TestRoyalePopulateObstacles(t *testing.T) {
	tests := []struct {
		Width               int32
		Height              int32
		Turn                int32
		ShrinkEveryNTurns   int32
		Error               error
		ExpectedOutOfBounds []Point
	}{
		{Error: errors.New("royale game must shrink at least every 1 turn")},
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

		err := r.populateOutOfBounds(b)
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

func TestRoyaleEliminateOutOfBounds(t *testing.T) {
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
			ExpectedEliminatedCauses: []string{EliminatedByOutOfBounds},
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
			ExpectedEliminatedCauses: []string{NotEliminated, EliminatedByOutOfBounds},
			ExpectedEliminatedByIDs:  []string{"", ""},
		},
	}

	for _, test := range tests {
		b := &BoardState{Snakes: test.Snakes}
		r := RoyaleRuleset{OutOfBounds: test.OutOfBounds}
		err := r.eliminateOutOfBounds(b)
		require.NoError(t, err)

		for i, snake := range b.Snakes {
			require.Equal(t, test.ExpectedEliminatedCauses[i], snake.EliminatedCause)
		}

	}
}
