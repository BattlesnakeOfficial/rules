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

func TestRulesetBuilderInternals(t *testing.T) {

	// test Standard and Royale with seed
	rsb := NewRulesetBuilder().WithSeed(3).WithParams(map[string]string{ParamGameType: GameTypeRoyale})
	require.Equal(t, int64(3), rsb.seed)
	require.Equal(t, GameTypeRoyale, rsb.Ruleset().Name())
	require.Equal(t, int64(0), rsb.Ruleset().(*RoyaleRuleset).Seed)
	require.Equal(t, int64(3), rsb.Ruleset().(*RoyaleRuleset).StandardRuleset.seed)

	// test squad configuration
	rsb = NewRulesetBuilder().
		WithParams(map[string]string{
			ParamGameType: GameTypeSquad,
		}).
		AddSnakeToSquad("snek1", "squad1").
		AddSnakeToSquad("snek2", "squad1").
		AddSnakeToSquad("snek3", "squad2").
		AddSnakeToSquad("snek4", "squad2")

	require.NotNil(t, rsb.Ruleset())
	require.Equal(t, GameTypeSquad, rsb.Ruleset().Name())
	require.Equal(t, 4, len(rsb.squads))
	require.Equal(t, "squad1", rsb.Ruleset().(*SquadRuleset).SquadMap["snek1"])
	require.Equal(t, "squad1", rsb.Ruleset().(*SquadRuleset).SquadMap["snek2"])
	require.Equal(t, "squad2", rsb.Ruleset().(*SquadRuleset).SquadMap["snek3"])
	require.Equal(t, "squad2", rsb.Ruleset().(*SquadRuleset).SquadMap["snek4"])

	// test parameter merging
	rsb = NewRulesetBuilder().
		WithParams(map[string]string{
			"someSetting":    "some value",
			"anotherSetting": "another value",
		}).
		WithParams(map[string]string{
			"anotherSetting": "overridden value",
			"aNewSetting":    "a new value",
		})

	require.Equal(t, map[string]string{
		"someSetting":    "some value",
		"anotherSetting": "overridden value",
		"aNewSetting":    "a new value",
	}, rsb.params, "multiple calls to WithParams should merge parameters")
}
