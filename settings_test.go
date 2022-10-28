package rules_test

import (
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/stretchr/testify/assert"
)

func TestSettings(t *testing.T) {
	params := map[string]string{
		"invalidSetting": "abcd",
		"intSetting":     "1234",
		"boolSetting":    "true",
	}

	settings := rules.NewSettings(params)

	assert.Equal(t, 4567, settings.Int("missingIntSetting", 4567))
	assert.Equal(t, 4567, settings.Int("invalidSetting", 4567))
	assert.Equal(t, 1234, settings.Int("intSetting", 4567))

	assert.Equal(t, false, settings.Bool("missingBoolSetting", false))
	assert.Equal(t, true, settings.Bool("missingBoolSetting", true))
	assert.Equal(t, false, settings.Bool("invalidSetting", true))
	assert.Equal(t, true, settings.Bool("boolSetting", true))

	assert.Equal(t, 4567, rules.NewSettingsWithParams("newIntSetting").Int("newIntSetting", 4567))
	assert.Equal(t, 1234, rules.NewSettingsWithParams("newIntSetting", "1234").Int("newIntSetting", 4567))
	assert.Equal(t, 4567, rules.NewSettingsWithParams("x", "y", "newIntSetting").Int("newIntSetting", 4567))
}
