package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type gameTestCase struct {
	name          string
	prevState     *BoardState
	moves         []SnakeMove
	expectedError error
	expectedState *BoardState
}

// requireValidNextState requires that the ruleset produces a valid next state
func (gc *gameTestCase) requireValidNextState(t *testing.T, r Ruleset) {
	t.Run(gc.name, func(t *testing.T) {
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
	})
}
