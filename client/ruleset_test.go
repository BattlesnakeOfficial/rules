package client_test

import (
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/BattlesnakeOfficial/rules/client"
	"github.com/stretchr/testify/require"
)

func TestGetRuleset(t *testing.T) {
	// test nil safety / defaults
	rs := client.GetRuleset(0, nil, nil)
	require.NotNil(t, rs)

	// make sure it works okay for lots of game types
	for _, expected := range []struct {
		GameType string
		Snakes   []client.SquadSnake
	}{
		{GameType: rules.GameTypeStanadard},
		{GameType: rules.GameTypeWrapped},
		{GameType: rules.GameTypeRoyale},
		{GameType: rules.GameTypeSolo},
		{GameType: rules.GameTypeSquad, Snakes: []client.SquadSnake{
			client.Snake{ID: "one", Squad: "s1"},
			client.Snake{ID: "two", Squad: "s1"},
			client.Snake{ID: "three", Squad: "s2"},
			client.Snake{ID: "four", Squad: "s2"},
			client.Snake{ID: "five", Squad: "s3"},
			client.Snake{ID: "six", Squad: "s3"},
			client.Snake{ID: "seven", Squad: "s4"},
			client.Snake{ID: "eight", Squad: "s4"},
		}},
		{GameType: rules.GameTypeConstrictor},
	} {
		t.Run(expected.GameType, func(t *testing.T) {
			rs = client.GetRuleset(0, map[string]string{
				client.SettingGameType: expected.GameType,
			}, nil)
			require.NotNil(t, rs)
		})
	}
}
