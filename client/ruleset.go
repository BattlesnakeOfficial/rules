package client

import (
	"strconv"

	"github.com/BattlesnakeOfficial/rules"
)

type ConfigParameter string

const (
	GameType            ConfigParameter = "name"
	FoodSpawnChance     ConfigParameter = "foodSpawnChance"
	MinimumFood         ConfigParameter = "minimumFood"
	HazardDamagePerTurn ConfigParameter = "damagePerTurn"
	ShrinkEveryNTurns   ConfigParameter = "shrinkEveryNTurns"
	AllowBodyCollisions ConfigParameter = "allowBodyCollisions"
	SharedElimination   ConfigParameter = "sharedElimination"
	SharedHealth        ConfigParameter = "sharedHealth"
	SharedLength        ConfigParameter = "sharedLength"
)

// GetRuleset constructs a ruleset from the parameters passed when creating a
// new game, and returns a ruleset customised by those parameters.
func GetRuleset(seed int64, config map[ConfigParameter]string, snakes []SquadSnake) rules.Ruleset {

	standardRuleset := &rules.StandardRuleset{
		FoodSpawnChance:     optionFromRulesetInt(config, FoodSpawnChance, 0),
		MinimumFood:         optionFromRulesetInt(config, MinimumFood, 0),
		HazardDamagePerTurn: optionFromRulesetInt(config, HazardDamagePerTurn, 0),
	}

	name, ok := config[GameType]
	if !ok {
		return standardRuleset
	}

	switch name {
	case rules.Constrictor:
		return &rules.ConstrictorRuleset{
			StandardRuleset: *standardRuleset,
		}
	case rules.Royale:
		return &rules.RoyaleRuleset{
			StandardRuleset:   *standardRuleset,
			Seed:              seed,
			ShrinkEveryNTurns: optionFromRulesetInt(config, ShrinkEveryNTurns, 0),
		}
	case rules.Solo:
		return &rules.SoloRuleset{
			StandardRuleset: *standardRuleset,
		}
	case rules.Wrapped:
		return &rules.WrappedRuleset{
			StandardRuleset: *standardRuleset,
		}
	case rules.Squad:
		squadMap := map[string]string{}
		for _, snake := range snakes {
			squadMap[snake.GetID()] = snake.GetSquad()
		}
		return &rules.SquadRuleset{
			StandardRuleset:     *standardRuleset,
			SquadMap:            squadMap,
			AllowBodyCollisions: optionFromRulesetBool(config, AllowBodyCollisions, true),
			SharedElimination:   optionFromRulesetBool(config, SharedElimination, true),
			SharedHealth:        optionFromRulesetBool(config, SharedHealth, true),
			SharedLength:        optionFromRulesetBool(config, SharedLength, true),
		}
	}
	return standardRuleset
}

func optionFromRulesetBool(ruleset map[ConfigParameter]string, optionName ConfigParameter, defaultValue bool) bool {
	if val, ok := ruleset[optionName]; ok {
		return val == "true"
	}
	return defaultValue
}

func optionFromRulesetInt(ruleset map[ConfigParameter]string, optionName ConfigParameter, defaultValue int32) int32 {
	if val, ok := ruleset[optionName]; ok {
		i, err := strconv.Atoi(val)
		if err == nil {
			return int32(i)
		}
	}
	return defaultValue
}
