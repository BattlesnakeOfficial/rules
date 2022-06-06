package rules_test

import (
	"fmt"
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

func TestRulesetBuilder(t *testing.T) {
	// Test that a fresh instance can produce a Ruleset
	require.NotNil(t, rules.NewRulesetBuilder().Ruleset())
	require.Equal(t, rules.GameTypeStandard, rules.NewRulesetBuilder().Ruleset().Name(), "should default to standard game")

	// test nil safety / defaults
	require.NotNil(t, rules.NewRulesetBuilder().Ruleset())

	// make sure it works okay for lots of game types
	expectedResults := []struct {
		GameType string
	}{
		{GameType: rules.GameTypeStandard},
		{GameType: rules.GameTypeWrapped},
		{GameType: rules.GameTypeRoyale},
		{GameType: rules.GameTypeSolo},
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

			require.NotNil(t, rsb.Ruleset())
			require.Equal(t, expected.GameType, rsb.Ruleset().Name())
			// All the standard settings should always be copied over
			require.Equal(t, 10, rsb.Ruleset().Settings().FoodSpawnChance)
			require.Equal(t, 12, rsb.Ruleset().Settings().HazardDamagePerTurn)
			require.Equal(t, 5, rsb.Ruleset().Settings().MinimumFood)
			require.Equal(t, "test", rsb.Ruleset().Settings().HazardMap)
			require.Equal(t, "tester", rsb.Ruleset().Settings().HazardMapAuthor)
		})
	}
}

func TestRulesetBuilderGameOver(t *testing.T) {
	settings := rules.Settings{
		RoyaleSettings: rules.RoyaleSettings{
			ShrinkEveryNTurns: 12,
		},
	}
	moves := []rules.SnakeMove{
		{ID: "1", Move: "up"},
	}
	boardState := rules.NewBoardState(7, 7)
	boardState.Snakes = append(boardState.Snakes, rules.Snake{
		ID: "1",
		Body: []rules.Point{
			{X: 3, Y: 3},
			{X: 3, Y: 3},
			{X: 3, Y: 3},
		},
		Health: 100,
	})

	tests := []struct {
		gameType string
		solo     bool
		gameOver bool
	}{
		{
			gameType: rules.GameTypeStandard,
			solo:     false,
			gameOver: true,
		},
		{
			gameType: rules.GameTypeConstrictor,
			solo:     false,
			gameOver: true,
		},
		{
			gameType: rules.GameTypeRoyale,
			solo:     false,
			gameOver: true,
		},
		{
			gameType: rules.GameTypeWrapped,
			solo:     false,
			gameOver: true,
		},
		{
			gameType: rules.GameTypeSolo,
			solo:     false,
			gameOver: false,
		},
		{
			gameType: rules.GameTypeStandard,
			solo:     true,
			gameOver: false,
		},
		{
			gameType: rules.GameTypeConstrictor,
			solo:     true,
			gameOver: false,
		},
		{
			gameType: rules.GameTypeRoyale,
			solo:     true,
			gameOver: false,
		},
		{
			gameType: rules.GameTypeWrapped,
			solo:     true,
			gameOver: false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v_%v", test.gameType, test.solo), func(t *testing.T) {
			rsb := rules.NewRulesetBuilder().WithParams(map[string]string{
				rules.ParamGameType: test.gameType,
			}).WithSolo(test.solo)

			ruleset := rsb.Ruleset()

			gameOver, _, err := ruleset.Execute(boardState, settings, moves)

			require.NoError(t, err)
			require.Equal(t, test.gameOver, gameOver)
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

func TestRulesetBuilderGetRand(t *testing.T) {
	var seed int64 = 12345
	var turn int = 5
	ruleset := rules.NewRulesetBuilder().WithSeed(seed).PipelineRuleset("example", rules.NewPipeline(rules.StageGameOverStandard))

	rand1 := ruleset.Settings().GetRand(turn)

	// Should produce a predictable series of numbers based on a seed
	require.Equal(t, 83, rand1.Intn(100))
	require.Equal(t, 15, rand1.Intn(100))

	// Should produce the same number if re-initialized
	require.Equal(
		t,
		ruleset.Settings().GetRand(turn).Intn(100),
		ruleset.Settings().GetRand(turn).Intn(100),
	)

	// Should produce a different series of numbers for another turn
	require.Equal(t, 69, rand1.Intn(100))
	require.Equal(t, 86, rand1.Intn(100))
}
