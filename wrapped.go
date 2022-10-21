package rules

var wrappedRulesetStages = []string{
	StageGameOverStandard,
	StageMovementWrapBoundaries,
	StageStarvationStandard,
	StageHazardDamageStandard,
	StageFeedSnakesStandard,
	StageEliminationStandard,
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

func wrap(value, min, max int) int {
	if value < min {
		return max
	}
	if value > max {
		return min
	}
	return value
}
