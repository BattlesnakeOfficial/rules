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
			boardWidth:      7,
			boardHeight:     7,
			expectedError:   nil,
			expectedHazards: maps.CastleWallPositions.Hazards[rules.Point{X: 7, Y: 7}],
		},
		{
			boardWidth:      11,
			boardHeight:     11,
			expectedError:   nil,
			expectedHazards: maps.CastleWallPositions.Hazards[rules.Point{X: 11, Y: 11}],
		},
		{
			boardWidth:      13,
			boardHeight:     13,
			expectedError:   nil,
			expectedHazards: maps.CastleWallPositions.Hazards[rules.Point{X: 13, Y: 13}],
		},
		{
			boardWidth:      15,
			boardHeight:     15,
			expectedError:   nil,
			expectedHazards: maps.CastleWallPositions.Hazards[rules.Point{X: 15, Y: 15}],
		},
		{
			boardWidth:      17,
			boardHeight:     17,
			expectedError:   nil,
			expectedHazards: maps.CastleWallPositions.Hazards[rules.Point{X: 17, Y: 17}],
		},
		{
			boardWidth:      19,
			boardHeight:     19,
			expectedError:   nil,
			expectedHazards: maps.CastleWallPositions.Hazards[rules.Point{X: 19, Y: 19}]},
		{
			boardWidth:      21,
			boardHeight:     21,
			expectedError:   nil,
			expectedHazards: maps.CastleWallPositions.Hazards[rules.Point{X: 21, Y: 21}],
		},
		{
			boardWidth:      23,
			boardHeight:     23,
			expectedError:   nil,
			expectedHazards: maps.CastleWallPositions.Hazards[rules.Point{X: 23, Y: 23}],
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
