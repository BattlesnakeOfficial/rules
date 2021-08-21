package rules

type ConstrictorRuleset struct {
	StandardRuleset
}

func (r *ConstrictorRuleset) ModifyInitialBoardState(initialBoardState *BoardState) (*BoardState, error) {
	initialBoardState, err := r.StandardRuleset.ModifyInitialBoardState(initialBoardState)
	if err != nil {
		return nil, err
	}
	err = r.applyConstrictorRules(initialBoardState)
	if err != nil {
		return nil, err
	}

	return initialBoardState, nil
}

func (r *ConstrictorRuleset) CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {

	nextState, err := r.StandardRuleset.CreateNextBoardState(prevState, moves)
	if err != nil {
		return nil, err
	}

	err = r.applyConstrictorRules(nextState)
	if err != nil {
		return nil, err
	}

	return nextState, nil
}

func (r *ConstrictorRuleset) applyConstrictorRules(b *BoardState) error {
	// Remove all food from the board
	b.Food = []Point{}

	// Set all snakes to max health and ensure they grow next turn
	for i := 0; i < len(b.Snakes); i++ {
		b.Snakes[i].Health = SnakeMaxHealth

		tail := b.Snakes[i].Body[len(b.Snakes[i].Body)-1]
		subTail := b.Snakes[i].Body[len(b.Snakes[i].Body)-2]
		if tail != subTail {
			r.growSnake(&b.Snakes[i])
		}
	}

	return nil
}
