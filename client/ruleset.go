package client

import (
	"strconv"

	"github.com/BattlesnakeOfficial/rules"
)

const (
	SettingGameType            = "name"
	SettingFoodSpawnChance     = "foodSpawnChance"
	SettingMinimumFood         = "minimumFood"
	SettingHazardDamagePerTurn = "damagePerTurn"
	SettingShrinkEveryNTurns   = "shrinkEveryNTurns"
	SettingAllowBodyCollisions = "allowBodyCollisions"
	SettingSharedElimination   = "sharedElimination"
	SettingSharedHealth        = "sharedHealth"
	SettingSharedLength        = "sharedLength"
)

// GetRuleset constructs a ruleset from the parameters passed when creating a
// new game, and returns a ruleset customised by those parameters.
func GetRuleset(seed int64, config map[string]string, snakes []SquadSnake) rules.Ruleset {

	standardRuleset := &rules.StandardRuleset{
		FoodSpawnChance:     optionFromRulesetInt(config, SettingFoodSpawnChance, 0),
		MinimumFood:         optionFromRulesetInt(config, SettingMinimumFood, 0),
		HazardDamagePerTurn: optionFromRulesetInt(config, SettingHazardDamagePerTurn, 0),
	}

	name, ok := config[SettingGameType]
	if !ok {
		return standardRuleset
	}

	switch name {
	case rules.GameTypeConstrictor:
		return &rules.ConstrictorRuleset{
			StandardRuleset: *standardRuleset,
		}
	case rules.GameTypeRoyale:
		return &rules.RoyaleRuleset{
			StandardRuleset:   *standardRuleset,
			Seed:              seed,
			ShrinkEveryNTurns: optionFromRulesetInt(config, SettingShrinkEveryNTurns, 0),
		}
	case rules.GameTypeSolo:
		return &rules.SoloRuleset{
			StandardRuleset: *standardRuleset,
		}
	case rules.GameTypeWrapped:
		return &rules.WrappedRuleset{
			StandardRuleset: *standardRuleset,
		}
	case rules.GameTypeSquad:
		squadMap := map[string]string{}
		for _, snake := range snakes {
			squadMap[snake.GetID()] = snake.GetSquad()
		}
		return &rules.SquadRuleset{
			StandardRuleset:     *standardRuleset,
			SquadMap:            squadMap,
			AllowBodyCollisions: optionFromRulesetBool(config, SettingAllowBodyCollisions, true),
			SharedElimination:   optionFromRulesetBool(config, SettingSharedElimination, true),
			SharedHealth:        optionFromRulesetBool(config, SettingSharedHealth, true),
			SharedLength:        optionFromRulesetBool(config, SettingSharedLength, true),
		}
	}
	return standardRuleset
}

func optionFromRulesetBool(ruleset map[string]string, optionName string, defaultValue bool) bool {
	if val, ok := ruleset[optionName]; ok {
		return val == "true"
	}
	return defaultValue
}

func optionFromRulesetInt(ruleset map[string]string, optionName string, defaultValue int32) int32 {
	if val, ok := ruleset[optionName]; ok {
		i, err := strconv.Atoi(val)
		if err == nil {
			return int32(i)
		}
	}
	return defaultValue
}
