package rules

import (
	"math/rand"
	"sort"
)

var standardRulesetStages = []string{
	StageGameOverStandard,
	StageMovementStandard,
	StageStarvationStandard,
	StageHazardDamageStandard,
	StageFeedSnakesStandard,
	StageEliminationStandard,
}

func MoveSnakesStandard(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
	if IsInitialization(b, settings, moves) {
		return false, nil
	}

	// no-op when moves are empty
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

func ReduceSnakeHealthStandard(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
	if IsInitialization(b, settings, moves) {
		return false, nil
	}
	for i := 0; i < len(b.Snakes); i++ {
		if b.Snakes[i].EliminatedCause == NotEliminated {
			b.Snakes[i].Health = b.Snakes[i].Health - 1
		}
	}
	return false, nil
}

func DamageHazardsStandard(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
	if IsInitialization(b, settings, moves) {
		return false, nil
	}
	hazardDamage := settings.Int(ParamHazardDamagePerTurn, 0)
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
				snake.Health = snake.Health - hazardDamage
				if snake.Health < 0 {
					snake.Health = 0
				}
				if snake.Health > SnakeMaxHealth {
					snake.Health = SnakeMaxHealth
				}
				if snakeIsOutOfHealth(snake) {
					EliminateSnake(snake, EliminatedByHazard, "", b.Turn+1)
				}
			}
		}
	}

	return false, nil
}

func EliminateSnakesStandard(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
	if IsInitialization(b, settings, moves) {
		return false, nil
	}
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
			EliminateSnake(snake, EliminatedByOutOfHealth, "", b.Turn+1)
			continue
		}

		if snakeIsOutOfBounds(snake, b.Width, b.Height) {
			EliminateSnake(snake, EliminatedByOutOfBounds, "", b.Turn+1)
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
				EliminateSnake(snake, elimination.Cause, elimination.By, b.Turn+1)
				break
			}
		}
	}

	return false, nil
}

func snakeIsOutOfHealth(s *Snake) bool {
	return s.Health <= 0
}

func snakeIsOutOfBounds(s *Snake, boardWidth int, boardHeight int) bool {
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

func FeedSnakesStandard(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
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

// Deprecated: handled by maps.Standard
func SpawnFoodStandard(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
	if IsInitialization(b, settings, moves) {
		return false, nil
	}
	minimumFood := settings.Int(ParamMinimumFood, 0)
	foodSpawnChance := settings.Int(ParamFoodSpawnChance, 0)
	numCurrentFood := int(len(b.Food))
	if numCurrentFood < minimumFood {
		return false, PlaceFoodRandomly(GlobalRand, b, minimumFood-numCurrentFood)
	}
	if foodSpawnChance > 0 && int(rand.Intn(100)) < foodSpawnChance {
		return false, PlaceFoodRandomly(GlobalRand, b, 1)
	}
	return false, nil
}

func GameOverStandard(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
	numSnakesRemaining := 0
	for i := 0; i < len(b.Snakes); i++ {
		if b.Snakes[i].EliminatedCause == NotEliminated {
			numSnakesRemaining++
		}
	}
	return numSnakesRemaining <= 1, nil
}
