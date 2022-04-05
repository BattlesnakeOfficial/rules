package rules

type ConstrictorRuleset struct {
	StandardRuleset
}

func (r *ConstrictorRuleset) Name() string { return GameTypeConstrictor }

func (r ConstrictorRuleset) Pipeline() (*Pipeline, error) {
	// The constrictor pipeline extends the standard pipeline
	standard, err := r.StandardRuleset.Pipeline()
	if err != nil {
		return nil, err
	}

	constrictor, err := NewPipeline(
		"removefood.constrictor",
		"growsnake.constrictor",
	)
	if err != nil {
		return nil, err
	}
	return standard.Append(constrictor), nil
}

func (r *ConstrictorRuleset) ModifyInitialBoardState(initialBoardState *BoardState) (*BoardState, error) {
	initialBoardState, err := r.StandardRuleset.ModifyInitialBoardState(initialBoardState)
	if err != nil {
		return nil, err
	}

	r.removeFood(initialBoardState)

	err = r.applyConstrictorRules(initialBoardState)
	if err != nil {
		return nil, err
	}

	return initialBoardState, nil
}

func (r *ConstrictorRuleset) CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
	nextState := prevState.Clone()

	p, err := r.Pipeline()
	if err != nil {
		return nil, err
	}
	_, err = p.Execute(nextState, r.Settings(), moves)

	return nextState, err
}

func (r *ConstrictorRuleset) removeFood(b *BoardState) {
	_, _ = r.callStageFunc(RemoveFoodConstrictor, b, []SnakeMove{})
}

func RemoveFoodConstrictor(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
	// Remove all food from the board
	b.Food = []Point{}

	return false, nil
}

func (r *ConstrictorRuleset) applyConstrictorRules(b *BoardState) error {
	_, err := r.callStageFunc(GrowSnakesConstrictor, b, []SnakeMove{})
	return err
}

func GrowSnakesConstrictor(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
	// Set all snakes to max health and ensure they grow next turn
	for i := 0; i < len(b.Snakes); i++ {
		b.Snakes[i].Health = SnakeMaxHealth

		tail := b.Snakes[i].Body[len(b.Snakes[i].Body)-1]
		subTail := b.Snakes[i].Body[len(b.Snakes[i].Body)-2]
		if tail != subTail {
			growSnake(&b.Snakes[i])
		}
	}

	return false, nil
}
