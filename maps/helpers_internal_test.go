package maps

import (
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/stretchr/testify/require"
)

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
