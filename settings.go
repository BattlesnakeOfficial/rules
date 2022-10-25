package rules

import "strconv"

// Settings contains all settings relevant to a game.
// The settings are stored as raw string values, which should not be accessed
// directly. Calling code should instead use the Int/Bool methods to parse them.
type Settings struct {
	rawValues map[string]string

	rand Rand
	seed int64
}

func NewSettings(params map[string]string) Settings {
	rawValues := make(map[string]string, len(params))

	// Copy incoming params into a new map
	for key, value := range params {
		rawValues[key] = value
	}

	return Settings{
		rawValues: rawValues,
	}
}

func NewSettingsWithParams(params ...string) Settings {
	rawValues := map[string]string{}

	for index := 1; index < len(params); index += 2 {
		rawValues[params[index-1]] = params[index]
	}

	return Settings{
		rawValues: rawValues,
	}
}

// Get a random number generator initialized based on the seed and current turn.
func (settings Settings) GetRand(turn int) Rand {
	// Allow overriding the random generator for testing
	if settings.rand != nil {
		return settings.rand
	}

	if settings.seed != 0 {
		return NewSeedRand(settings.seed + int64(turn))
	}

	// Default to global random number generator if neither seed or rand are set.
	return GlobalRand
}

func (settings Settings) WithRand(rand Rand) Settings {
	settings.rand = rand
	return settings
}

func (settings Settings) Seed() int64 {
	return settings.seed
}

func (settings Settings) WithSeed(seed int64) Settings {
	settings.seed = seed
	return settings
}

// Bool returns the boolean value for the specified parameter.
// If the parameter doesn't exist, the default value will be returned.
// If the parameter does exist, but is not "true", false will be returned.
func (settings Settings) Bool(paramName string, defaultValue bool) bool {
	if val, ok := settings.rawValues[paramName]; ok {
		return val == "true"
	}
	return defaultValue
}

// Int returns the int value for the specified parameter.
// If the parameter doesn't exist, the default value will be returned.
// If the parameter does exist, but is not a valid int, the default value will be returned.
func (settings Settings) Int(paramName string, defaultValue int) int {
	if val, ok := settings.rawValues[paramName]; ok {
		i, err := strconv.Atoi(val)
		if err == nil {
			return i
		}
	}
	return defaultValue
}
