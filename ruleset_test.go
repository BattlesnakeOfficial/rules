package rules_test

import (
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStandardRulesetSettings(t *testing.T) {
	ruleset := rules.StandardRuleset{
		MinimumFood:         5,
		FoodSpawnChance:     10,
		HazardDamagePerTurn: 10,
	}
	assert.Equal(t, ruleset.MinimumFood, ruleset.Settings().MinimumFood)
	assert.Equal(t, ruleset.FoodSpawnChance, ruleset.Settings().FoodSpawnChance)
	assert.Equal(t, ruleset.HazardDamagePerTurn, ruleset.Settings().HazardDamagePerTurn)
}

func TestWrappedRulesetSettings(t *testing.T) {
	ruleset := rules.WrappedRuleset{
		StandardRuleset: rules.StandardRuleset{
			MinimumFood:         5,
			FoodSpawnChance:     10,
			HazardDamagePerTurn: 10,
		},
	}
	assert.Equal(t, ruleset.MinimumFood, ruleset.Settings().MinimumFood)
	assert.Equal(t, ruleset.FoodSpawnChance, ruleset.Settings().FoodSpawnChance)
	assert.Equal(t, ruleset.HazardDamagePerTurn, ruleset.Settings().HazardDamagePerTurn)
}

func TestSoloRulesetSettings(t *testing.T) {
	ruleset := rules.SoloRuleset{
		StandardRuleset: rules.StandardRuleset{
			MinimumFood:         5,
			FoodSpawnChance:     10,
			HazardDamagePerTurn: 10,
		},
	}
	assert.Equal(t, ruleset.MinimumFood, ruleset.Settings().MinimumFood)
	assert.Equal(t, ruleset.FoodSpawnChance, ruleset.Settings().FoodSpawnChance)
	assert.Equal(t, ruleset.HazardDamagePerTurn, ruleset.Settings().HazardDamagePerTurn)
}

func TestRoyaleRulesetSettings(t *testing.T) {
	ruleset := rules.RoyaleRuleset{
		Seed:              30,
		ShrinkEveryNTurns: 12,
		StandardRuleset: rules.StandardRuleset{
			MinimumFood:         5,
			FoodSpawnChance:     10,
			HazardDamagePerTurn: 10,
		},
	}
	assert.Equal(t, ruleset.ShrinkEveryNTurns, ruleset.Settings().RoyaleSettings.ShrinkEveryNTurns)
	assert.Equal(t, ruleset.MinimumFood, ruleset.Settings().MinimumFood)
	assert.Equal(t, ruleset.FoodSpawnChance, ruleset.Settings().FoodSpawnChance)
	assert.Equal(t, ruleset.HazardDamagePerTurn, ruleset.Settings().HazardDamagePerTurn)
}

func TestConstrictorRulesetSettings(t *testing.T) {
	ruleset := rules.ConstrictorRuleset{
		StandardRuleset: rules.StandardRuleset{
			MinimumFood:         5,
			FoodSpawnChance:     10,
			HazardDamagePerTurn: 10,
		},
	}
	assert.Equal(t, ruleset.MinimumFood, ruleset.Settings().MinimumFood)
	assert.Equal(t, ruleset.FoodSpawnChance, ruleset.Settings().FoodSpawnChance)
	assert.Equal(t, ruleset.HazardDamagePerTurn, ruleset.Settings().HazardDamagePerTurn)
}

func TestSquadRulesetSettings(t *testing.T) {
	ruleset := rules.SquadRuleset{
		AllowBodyCollisions: true,
		SharedElimination:   false,
		SharedHealth:        true,
		SharedLength:        false,
		StandardRuleset: rules.StandardRuleset{
			MinimumFood:         5,
			FoodSpawnChance:     10,
			HazardDamagePerTurn: 10,
		},
	}
	assert.Equal(t, ruleset.AllowBodyCollisions, ruleset.Settings().SquadSettings.AllowBodyCollisions)
	assert.Equal(t, ruleset.SharedElimination, ruleset.Settings().SquadSettings.SharedElimination)
	assert.Equal(t, ruleset.SharedHealth, ruleset.Settings().SquadSettings.SharedHealth)
	assert.Equal(t, ruleset.SharedLength, ruleset.Settings().SquadSettings.SharedLength)

	assert.Equal(t, ruleset.MinimumFood, ruleset.Settings().MinimumFood)
	assert.Equal(t, ruleset.FoodSpawnChance, ruleset.Settings().FoodSpawnChance)
	assert.Equal(t, ruleset.HazardDamagePerTurn, ruleset.Settings().HazardDamagePerTurn)
}

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
