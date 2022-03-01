package rules

import (
	"math/rand"
	"sort"
)

type StandardRuleset struct {
	FoodSpawnChance     int32 // [0, 100]
	MinimumFood         int32
	HazardDamagePerTurn int32
}

func (r *StandardRuleset) Name() string { return "standard" }

func (r *StandardRuleset) ModifyInitialBoardState(initialState *BoardState) (*BoardState, error) {
	// No-op
	return initialState, nil
}

func (r *StandardRuleset) CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
	// We specifically want to copy prevState, so as not to alter it directly.
	nextState := prevState.Clone()

	err := r.moveSnakes(nextState, moves)
	if err != nil {
		return nil, err
	}

	err = r.reduceSnakeHealth(nextState)
	if err != nil {
		return nil, err
	}

	err = r.maybeDamageHazards(nextState)
	if err != nil {
		return nil, err
	}

	// bvanvugt: We specifically want this to happen before elimination for two reasons:
	// 1) We want snakes to be able to eat on their very last turn and still survive.
	// 2) So that head-to-head collisions on food still remove the food.
	//    This does create an artifact though, where head-to-head collisions
	//    of equal length actually show length + 1 and full health, as if both snakes ate.
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

func (r *StandardRuleset) moveSnakes(b *BoardState, moves []SnakeMove) error {
	_, err := r.callStageFunc(MoveSnakesStandard, b, moves)
	return err
}

func MoveSnakesStandard(b *BoardState, settings RulesetSettings, moves []SnakeMove) (bool, error) {
	// If no moves are passed, pass on modifying the initial board state
	if len(moves) == 0 {
		return false, nil
	}

	// Sanity check that all non-eliminated snakes have moves and bodies.
	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause != NotEliminated {
			continue
		}

		if len(snake.Body) == 0 {
			return false, ErrorZeroLengthSnake
		}
		moveFound := false
		for _, move := range moves {
			if snake.ID == move.ID {
				moveFound = true
				break
			}
		}
		if !moveFound {
			return false, ErrorNoMoveFound
		}
	}

	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause != NotEliminated {
			continue
		}

		for _, move := range moves {
			if move.ID == snake.ID {
				appliedMove := move.Move
				switch move.Move {
				case MoveUp, MoveDown, MoveRight, MoveLeft:
					break
				default:
					appliedMove = getDefaultMove(snake.Body)
				}

				newHead := Point{}
				switch appliedMove {
				// Guaranteed to be one of these options given the clause above
				case MoveUp:
					newHead.X = snake.Body[0].X
					newHead.Y = snake.Body[0].Y + 1
				case MoveDown:
					newHead.X = snake.Body[0].X
					newHead.Y = snake.Body[0].Y - 1
				case MoveLeft:
					newHead.X = snake.Body[0].X - 1
					newHead.Y = snake.Body[0].Y
				case MoveRight:
					newHead.X = snake.Body[0].X + 1
					newHead.Y = snake.Body[0].Y
				}

				// Append new head, pop old tail
				snake.Body = append([]Point{newHead}, snake.Body[:len(snake.Body)-1]...)
			}
		}
	}
	return false, nil
}

func getDefaultMove(snakeBody []Point) string {
	if len(snakeBody) >= 2 {
		// Use neck to determine last move made
		head, neck := snakeBody[0], snakeBody[1]
		// Situations where neck is next to head
		if head.X == neck.X+1 {
			return MoveRight
		} else if head.X == neck.X-1 {
			return MoveLeft
		} else if head.Y == neck.Y+1 {
			return MoveUp
		} else if head.Y == neck.Y-1 {
			return MoveDown
		}
		// Consider the wrapped cases using zero axis to anchor
		if head.X == 0 && neck.X > 0 {
			return MoveRight
		} else if neck.X == 0 && head.X > 0 {
			return MoveLeft
		} else if head.Y == 0 && neck.Y > 0 {
			return MoveUp
		} else if neck.Y == 0 && head.Y > 0 {
			return MoveDown
		}
	}
	return MoveUp
}

func (r *StandardRuleset) reduceSnakeHealth(b *BoardState) error {
	_, err := r.callStageFunc(ReduceSnakeHealthStandard, b, []SnakeMove{})
	return err
}

