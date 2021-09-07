package commands

import (
	"testing"

	"github.com/BattlesnakeOfficial/rules"
)

func TestGetIndividualBoardStateForSnake(t *testing.T) {
	s1 := rules.Snake{ID: "one", Body: []rules.Point{{X: 3, Y: 3}}}
	s2 := rules.Snake{ID: "two", Body: []rules.Point{{X: 4, Y: 3}}}
	state := &rules.BoardState{
		Height: 11,
		Width:  11,
		Snakes: []rules.Snake{s1, s2},
	}
	snake := Battlesnake{Name: "one", URL: "", ID: "one"}
	requestBody := getIndividualBoardStateForSnake(state, snake, &rules.StandardRuleset{})

	rules.RequireJSONMatchesFixture(t, "testdata/snake_request_body.json", string(requestBody))
}
