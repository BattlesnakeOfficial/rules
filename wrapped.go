package rules

type WrappedRuleset struct {
	StandardRuleset
}

func (r *WrappedRuleset) Name() string { return GameTypeWrapped }

func (r WrappedRuleset) Pipeline() (*Pipeline, error) {
	return NewPipeline(
		"movement.wrapped",
		"reducehealth.standard",
		"hazarddamage.standard",
		"eatfood.standard",
		"placefood.standard",
		"eliminatesnake.standard",
	)
}

func (r *WrappedRuleset) CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
	nextState := prevState.Clone()

	p, err := r.Pipeline()
	if err != nil {
		return nil, err
	}
	_, err = p.Execute(nextState, r.Settings(), moves)

	return nextState, err
}

func (r *WrappedRuleset) moveSnakes(b *BoardState, moves []SnakeMove) error {
	_, err := r.callStageFunc(MoveSnakesWrapped, b, moves)
	return err
}

func MoveSnakesWrapped(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
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
