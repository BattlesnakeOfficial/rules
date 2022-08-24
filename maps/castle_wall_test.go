package maps_test

import (
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/BattlesnakeOfficial/rules/maps"
	"github.com/stretchr/testify/require"
)

func TestCastleWallMap(t *testing.T) {
	tests := []struct {
		boardWidth      int
		boardHeight     int
		expectedError   error
		expectedHazards []rules.Point
	}{
		{
			boardWidth:      11,
			boardHeight:     11,
			expectedError:   nil,
			expectedHazards: maps.CastleWallPositions.Hazards[rules.Point{X: 11, Y: 11}],
		},
		{
			boardWidth:      19,
			boardHeight:     19,
			expectedError:   nil,
			expectedHazards: maps.CastleWallPositions.Hazards[rules.Point{X: 19, Y: 19}],
		},
		{
			boardWidth:      25,
			boardHeight:     25,
			expectedError:   nil,
			expectedHazards: maps.CastleWallPositions.Hazards[rules.Point{X: 25, Y: 25}],
		},
	}
	for _, test := range tests {
		m := maps.CastleWallMap{}
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
