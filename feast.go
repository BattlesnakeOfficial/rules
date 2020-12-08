package rules

import ()

type FeastRuleset struct {
	StandardRuleset
}

func (r *FeastRuleset) CreateInitialBoardState(width int32, height int32, snakeIDs []string) (*BoardState, error) {
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

func (r *FeastRuleset) CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
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

func (r *FeastRuleset) fillBoardWithFood(b *BoardState) error {
	unoccupiedPoints := r.getUnoccupiedPoints(b, true)
	b.Food = append(b.Food, unoccupiedPoints...)
	return nil
}
