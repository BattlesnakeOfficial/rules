package rules

var constrictorRulesetStages = []string{
	StageMovementStandard,
	StageStarvationStandard,
	StageHazardDamageStandard,
	StageFeedSnakesStandard,
	StageEliminationStandard,
	StageSpawnFoodNoFood,
	StageModifySnakesAlwaysGrow,
	StageGameOverStandard,
}

type ConstrictorRuleset struct {
	StandardRuleset
}

func (r *ConstrictorRuleset) Name() string { return GameTypeConstrictor }

func (r ConstrictorRuleset) Execute(bs *BoardState, s Settings, sm []SnakeMove) (bool, *BoardState, error) {
	return NewPipeline(constrictorRulesetStages...).Execute(bs, s, sm)
}

func (r *ConstrictorRuleset) ModifyInitialBoardState(initialBoardState *BoardState) (*BoardState, error) {
	_, nextState, err := r.Execute(initialBoardState, r.Settings(), nil)
	return nextState, err
}

func (r *ConstrictorRuleset) CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
	_, nextState, err := r.Execute(prevState, r.Settings(), moves)

	return nextState, err
}

func (r *ConstrictorRuleset) IsGameOver(b *BoardState) (bool, error) {
	return GameOverStandard(b, r.Settings(), nil)
}

func RemoveFoodConstrictor(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
	// Remove all food from the board
	b.Food = []Point{}

	return false, nil
}

func GrowSnakesConstrictor(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
	// Set all snakes to max health and ensure they grow next turn
	for i := 0; i < len(b.Snakes); i++ {
		if len(b.Snakes[i].Body) <= 0 {
			return false, ErrorZeroLengthSnake
		}
		b.Snakes[i].Health = SnakeMaxHealth

		tail := b.Snakes[i].Body[len(b.Snakes[i].Body)-1]
		subTail := b.Snakes[i].Body[len(b.Snakes[i].Body)-2]
		if tail != subTail {
			growSnake(&b.Snakes[i])
		}
	}

	return false, nil
}
