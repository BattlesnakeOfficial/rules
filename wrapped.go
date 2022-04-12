package rules

var wrappedRulesetStages = []string{
	StageMovementWrapBoundaries,
	StageStarvationStandard,
	StageHazardDamageStandard,
	StageFeedSnakesStandard,
	StageSpawnFoodStandard,
	StageEliminationStandard,
}

type WrappedRuleset struct {
	StandardRuleset
}

func (r *WrappedRuleset) Name() string { return GameTypeWrapped }

func (r WrappedRuleset) Execute(bs *BoardState, s Settings, sm []SnakeMove) (bool, *BoardState, error) {
	return NewPipeline(wrappedRulesetStages...).Execute(bs, s, sm)
}

func (r *WrappedRuleset) CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
	_, nextState, err := r.Execute(prevState, r.Settings(), moves)

	return nextState, err
}

func MoveSnakesWrapped(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
	if IsInitialization(b, settings, moves) {
		return false, nil
	}

	_, err := MoveSnakesStandard(b, settings, moves)
	if err != nil {
		return false, err
	}

	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause != NotEliminated {
			continue
		}
		snake.Body[0].X = wrap(snake.Body[0].X, 0, b.Width-1)
		snake.Body[0].Y = wrap(snake.Body[0].Y, 0, b.Height-1)
	}

	return false, nil
}

func wrap(value, min, max int32) int32 {
	if value < min {
		return max
	}
	if value > max {
		return min
	}
	return value
}
