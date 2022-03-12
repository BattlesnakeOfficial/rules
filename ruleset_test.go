package rules

import (
	"testing"

	"github.com/stretchr/testify/require"

	// included to allow using -update-fixtures for every package without errors
	_ "github.com/BattlesnakeOfficial/rules/test"
)

func TestRulesetError(t *testing.T) {
	err := (error)(RulesetError("test error string"))
	require.Equal(t, "test error string", err.Error())
}

func TestRulesBuilder(t *testing.T) {
	// test nil safety / defaults
	require.NotNil(t, NewBuilder().Ruleset())

	// test seed
	require.Equal(t, int64(3), NewBuilder().WithSeed(3).seed)

	// make sure it works okay for lots of game types
	expectedResults := []struct {
		GameType string
		Snakes   map[string]string
	}{
		{GameType: GameTypeStanadard},
		{GameType: GameTypeWrapped},
		{GameType: GameTypeRoyale},
		{GameType: GameTypeSolo},
		{GameType: GameTypeSquad, Snakes: map[string]string{
			"one":   "s1",
			"two":   "s1",
			"three": "s2",
			"four":  "s2",
			"five":  "s3",
			"six":   "s3",
			"seven": "s4",
			"eight": "s4",
		}},
		{GameType: GameTypeConstrictor},
	}

	for _, expected := range expectedResults {
		t.Run(expected.GameType, func(t *testing.T) {
			rsb := NewBuilder()

			rsb.WithParams(map[string]string{
				ParamGameType: expected.GameType,
			})

			// add any snake squads
			for id, squad := range expected.Snakes {
				rsb = rsb.AddSnakeToSquad(id, squad)
			}

			require.NotNil(t, rsb.Ruleset())
			require.Equal(t, expected.GameType, rsb.Ruleset().Name())
			require.Equal(t, len(expected.Snakes), len(rsb.squads))
		})
	}

}
