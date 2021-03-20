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
	nextBoardState, err := r.StandardRuleset.CreateNextBoardState(prevState, moves)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(nextBoardState.Snakes); i++ {
		snake := &nextBoardState.Snakes[i]
		if snake.EliminatedCause == EliminatedByOutOfBounds {
			snake.EliminatedCause = NotEliminated
			snake.EliminatedBy = ""
			snake.Body[0].X = replace(snake.Body[0].X, 0, nextBoardState.Width-1)
			snake.Body[0].Y = replace(snake.Body[0].Y, 0, nextBoardState.Height-1)
		}
	}

	return nextBoardState, nil
}
