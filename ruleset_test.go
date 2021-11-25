package rules

import (
	"testing"

	"github.com/stretchr/testify/require"

	// included to allow using -update-fixtures for every package without errors
	_ "github.com/BattlesnakeOfficial/rules/test"
)

func TestRulesetError(t *testing.T) {
	err := (error)(RulesetError("test error string"))
	require.Equal(t, "test error string", err.Error())
}
