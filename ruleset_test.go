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
		HazardMap:           "hz_spiral",
		HazardMapAuthor:     "altersaddle",
	}
	assert.Equal(t, ruleset.MinimumFood, ruleset.Settings().MinimumFood)
	assert.Equal(t, ruleset.FoodSpawnChance, ruleset.Settings().FoodSpawnChance)
	assert.Equal(t, ruleset.HazardDamagePerTurn, ruleset.Settings().HazardDamagePerTurn)
	assert.Equal(t, ruleset.HazardMap, ruleset.Settings().HazardMap)
	assert.Equal(t, ruleset.HazardMapAuthor, ruleset.Settings().HazardMapAuthor)
}

func TestWrappedRulesetSettings(t *testing.T) {
	ruleset := rules.WrappedRuleset{
		StandardRuleset: rules.StandardRuleset{
			MinimumFood:         5,
			FoodSpawnChance:     10,
			HazardDamagePerTurn: 10,
			HazardMap:           "hz_spiral",
			HazardMapAuthor:     "altersaddle",
		},
	}
	assert.Equal(t, ruleset.MinimumFood, ruleset.Settings().MinimumFood)
	assert.Equal(t, ruleset.FoodSpawnChance, ruleset.Settings().FoodSpawnChance)
	assert.Equal(t, ruleset.HazardDamagePerTurn, ruleset.Settings().HazardDamagePerTurn)
	assert.Equal(t, ruleset.HazardMap, ruleset.Settings().HazardMap)
	assert.Equal(t, ruleset.HazardMapAuthor, ruleset.Settings().HazardMapAuthor)
}

func TestSoloRulesetSettings(t *testing.T) {
	ruleset := rules.SoloRuleset{
		StandardRuleset: rules.StandardRuleset{
			MinimumFood:         5,
			FoodSpawnChance:     10,
			HazardDamagePerTurn: 10,
			HazardMap:           "hz_spiral",
			HazardMapAuthor:     "altersaddle",
		},
	}
	assert.Equal(t, ruleset.MinimumFood, ruleset.Settings().MinimumFood)
	assert.Equal(t, ruleset.FoodSpawnChance, ruleset.Settings().FoodSpawnChance)
	assert.Equal(t, ruleset.HazardDamagePerTurn, ruleset.Settings().HazardDamagePerTurn)
	assert.Equal(t, ruleset.HazardMap, ruleset.Settings().HazardMap)
	assert.Equal(t, ruleset.HazardMapAuthor, ruleset.Settings().HazardMapAuthor)
}

func TestRoyaleRulesetSettings(t *testing.T) {
	ruleset := rules.RoyaleRuleset{
		Seed:              30,
		ShrinkEveryNTurns: 12,
		StandardRuleset: rules.StandardRuleset{
			MinimumFood:         5,
			FoodSpawnChance:     10,
			HazardDamagePerTurn: 10,
			HazardMap:           "hz_spiral",
			HazardMapAuthor:     "altersaddle",
		},
	}
	assert.Equal(t, ruleset.ShrinkEveryNTurns, ruleset.Settings().RoyaleSettings.ShrinkEveryNTurns)
	assert.Equal(t, ruleset.MinimumFood, ruleset.Settings().MinimumFood)
	assert.Equal(t, ruleset.FoodSpawnChance, ruleset.Settings().FoodSpawnChance)
	assert.Equal(t, ruleset.HazardDamagePerTurn, ruleset.Settings().HazardDamagePerTurn)
	assert.Equal(t, ruleset.HazardMap, ruleset.Settings().HazardMap)
	assert.Equal(t, ruleset.HazardMapAuthor, ruleset.Settings().HazardMapAuthor)
}

