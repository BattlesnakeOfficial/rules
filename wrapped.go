package rules

type WrappedRuleset struct {
	StandardRuleset
}

func replace(value, min, max int32) int32 {
	if value < min {
		return max
	}
	if value > max {
		return min
	}
	return value
}

func (r *WrappedRuleset) CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
	nextState := prevState.Copy()

	err := r.moveSnakes(nextState, moves)
	if err != nil {
		return nil, err
	}

	err = r.reduceSnakeHealth(nextState)
	if err != nil {
		return nil, err
	}

	err = r.maybeFeedSnakes(nextState)
	if err != nil {
		return nil, err
	}

	err = r.maybeSpawnFood(nextState)
	if err != nil {
		return nil, err
	}

	err = r.maybeEliminateSnakes(nextState)
	if err != nil {
		return nil, err
	}

	return nextState, nil
}

func (r *WrappedRuleset) moveSnakes(b *BoardState, moves []SnakeMove) error {
	err := r.StandardRuleset.moveSnakes(b, moves)
	if err != nil {
		return err
	}

	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		snake.Body[0].X = replace(snake.Body[0].X, 0, b.Width-1)
		snake.Body[0].Y = replace(snake.Body[0].Y, 0, b.Height-1)
	}

	return nil
}
