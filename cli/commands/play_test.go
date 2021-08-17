package commands

import (
	"testing"

	"github.com/BattlesnakeOfficial/rules"
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
	snake := Battlesnake{Name: "one", URL: "", ID: "one"}
	requestBody := getIndividualBoardStateForSnake(state, snake, &rules.StandardRuleset{})

	expected := "{\"game\":{\"id\":\"\",\"timeout\":500,\"ruleset\":{\"name\":\"standard\",\"version\":\"cli\"}},\"turn\":0,\"board\":{\"height\":11,\"width\":11,\"food\":[],\"hazards\":[],\"snakes\":[{\"id\":\"one\",\"name\":\"\",\"health\":0,\"body\":[{\"x\":3,\"y\":3}],\"latency\":\"0\",\"head\":{\"x\":3,\"y\":3},\"length\":1,\"shout\":\"\",\"squad\":\"\"},{\"id\":\"two\",\"name\":\"\",\"health\":0,\"body\":[{\"x\":4,\"y\":3}],\"latency\":\"0\",\"head\":{\"x\":4,\"y\":3},\"length\":1,\"shout\":\"\",\"squad\":\"\"}]},\"you\":{\"id\":\"one\",\"name\":\"\",\"health\":0,\"body\":[{\"x\":3,\"y\":3}],\"latency\":\"0\",\"head\":{\"x\":3,\"y\":3},\"length\":1,\"shout\":\"\",\"squad\":\"\"}}"
	require.Equal(t, expected, string(requestBody))
}
