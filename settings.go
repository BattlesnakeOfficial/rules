package rules

// Settings contains all settings relevant to a game.
// It is used by game logic to take a previous game state and produce a next game state.
type Settings struct {
	FoodSpawnChance     int            `json:"foodSpawnChance"`
	MinimumFood         int            `json:"minimumFood"`
	HazardDamagePerTurn int            `json:"hazardDamagePerTurn"`
	HazardMap           string         `json:"hazardMap"`
	HazardMapAuthor     string         `json:"hazardMapAuthor"`
	RoyaleSettings      RoyaleSettings `json:"royale"`
	SquadSettings       SquadSettings  `json:"squad"` // Deprecated, provided with default fields for API compatibility

	rand Rand
	seed int64
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

// RoyaleSettings contains settings that are specific to the "royale" game mode
type RoyaleSettings struct {
	ShrinkEveryNTurns int `json:"shrinkEveryNTurns"`
}

// SquadSettings contains settings that are specific to the "squad" game mode
type SquadSettings struct {
	AllowBodyCollisions bool `json:"allowBodyCollisions"`
	SharedElimination   bool `json:"sharedElimination"`
	SharedHealth        bool `json:"sharedHealth"`
	SharedLength        bool `json:"sharedLength"`
}