func ReduceSnakeHealthStandard(b *BoardState, settings RulesetSettings, moves []SnakeMove) (bool, error) {
	for i := 0; i < len(b.Snakes); i++ {
		if b.Snakes[i].EliminatedCause == NotEliminated {
			b.Snakes[i].Health = b.Snakes[i].Health - 1
		}
	}
	return false, nil
}

func (r *StandardRuleset) maybeDamageHazards(b *BoardState) error {
	_, err := r.callStageFunc(DamageHazardsStandard, b, []SnakeMove{})
	return err
}

func DamageHazardsStandard(b *BoardState, settings RulesetSettings, moves []SnakeMove) (bool, error) {
	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause != NotEliminated {
			continue
		}
		head := snake.Body[0]
		for _, p := range b.Hazards {
			if head == p {
				// If there's a food in this square, don't reduce health
				foundFood := false
				for _, food := range b.Food {
					if p == food {
						foundFood = true
					}
				}
				if foundFood {
					continue
				}

				// Snake is in a hazard, reduce health
				snake.Health = snake.Health - settings.HazardDamagePerTurn
				if snake.Health < 0 {
					snake.Health = 0
				}
				if snakeIsOutOfHealth(snake) {
					snake.EliminatedCause = EliminatedByOutOfHealth
				}
			}
		}
	}

	return false, nil
}

func (r *StandardRuleset) maybeEliminateSnakes(b *BoardState) error {
	_, err := r.callStageFunc(EliminateSnakesStandard, b, []SnakeMove{})
	return err
}

func EliminateSnakesStandard(b *BoardState, settings RulesetSettings, moves []SnakeMove) (bool, error) {
	// First order snake indices by length.
	// In multi-collision scenarios we want to always attribute elimination to the longest snake.
	snakeIndicesByLength := make([]int, len(b.Snakes))
	for i := 0; i < len(b.Snakes); i++ {
		snakeIndicesByLength[i] = i
	}
	sort.Slice(snakeIndicesByLength, func(i int, j int) bool {
		lenI := len(b.Snakes[snakeIndicesByLength[i]].Body)
		lenJ := len(b.Snakes[snakeIndicesByLength[j]].Body)
		return lenI > lenJ
	})

	// First, iterate over all non-eliminated snakes and eliminate the ones
	// that are out of health or have moved out of bounds.
	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause != NotEliminated {
			continue
		}
		if len(snake.Body) <= 0 {
			return false, ErrorZeroLengthSnake
		}

		if snakeIsOutOfHealth(snake) {
			snake.EliminatedCause = EliminatedByOutOfHealth
			continue
		}

		if snakeIsOutOfBounds(snake, b.Width, b.Height) {
			snake.EliminatedCause = EliminatedByOutOfBounds
			continue
		}
	}

	// Next, look for any collisions. Note we apply collision eliminations
	// after this check so that snakes can collide with each other and be properly eliminated.
	type CollisionElimination struct {
		ID    string
		Cause string
		By    string
	}
	collisionEliminations := []CollisionElimination{}
	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause != NotEliminated {
			continue
		}
		if len(snake.Body) <= 0 {
			return false, ErrorZeroLengthSnake
		}

		// Check for self-collisions first
		if snakeHasBodyCollided(snake, snake) {
			collisionEliminations = append(collisionEliminations, CollisionElimination{
				ID:    snake.ID,
				Cause: EliminatedBySelfCollision,
				By:    snake.ID,
			})
			continue
		}

		// Check for body collisions with other snakes second
		hasBodyCollided := false
		for _, otherIndex := range snakeIndicesByLength {
			other := &b.Snakes[otherIndex]
			if other.EliminatedCause != NotEliminated {
				continue
			}
			if snake.ID != other.ID && snakeHasBodyCollided(snake, other) {
				collisionEliminations = append(collisionEliminations, CollisionElimination{
					ID:    snake.ID,
					Cause: EliminatedByCollision,
					By:    other.ID,
				})
				hasBodyCollided = true
				break
			}
		}
		if hasBodyCollided {
			continue
		}

		// Check for head-to-heads last
		hasHeadCollided := false
		for _, otherIndex := range snakeIndicesByLength {
			other := &b.Snakes[otherIndex]
			if other.EliminatedCause != NotEliminated {
				continue
			}
			if snake.ID != other.ID && snakeHasLostHeadToHead(snake, other) {
				collisionEliminations = append(collisionEliminations, CollisionElimination{
					ID:    snake.ID,
					Cause: EliminatedByHeadToHeadCollision,
					By:    other.ID,
				})
				hasHeadCollided = true
				break
			}
		}
		if hasHeadCollided {
			continue
		}
	}

	// Apply collision eliminations
	for _, elimination := range collisionEliminations {
		for i := 0; i < len(b.Snakes); i++ {
			snake := &b.Snakes[i]
			if snake.ID == elimination.ID {
				snake.EliminatedCause = elimination.Cause
				snake.EliminatedBy = elimination.By
				break
			}
		}
	}

	return false, nil
}

