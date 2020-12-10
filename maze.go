package rules

import ()

type MazeRuleset struct {
	StandardRuleset
}

func (r *MazeRuleset) CreateInitialBoardState(width int32, height int32, snakeIDs []string) (*BoardState, error) {
	initialBoardState, err := r.StandardRuleset.CreateInitialBoardState(width, height, snakeIDs)
	if err != nil {
		return nil, err
	}

	err = r.fillBoardWithFood(initialBoardState)
	if err != nil {
		return nil, err
	}

	return initialBoardState, nil
}

func (r *MazeRuleset) CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
	nextState, err := r.StandardRuleset.CreateNextBoardState(prevState, moves)
	if err != nil {
		return nil, err
	}

	err = r.fillBoardWithFood(nextState)
	if err != nil {
		return nil, err
	}

	return nextState, nil
}

func (r *MazeRuleset) fillBoardWithFood(b *BoardState) error {
	unoccupiedPoints := r.getUnoccupiedPoints(b, true)
	b.Food = append(b.Food, unoccupiedPoints...)
	return nil
}
