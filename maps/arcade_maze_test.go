package maps_test

import (
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/BattlesnakeOfficial/rules/maps"
	"github.com/stretchr/testify/require"
)

func TestArcadeMazeMap(t *testing.T) {
	tests := []struct {
		boardWidth      int
		boardHeight     int
		expectedError   error
		expectedHazards []rules.Point
	}{
		{
			boardWidth:      19,
			boardHeight:     21,
			expectedError:   nil,
			expectedHazards: maps.ArcadeMazeHazards,
		},
		{
			boardWidth:      18,
			boardHeight:     21,
			expectedError:   rules.RulesetError("This map can only be played on a 19X21 board"),
			expectedHazards: nil,
		},
		{
			boardWidth:      20,
			boardHeight:     21,
			expectedError:   rules.RulesetError("This map can only be played on a 19X21 board"),
			expectedHazards: nil,
		},
		{
			boardWidth:      19,
			boardHeight:     20,
			expectedError:   rules.RulesetError("This map can only be played on a 19X21 board"),
			expectedHazards: nil,
		},
		{
			boardWidth:      19,
			boardHeight:     22,
			expectedError:   rules.RulesetError("This map can only be played on a 19X21 board"),
			expectedHazards: nil,
		},
	}
	for _, test := range tests {
		m := maps.ArcadeMazeMap{}
		boardState := rules.NewBoardState(test.boardWidth, test.boardHeight)
		settings := rules.Settings{}
		editor := maps.NewBoardStateEditor(boardState)

		err := m.SetupBoard(boardState, settings, editor)
		if test.expectedError != nil {
			require.Equal(t, test.expectedError, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, test.expectedHazards, boardState.Hazards)

			for _, snake := range boardState.Snakes {
				require.Equal(t, 3, len(snake.Body))
			}
		}
	}
}