func snakeIsOutOfHealth(s *Snake) bool {
	return s.Health <= 0
}

func snakeIsOutOfBounds(s *Snake, boardWidth int32, boardHeight int32) bool {
	for _, point := range s.Body {
		if (point.X < 0) || (point.X >= boardWidth) {
			return true
		}
		if (point.Y < 0) || (point.Y >= boardHeight) {
			return true
		}
	}
	return false
}

func snakeHasBodyCollided(s *Snake, other *Snake) bool {
	head := s.Body[0]
	for i, body := range other.Body {
		if i == 0 {
			continue
		} else if head.X == body.X && head.Y == body.Y {
			return true
		}
	}
	return false
}

func snakeHasLostHeadToHead(s *Snake, other *Snake) bool {
	if s.Body[0].X == other.Body[0].X && s.Body[0].Y == other.Body[0].Y {
		return len(s.Body) <= len(other.Body)
	}
	return false
}

func (r *StandardRuleset) maybeFeedSnakes(b *BoardState) error {
	_, err := r.callStageFunc(FeedSnakesStandard, b, []SnakeMove{})
	return err
}

func FeedSnakesStandard(b *BoardState, settings RulesetSettings, moves []SnakeMove) (bool, error) {
	newFood := []Point{}
	for _, food := range b.Food {
		foodHasBeenEaten := false
		for i := 0; i < len(b.Snakes); i++ {
			snake := &b.Snakes[i]

			// Ignore eliminated and zero-length snakes, they can't eat.
			if snake.EliminatedCause != NotEliminated || len(snake.Body) == 0 {
				continue
			}

			if snake.Body[0].X == food.X && snake.Body[0].Y == food.Y {
				feedSnake(snake)
				foodHasBeenEaten = true
			}
		}
		// Persist food to next BoardState if not eaten
		if !foodHasBeenEaten {
			newFood = append(newFood, food)
		}
	}

	b.Food = newFood
	return false, nil
}

func feedSnake(snake *Snake) {
	growSnake(snake)
	snake.Health = SnakeMaxHealth
}

func growSnake(snake *Snake) {
	if len(snake.Body) > 0 {
		snake.Body = append(snake.Body, snake.Body[len(snake.Body)-1])
	}
}

func (r *StandardRuleset) maybeSpawnFood(b *BoardState) error {
	_, err := r.callStageFunc(SpawnFoodStandard, b, []SnakeMove{})
	return err
}

func SpawnFoodStandard(b *BoardState, settings RulesetSettings, moves []SnakeMove) (bool, error) {
	numCurrentFood := int32(len(b.Food))
	if numCurrentFood < settings.MinimumFood {
		return false, PlaceFoodRandomly(b, settings.MinimumFood-numCurrentFood)
	}
	if settings.FoodSpawnChance > 0 && int32(rand.Intn(100)) < settings.FoodSpawnChance {
		return false, PlaceFoodRandomly(b, 1)
	}
	return false, nil
}

func (r *StandardRuleset) IsGameOver(b *BoardState) (bool, error) {
	return r.callStageFunc(GameOverStandard, b, []SnakeMove{})
}

func GameOverStandard(b *BoardState, settings RulesetSettings, moves []SnakeMove) (bool, error) {
	numSnakesRemaining := 0
	for i := 0; i < len(b.Snakes); i++ {
		if b.Snakes[i].EliminatedCause == NotEliminated {
			numSnakesRemaining++
		}
	}
	return numSnakesRemaining <= 1, nil
}

// Adaptor for integrating stages into StandardRuleset
func (r *StandardRuleset) callStageFunc(stage StageFunc, boardState *BoardState, moves []SnakeMove) (bool, error) {
	return stage(boardState, RulesetSettings{
		FoodSpawnChance:     r.FoodSpawnChance,
		MinimumFood:         r.MinimumFood,
		HazardDamagePerTurn: r.HazardDamagePerTurn,
	}, moves)
}
