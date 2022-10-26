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

func TestRulesetBuilderInternals(t *testing.T) {
	// test Royale with seed
	rsb := NewRulesetBuilder().WithSeed(3)
	require.Equal(t, int64(3), rsb.seed)
	require.Equal(t, GameTypeRoyale, rsb.NamedRuleset(GameTypeRoyale).Name())
	require.Equal(t, int64(3), rsb.NamedRuleset(GameTypeRoyale).Settings().Seed())

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
