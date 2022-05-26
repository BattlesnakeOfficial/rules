package maps

import (
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/stretchr/testify/require"
)

func TestDrawRing(t *testing.T) {
	ring := drawRing(11, 11, 2, 2)

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
