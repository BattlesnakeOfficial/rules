package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type gameTestCase struct {
	prevState     *BoardState
	moves         []SnakeMove
	expectedError error
	expectedState *BoardState
}

func (gc *gameTestCase) requireCasesEqual(t *testing.T, r Ruleset) {
	prev := gc.prevState.Clone() // clone to protect against mutation (so we can ru-use test cases)
	nextState, err := r.CreateNextBoardState(prev, gc.moves)
	require.Equal(t, gc.expectedError, err)
	if gc.expectedState != nil {
		require.Equal(t, gc.expectedState.Width, nextState.Width)
		require.Equal(t, gc.expectedState.Height, nextState.Height)
		require.Equal(t, gc.expectedState.Food, nextState.Food)
		require.Equal(t, gc.expectedState.Snakes, nextState.Snakes)
		require.Equal(t, gc.expectedState.Hazards, nextState.Hazards)
	}
}
