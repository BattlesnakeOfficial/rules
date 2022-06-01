package maps

import (
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/stretchr/testify/require"
)

func TestMaxInt(t *testing.T) {
	// simple case
	n := maxInt(0, 1, 2, 3)
	require.Equal(t, 3, n)

	// use negative, out of order
	n = maxInt(0, -1, 200, 3)
	require.Equal(t, 200, n)

	// use only 1 value, and negative
	n = maxInt(-99)
	require.Equal(t, -99, n)

	// use duplicate values
	n = maxInt(3, 3, 3)
	require.Equal(t, 3, n)

	// use duplicate and other values
	n = maxInt(-1, 3, 5, 3, 3, 2)
	require.Equal(t, 5, n)
}

func TestIsOnBoard(t *testing.T) {
	// a few spot checks
	require.True(t, isOnBoard(11, 11, 0, 0))
	require.False(t, isOnBoard(11, 11, -1, 0))
	require.True(t, isOnBoard(11, 11, 10, 10))
	require.False(t, isOnBoard(11, 11, 11, 11))
	require.True(t, isOnBoard(2, 2, 1, 1))

	// exhaustive check on a small, non-square board
	for x := 0; x < 4; x++ {
		for y := 0; y < 9; y++ {
			require.True(t, isOnBoard(4, 9, x, y))
		}
	}
}

func TestDrawRing(t *testing.T) {
	_, err := drawRing(0, 11, 2, 2)
	require.Equal(t, "board width too small", err.Error())

	_, err = drawRing(11, 0, 2, 2)
	require.Equal(t, "board height too small", err.Error())

	_, err = drawRing(11, 11, 10, 2)
	require.Equal(t, "horizontal offset too large", err.Error())

	_, err = drawRing(11, 11, 2, 10)
	require.Equal(t, "vertical offset too large", err.Error())

	_, err = drawRing(11, 11, 0, 2)
	require.Equal(t, "horizontal offset too small", err.Error())

	_, err = drawRing(11, 11, 2, 0)
	require.Equal(t, "vertical offset too small", err.Error())

	_, err = drawRing(19, 1, 4, 4)
	require.Equal(t, "vertical offset too large", err.Error())

	_, err = drawRing(19, 1, 6, 6)
	require.Equal(t, "vertical offset too large", err.Error())

	_, err = drawRing(14, 7, 6, 6)
	require.Equal(t, "vertical offset too large", err.Error())

	_, err = drawRing(18, 10, 8, 8)
	require.NoError(t, err)

	ring, err := drawRing(11, 11, 2, 2)
	require.NoError(t, err)

	// ring should not be empty
	require.NotEmpty(t, ring)

	// should have exactly 32 points in this ring
	require.Len(t, ring, 32)

	// ensure no duplicates
	seen := map[rules.Point]struct{}{}
	for _, p := range ring {
		// _, ok := seen[p]
		require.NotContains(t, seen, p)
		seen[p] = struct{}{}
	}

	// spot check a few known points
	require.Contains(t, seen, rules.Point{X: 1, Y: 1}, "bottom left")
	require.Contains(t, seen, rules.Point{X: 1, Y: 9}, "top left")
	require.Contains(t, seen, rules.Point{X: 9, Y: 1}, "bottom right")
	require.Contains(t, seen, rules.Point{X: 9, Y: 9}, "top right")
	require.Contains(t, seen, rules.Point{X: 1, Y: 5})
	require.Contains(t, seen, rules.Point{X: 6, Y: 1})
	require.Contains(t, seen, rules.Point{X: 8, Y: 9})
}
