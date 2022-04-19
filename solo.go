package rules

var soloRulesetStages = []string{
	StageMovementStandard,
	StageStarvationStandard,
	StageHazardDamageStandard,
	StageFeedSnakesStandard,
	StageSpawnFoodStandard,
	StageEliminationStandard,
	StageGameOverSoloSnake,
}

type SoloRuleset struct {
	StandardRuleset
}

func (r *SoloRuleset) Name() string { return GameTypeSolo }

func (r SoloRuleset) Execute(bs *BoardState, s Settings, sm []SnakeMove) (bool, *BoardState, error) {
	return NewPipeline(soloRulesetStages...).Execute(bs, s, sm)
}

func (r *SoloRuleset) IsGameOver(b *BoardState) (bool, error) {
	return GameOverSolo(b, r.Settings(), nil)
}

func GameOverSolo(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
	for i := 0; i < len(b.Snakes); i++ {
		if b.Snakes[i].EliminatedCause == NotEliminated {
			return false, nil
		}
	}
	return true, nil
}
