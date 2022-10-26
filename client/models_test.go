package client

import (
	"encoding/json"
	"testing"

	"github.com/BattlesnakeOfficial/rules/test"
	"github.com/stretchr/testify/require"
)

func TestBuildSnakeRequestJSON(t *testing.T) {
	snakeRequest := exampleSnakeRequest()
	data, err := json.MarshalIndent(snakeRequest, "", "  ")
	require.NoError(t, err)

	test.RequireJSONMatchesFixture(t, "testdata/snake_request.json", string(data))
}

func TestBuildSnakeRequestJSONEmptyRulesetSettings(t *testing.T) {
	snakeRequest := exampleSnakeRequest()
	snakeRequest.Game.Ruleset.Settings = RulesetSettings{}
	data, err := json.MarshalIndent(snakeRequest, "", "  ")
	require.NoError(t, err)

	test.RequireJSONMatchesFixture(t, "testdata/snake_request_empty_ruleset_settings.json", string(data))
}
