package rules

import (
	"testing"

	"github.com/stretchr/testify/require"

	// included to allow using -update-fixtures for every package without errors
	_ "github.com/BattlesnakeOfficial/rules/test"
)

func TestParamInt32(t *testing.T) {
	require.Equal(t, int32(5), paramsInt32(nil, "test", 5), "nil map")
	require.Equal(t, int32(10), paramsInt32(map[string]string{}, "foo", 10), "empty map")
	require.Equal(t, int32(10), paramsInt32(map[string]string{"hullo": "there"}, "hullo", 10), "invalid value")
	require.Equal(t, int32(20), paramsInt32(map[string]string{"bonjour": "20"}, "bonjour", 20), "valid value")
}

func TestParamBool(t *testing.T) {
	// missing values default to specified value
	require.Equal(t, true, paramsBool(nil, "test", true), "nil map true")
	require.Equal(t, false, paramsBool(nil, "test", false), "nil map false")

	// missing values default to specified value
	require.Equal(t, true, paramsBool(map[string]string{}, "foo", true), "empty map true")
	require.Equal(t, false, paramsBool(map[string]string{}, "foo", false), "empty map false")

	// invalid values (exist but not booL) default to false
	require.Equal(t, false, paramsBool(map[string]string{"hullo": "there"}, "hullo", true), "invalid value default true")
	require.Equal(t, false, paramsBool(map[string]string{"hullo": "there"}, "hullo", false), "invalid value default false")

	// valid values ignore defaults
	require.Equal(t, false, paramsBool(map[string]string{"bonjour": "false"}, "bonjour", false), "valid value false")
	require.Equal(t, true, paramsBool(map[string]string{"bonjour": "true"}, "bonjour", false), "valid value true")
}

func TestRulesetError(t *testing.T) {
	err := (error)(RulesetError("test error string"))
	require.Equal(t, "test error string", err.Error())
}

func TestRulesetBuilder(t *testing.T) {
	// test nil safety / defaults
	require.NotNil(t, NewRulesetBuilder().Ruleset())

	// test seed
	require.Equal(t, int64(3), NewRulesetBuilder().WithSeed(3).seed)

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
			rsb := NewRulesetBuilder()

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
