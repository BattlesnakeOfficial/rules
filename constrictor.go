package rules

type ConstrictorRuleset struct {
	StandardRuleset
}

func (r *ConstrictorRuleset) Name() string { return GameTypeConstrictor }

func (r ConstrictorRuleset) Pipeline() (*Pipeline, error) {
	return NewPipeline(
		"snake.movement.standard",
		"health.reduce.standard",
		"hazard.damage.standard",
		"snake.eatfood.standard",
		"food.spawn.standard",
		"snake.eliminate.standard",
		"food.remove.constrictor",
		"snake.grow.constrictor",
		"gameover.standard",
	)
}

func (r *ConstrictorRuleset) ModifyInitialBoardState(initialBoardState *BoardState) (*BoardState, error) {
	p, err := r.Pipeline()
	if err != nil {
		return nil, err
	}

	_, nextState, err := p.Execute(initialBoardState, r.Settings(), nil)
	return nextState, err
}

func (r *ConstrictorRuleset) CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
	p, err := r.Pipeline()
	if err != nil {
		return nil, err
	}
	_, nextState, err := p.Execute(prevState, r.Settings(), moves)

	return nextState, err
}

func RemoveFoodConstrictor(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
	// Remove all food from the board
	b.Food = []Point{}

	return false, nil
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
