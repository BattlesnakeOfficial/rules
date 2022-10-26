package client

import "github.com/BattlesnakeOfficial/rules"

func exampleSnakeRequest() SnakeRequest {
	return SnakeRequest{
		Game: Game{
			ID: "game-id",
			Ruleset: Ruleset{
				Name:     "test-ruleset-name",
				Version:  "cli",
				Settings: ConvertRulesetSettings(exampleRulesetSettings),
			},
			Timeout: 33,
			Source:  "league",
			Map:     "standard",
		},
		Turn: 11,
		Board: Board{
			Height: 22,
			Width:  11,
			Snakes: []Snake{
				{
					ID:      "snake-0",
					Name:    "snake-0-name",
					Latency: "snake-0-latency",
					Health:  100,
					Body:    []Coord{{X: 1, Y: 2}, {X: 1, Y: 3}, {X: 1, Y: 4}},
					Head:    Coord{X: 1, Y: 2},
					Length:  3,
					Shout:   "snake-0-shout",
					Squad:   "",
					Customizations: Customizations{
						Head:  "safe",
						Tail:  "curled",
						Color: "#123456",
					},
				},
				{
					ID:      "snake-1",
					Name:    "snake-1-name",
					Latency: "snake-1-latency",
					Health:  200,
					Body:    []Coord{{X: 2, Y: 2}, {X: 2, Y: 3}, {X: 2, Y: 4}},
					Head:    Coord{X: 2, Y: 2},
					Length:  3,
					Shout:   "snake-1-shout",
					Squad:   "snake-1-squad",
					Customizations: Customizations{
						Head:  "silly",
						Tail:  "bolt",
						Color: "#654321",
					},
				},
			},
			Food:    []Coord{{X: 2, Y: 2}},
			Hazards: []Coord{{X: 8, Y: 8}, {X: 9, Y: 9}},
		},
		You: Snake{
			ID:      "snake-1",
			Name:    "snake-1-name",
			Latency: "snake-1-latency",
			Health:  200,
			Body:    []Coord{{X: 2, Y: 2}, {X: 2, Y: 3}, {X: 2, Y: 4}},
			Head:    Coord{X: 2, Y: 2},
			Length:  3,
			Shout:   "snake-1-shout",
			Squad:   "snake-1-squad",
			Customizations: Customizations{
				Head:  "silly",
				Tail:  "bolt",
				Color: "#654321",
			},
		},
	}
}

var exampleRulesetSettings = rules.NewSettings(map[string]string{
	rules.ParamFoodSpawnChance:     "10",
	rules.ParamMinimumFood:         "20",
	rules.ParamHazardDamagePerTurn: "30",
	rules.ParamShrinkEveryNTurns:   "40",
})
