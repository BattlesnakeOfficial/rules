package commands

import (
	"fmt"
	"os"
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/BattlesnakeOfficial/rules/test"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestPlayArgsConfigInitialised(t *testing.T) {
	oldArgs := os.Args
	oldRun := run
	defer func() {
		os.Args = oldArgs
		run = oldRun
	}()

	os.Args = []string{
		"",
		"play",
		"-g", "solo",
		"--foodSpawnChance", "2",
		"--minimumFood", "2",
		"--hazardDamagePerTurn", "2",
		"--shrinkEveryNTurns", "2",
	}
	run = func(cmd *cobra.Command, args []string) {
		// no-op
	}

	// validate initial assumptions
	require.NotEqual(t, "solo", defaultConfig[rules.ParamGameType])
	require.NotEqual(t, "2", defaultConfig[rules.ParamFoodSpawnChance])
	require.NotEqual(t, "2", defaultConfig[rules.ParamMinimumFood])
	require.NotEqual(t, "2", defaultConfig[rules.ParamHazardDamagePerTurn])
	require.NotEqual(t, "2", defaultConfig[rules.ParamShrinkEveryNTurns])

	err := playCmd.Execute()
	require.NoError(t, err)

	// check that default config is updated to reflect game type
	require.Equal(t, "solo", defaultConfig[rules.ParamGameType])
	require.Equal(t, "2", defaultConfig[rules.ParamFoodSpawnChance])
	require.Equal(t, "2", defaultConfig[rules.ParamMinimumFood])
	require.Equal(t, "2", defaultConfig[rules.ParamHazardDamagePerTurn])
	require.Equal(t, "2", defaultConfig[rules.ParamShrinkEveryNTurns])
}

func TestGetIndividualBoardStateForSnake(t *testing.T) {
	s1 := rules.Snake{ID: "one", Body: []rules.Point{{X: 3, Y: 3}}}
	s2 := rules.Snake{ID: "two", Body: []rules.Point{{X: 4, Y: 3}}}
	state := &rules.BoardState{
		Height: 11,
		Width:  11,
		Snakes: []rules.Snake{s1, s2},
	}
	s1State := SnakeState{
		ID:    "one",
		Name:  "ONE",
		URL:   "http://example1.com",
		Head:  "safe",
		Tail:  "curled",
		Color: "#123456",
	}
	s2State := SnakeState{
		ID:    "two",
		Name:  "TWO",
		URL:   "http://example2.com",
		Head:  "silly",
		Tail:  "bolt",
		Color: "#654321",
	}
	snakeStates := map[string]SnakeState{
		s1State.ID: s1State,
		s2State.ID: s2State,
	}
	initialiseGameConfig() // initialise default config
	snakeRequest := getIndividualBoardStateForSnake(state, s1State, snakeStates, getRuleset(0, snakeStates))
	requestBody := serialiseSnakeRequest(snakeRequest)

	test.RequireJSONMatchesFixture(t, "testdata/snake_request_body.json", string(requestBody))
}

func TestSettingsRequestSerialization(t *testing.T) {
	s1 := rules.Snake{ID: "one", Body: []rules.Point{{X: 3, Y: 3}}}
	s2 := rules.Snake{ID: "two", Body: []rules.Point{{X: 4, Y: 3}}}
	state := &rules.BoardState{
		Height: 11,
		Width:  11,
		Snakes: []rules.Snake{s1, s2},
	}
	s1State := SnakeState{
		ID:    "one",
		Name:  "ONE",
		URL:   "http://example1.com",
		Head:  "safe",
		Tail:  "curled",
		Color: "#123456",
	}
	s2State := SnakeState{
		ID:    "two",
		Name:  "TWO",
		URL:   "http://example2.com",
		Head:  "silly",
		Tail:  "bolt",
		Color: "#654321",
	}
	snakeStates := map[string]SnakeState{s1State.ID: s1State, s2State.ID: s2State}

	rsb := rules.NewRulesetBuilder().
		WithParams(map[string]string{
			// standard
			rules.ParamFoodSpawnChance:     "11",
			rules.ParamMinimumFood:         "7",
			rules.ParamHazardDamagePerTurn: "19",
			rules.ParamHazardMap:           "hz_spiral",
			rules.ParamHazardMapAuthor:     "altersaddle",
			// squad
			rules.ParamAllowBodyCollisions: "true",
			rules.ParamSharedElimination:   "false",
			rules.ParamSharedHealth:        "true",
			rules.ParamSharedLength:        "false",
			// royale
			rules.ParamShrinkEveryNTurns: "17",
		})

	for _, gt := range []string{
		rules.GameTypeStandard, rules.GameTypeRoyale, rules.GameTypeSolo,
		rules.GameTypeWrapped, rules.GameTypeSquad, rules.GameTypeConstrictor,
	} {
		t.Run(gt, func(t *testing.T) {
			// apply game type
			ruleset := rsb.WithParams(map[string]string{rules.ParamGameType: gt}).Ruleset()

			snakeRequest := getIndividualBoardStateForSnake(state, s1State, snakeStates, ruleset)
			requestBody := serialiseSnakeRequest(snakeRequest)
			t.Log(string(requestBody))

			test.RequireJSONMatchesFixture(t, fmt.Sprintf("testdata/snake_request_body_%s.json", gt), string(requestBody))
		})
	}
}
