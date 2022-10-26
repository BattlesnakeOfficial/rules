package rules

var soloRulesetStages = []string{
	StageGameOverSoloSnake,
	StageMovementStandard,
	StageStarvationStandard,
	StageHazardDamageStandard,
	StageFeedSnakesStandard,
	StageEliminationStandard,
}

func GameOverSolo(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
	for i := 0; i < len(b.Snakes); i++ {
		if b.Snakes[i].EliminatedCause == NotEliminated {
			return false, nil
		}
	}
	return true, nil
}
