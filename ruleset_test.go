package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRulesetError(t *testing.T) {
	err := (error)(RulesetError("test error string"))
	require.Equal(t, "test error string", err.Error())
}
