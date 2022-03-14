package rules_test

import (
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/stretchr/testify/require"
)

func TestRulesetBuilder(t *testing.T) {
	// Test that a fresh instance can produce a Ruleset
	require.NotNil(t, rules.NewRulesetBuilder().Ruleset())
	require.Equal(t, rules.GameTypeStanadard, rules.NewRulesetBuilder().Ruleset().Name(), "should default to standard game")

	// test nil safety / defaults
	require.NotNil(t, rules.NewRulesetBuilder().Ruleset())

	// make sure it works okay for lots of game types
	expectedResults := []struct {
		GameType string
		Snakes   map[string]string
	}{
		{GameType: rules.GameTypeStanadard},
		{GameType: rules.GameTypeWrapped},
		{GameType: rules.GameTypeRoyale},
		{GameType: rules.GameTypeSolo},
		{GameType: rules.GameTypeSquad, Snakes: map[string]string{
			"one":   "s1",
			"two":   "s1",
			"three": "s2",
			"four":  "s2",
			"five":  "s3",
			"six":   "s3",
			"seven": "s4",
			"eight": "s4",
		}},
		{GameType: rules.GameTypeConstrictor},
	}

	for _, expected := range expectedResults {
		t.Run(expected.GameType, func(t *testing.T) {
			rsb := rules.NewRulesetBuilder()

			rsb.WithParams(map[string]string{
				rules.ParamGameType: expected.GameType,
			})

			// add any snake squads
			for id, squad := range expected.Snakes {
				rsb = rsb.AddSnakeToSquad(id, squad)
			}

			require.NotNil(t, rsb.Ruleset())
			require.Equal(t, expected.GameType, rsb.Ruleset().Name())
		})
	}
}
