package commands

import (
	"fmt"
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/BattlesnakeOfficial/rules/client"
	"github.com/BattlesnakeOfficial/rules/test"
	"github.com/stretchr/testify/require"
)

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
			// royale
			rules.ParamShrinkEveryNTurns: "17",
		})

	for _, gt := range []string{
		rules.GameTypeStandard, rules.GameTypeRoyale, rules.GameTypeSolo,
		rules.GameTypeWrapped, rules.GameTypeConstrictor,
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

func TestConvertRulesSnakes(t *testing.T) {
	tests := []struct {
		name     string
		snakes   []rules.Snake
		state    map[string]SnakeState
		expected []client.Snake
	}{
		{
			name:     "empty",
			snakes:   []rules.Snake{},
			state:    map[string]SnakeState{},
			expected: []client.Snake{},
		},
		{
			name: "all properties",
			snakes: []rules.Snake{
				{ID: "one", Body: []rules.Point{{X: 3, Y: 3}, {X: 2, Y: 3}}, Health: 100},
			},
			state: map[string]SnakeState{
				"one": {
					ID:        "one",
					Name:      "ONE",
					URL:       "http://example1.com",
					Head:      "a",
					Tail:      "b",
					Color:     "#012345",
					LastMove:  "up",
					Character: '+',
				},
			},
			expected: []client.Snake{
				{
					ID:      "one",
					Name:    "ONE",
					Latency: "0",
					Health:  100,
					Body:    []client.Coord{{X: 3, Y: 3}, {X: 2, Y: 3}},
					Head:    client.Coord{X: 3, Y: 3},
					Length:  2,
					Shout:   "",
					Customizations: client.Customizations{
						Color: "#012345",
						Head:  "a",
						Tail:  "b",
					},
				},
			},
		},
		{
			name: "some eliminated",
			snakes: []rules.Snake{
				{
					ID:               "one",
					EliminatedCause:  rules.EliminatedByCollision,
					EliminatedOnTurn: 1,
					Body:             []rules.Point{{X: 3, Y: 3}},
				},
				{ID: "two", Body: []rules.Point{{X: 4, Y: 3}}},
			},
			state: map[string]SnakeState{
				"one": {ID: "one"},
				"two": {ID: "two"},
			},
			expected: []client.Snake{
				{
					ID:      "two",
					Latency: "0",
					Body:    []client.Coord{{X: 4, Y: 3}},
					Head:    client.Coord{X: 4, Y: 3},
					Length:  1,
				},
			},
		},
		{
			name: "all eliminated",
			snakes: []rules.Snake{
				{
					ID:               "one",
					EliminatedCause:  rules.EliminatedByCollision,
					EliminatedOnTurn: 1,
					Body:             []rules.Point{{X: 3, Y: 3}},
				},
			},
			state: map[string]SnakeState{
				"one": {ID: "one"},
			},
			expected: []client.Snake{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := convertRulesSnakes(test.snakes, test.state)
			require.Equal(t, test.expected, actual)
		})
	}
}
