package rules_test

import (
	"fmt"
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/stretchr/testify/require"
)

func TestRulesetBuilder(t *testing.T) {
	// Test that a fresh instance can produce a Ruleset
	require.NotNil(t, rules.NewRulesetBuilder().NamedRuleset(""))
	require.Equal(t, rules.GameTypeStandard, rules.NewRulesetBuilder().NamedRuleset("").Name(), "should default to standard game")

	// make sure it works okay for lots of game types
	expectedResults := []struct {
		GameType string
	}{
		{GameType: rules.GameTypeStandard},
		{GameType: rules.GameTypeWrapped},
		{GameType: rules.GameTypeRoyale},
		{GameType: rules.GameTypeSolo},
		{GameType: rules.GameTypeConstrictor},
		{GameType: rules.GameTypeWrappedConstrictor},
	}

	for _, expected := range expectedResults {
		t.Run(expected.GameType, func(t *testing.T) {
			rsb := rules.NewRulesetBuilder()

			rsb.WithParams(map[string]string{
				// apply the standard rule params
				rules.ParamFoodSpawnChance:     "10",
				rules.ParamMinimumFood:         "5",
				rules.ParamHazardDamagePerTurn: "12",
			})

			require.NotNil(t, rsb.NamedRuleset(expected.GameType))
			require.Equal(t, expected.GameType, rsb.NamedRuleset(expected.GameType).Name())
			// All the standard settings should always be copied over
			require.Equal(t, 10, rsb.NamedRuleset(expected.GameType).Settings().Int(rules.ParamFoodSpawnChance, 0))
			require.Equal(t, 12, rsb.NamedRuleset(expected.GameType).Settings().Int(rules.ParamHazardDamagePerTurn, 0))
			require.Equal(t, 5, rsb.NamedRuleset(expected.GameType).Settings().Int(rules.ParamMinimumFood, 0))
		})
	}
}

func TestRulesetBuilderGameOver(t *testing.T) {
	settings := rules.NewSettingsWithParams(rules.ParamShrinkEveryNTurns, "12")
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
			rsb := rules.NewRulesetBuilder().WithSettings(settings).WithSolo(test.solo)

			ruleset := rsb.NamedRuleset(test.gameType)

			gameOver, _, err := ruleset.Execute(boardState, moves)

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
	ended, err := stage(nil, rules.NewRulesetBuilder().NamedRuleset("").Settings(), nil)
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