func TestConstrictorRulesetSettings(t *testing.T) {
	ruleset := rules.ConstrictorRuleset{
		StandardRuleset: rules.StandardRuleset{
			MinimumFood:         5,
			FoodSpawnChance:     10,
			HazardDamagePerTurn: 10,
			HazardMap:           "hz_spiral",
			HazardMapAuthor:     "altersaddle",
		},
	}
	assert.Equal(t, ruleset.MinimumFood, ruleset.Settings().MinimumFood)
	assert.Equal(t, ruleset.FoodSpawnChance, ruleset.Settings().FoodSpawnChance)
	assert.Equal(t, ruleset.HazardDamagePerTurn, ruleset.Settings().HazardDamagePerTurn)
	assert.Equal(t, ruleset.HazardMap, ruleset.Settings().HazardMap)
	assert.Equal(t, ruleset.HazardMapAuthor, ruleset.Settings().HazardMapAuthor)
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
			HazardMap:           "hz_spiral",
			HazardMapAuthor:     "altersaddle",
		},
	}
	assert.Equal(t, ruleset.AllowBodyCollisions, ruleset.Settings().SquadSettings.AllowBodyCollisions)
	assert.Equal(t, ruleset.SharedElimination, ruleset.Settings().SquadSettings.SharedElimination)
	assert.Equal(t, ruleset.SharedHealth, ruleset.Settings().SquadSettings.SharedHealth)
	assert.Equal(t, ruleset.SharedLength, ruleset.Settings().SquadSettings.SharedLength)

	assert.Equal(t, ruleset.MinimumFood, ruleset.Settings().MinimumFood)
	assert.Equal(t, ruleset.FoodSpawnChance, ruleset.Settings().FoodSpawnChance)
	assert.Equal(t, ruleset.HazardDamagePerTurn, ruleset.Settings().HazardDamagePerTurn)
	assert.Equal(t, ruleset.HazardMap, ruleset.Settings().HazardMap)
	assert.Equal(t, ruleset.HazardMapAuthor, ruleset.Settings().HazardMapAuthor)
}

func TestRulesetBuilder(t *testing.T) {
	// Test that a fresh instance can produce a Ruleset
	require.NotNil(t, rules.NewRulesetBuilder().Ruleset())
	require.Equal(t, rules.GameTypeStandard, rules.NewRulesetBuilder().Ruleset().Name(), "should default to standard game")

	// test nil safety / defaults
	require.NotNil(t, rules.NewRulesetBuilder().Ruleset())

	// make sure it works okay for lots of game types
	expectedResults := []struct {
		GameType string
		Snakes   map[string]string
	}{
		{GameType: rules.GameTypeStandard},
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
				// apply the standard rule params
				rules.ParamGameType:            expected.GameType,
				rules.ParamFoodSpawnChance:     "10",
				rules.ParamMinimumFood:         "5",
				rules.ParamHazardDamagePerTurn: "12",
				rules.ParamHazardMap:           "test",
				rules.ParamHazardMapAuthor:     "tester",
			})

			// add any snake squads
			for id, squad := range expected.Snakes {
				rsb = rsb.AddSnakeToSquad(id, squad)
			}

			require.NotNil(t, rsb.Ruleset())
			require.Equal(t, expected.GameType, rsb.Ruleset().Name())
			// All the standard settings should always be copied over
			require.Equal(t, int32(10), rsb.Ruleset().Settings().FoodSpawnChance)
			require.Equal(t, int32(12), rsb.Ruleset().Settings().HazardDamagePerTurn)
			require.Equal(t, int32(5), rsb.Ruleset().Settings().MinimumFood)
			require.Equal(t, "test", rsb.Ruleset().Settings().HazardMap)
			require.Equal(t, "tester", rsb.Ruleset().Settings().HazardMapAuthor)
		})
	}
}

func TestStageFuncContract(t *testing.T) {
	//nolint:gosimple
	var stage rules.StageFunc
	stage = func(bs *rules.BoardState, s rules.Settings, sm []rules.SnakeMove) (bool, error) {
		return true, nil
	}
	ended, err := stage(nil, rules.NewRulesetBuilder().Ruleset().Settings(), nil)
	require.NoError(t, err)
	require.True(t, ended)
}
