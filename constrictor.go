package rules

var constrictorRulesetStages = []string{
	StageGameOverStandard,
	StageMovementStandard,
	StageStarvationStandard,
	StageHazardDamageStandard,
	StageFeedSnakesStandard,
	StageEliminationStandard,
	StageSpawnFoodNoFood,
	StageModifySnakesAlwaysGrow,
}

var wrappedConstrictorRulesetStages = []string{
	StageGameOverStandard,
	StageMovementWrapBoundaries,
	StageStarvationStandard,
	StageHazardDamageStandard,
	StageFeedSnakesStandard,
	StageEliminationStandard,
	StageSpawnFoodNoFood,
	StageModifySnakesAlwaysGrow,
}

func RemoveFoodConstrictor(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
	// Remove all food from the board
	b.Food = []Point{}

	return false, nil
}

func GrowSnakesConstrictor(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
	// Set all snakes to max health and ensure they grow next turn
	for i := 0; i < len(b.Snakes); i++ {
		if len(b.Snakes[i].Body) <= 0 {
			return false, ErrorZeroLengthSnake
		}
		b.Snakes[i].Health = SnakeMaxHealth

		tail := b.Snakes[i].Body[len(b.Snakes[i].Body)-1]
		subTail := b.Snakes[i].Body[len(b.Snakes[i].Body)-2]
		if tail != subTail {
			growSnake(&b.Snakes[i])
		}
	}

	return false, nil
}
