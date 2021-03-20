package rules

type WrappedRuleset struct {
	StandardRuleset
}

func (r *WrappedRuleset) snakeIsOutOfBounds(s *Snake, boardWidth int32, boardHeight int32) bool {
	return false
}

func replace(value, min, max int32) int32 {
	if value < min {return max}
	if value > max {return min}
	return value
}

func (r *WrappedRuleset) moveSnakes(b *BoardState, moves []SnakeMove) error {

	err := r.StandardRuleset.moveSnakes(b, moves)
	if err != nil {
		return err
	}

	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause != NotEliminated {
			continue
		}
		snake.Body[0].X = replace(snake.Body[0].X, 0, b.Width-1)
		snake.Body[0].Y = replace(snake.Body[0].Y, 0, b.Width-1)
	}
	return nil
}