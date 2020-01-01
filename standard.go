package rulesets

import (
// log "github.com/sirupsen/logrus"
)

type StandardRuleset struct{}

const (
	SNAKE_MAX_HEALTH = 100
	// bvanvugt - TODO: Just return formatted strings instead of codes?
	ELIMINATED_COLLISION      = "snake-collision"
	ELIMINATED_SELF_COLLISION = "snake-self-collision"
	ELIMINATED_STARVATION     = "starvation"
	ELIMINATED_HEAD_TO_HEAD   = "head-collision"
	ELIMINATED_OUT_OF_BOUNDS  = "wall-collision"
)

func (r *StandardRuleset) ResolveMoves(g *Game, prevGameState *GameState, moves []*SnakeMove) (nextGameState *GameState, err error) {
	nextGameState = &GameState{
		Snakes: prevGameState.Snakes,
		Food:   prevGameState.Food,
	}

	// TODO: Gut check the GameState?

	// TODO: LOG?
	err = r.moveSnakes(nextGameState, moves)
	if err != nil {
		return nil, err
	}

	// TODO: LOG?
	err = r.reduceSnakeHealth(nextGameState)
	if err != nil {
		return nil, err
	}

	// TODO: LOG?
	// bvanvugt: we specifically want this to happen before elimination
	// so that head-to-head collisions on food still remove the food.
	err = r.feedSnakes(nextGameState)
	if err != nil {
		return nil, err
	}

	// TODO: LOG?
	err = r.eliminateSnakes(nextGameState, g.Width, g.Height)
	if err != nil {
		return nil, err
	}

	// TODO: LOG?
	err = r.maybeSpawnFood(nextGameState)
	if err != nil {
		return nil, err
	}

	return nextGameState, nil
}

func (r *StandardRuleset) moveSnakes(gs *GameState, moves []*SnakeMove) error {
	for _, move := range moves {
		var newHead = &Point{}
		switch move.Move {
		case MOVE_DOWN:
			newHead.X = move.Snake.Body[0].X
			newHead.Y = move.Snake.Body[0].Y + 1
		case MOVE_LEFT:
			newHead.X = move.Snake.Body[0].X - 1
			newHead.Y = move.Snake.Body[0].Y
		case MOVE_RIGHT:
			newHead.X = move.Snake.Body[0].X + 1
			newHead.Y = move.Snake.Body[0].Y
		case MOVE_UP:
			newHead.X = move.Snake.Body[0].X
			newHead.Y = move.Snake.Body[0].Y - 1
		default:
			// Default to UP
			var dX int32 = 0
			var dY int32 = -1
			// If neck is available, use neck to determine last direction
			if len(move.Snake.Body) >= 2 {
				dX = move.Snake.Body[0].X - move.Snake.Body[1].X
				dY = move.Snake.Body[0].Y - move.Snake.Body[1].Y
				if dX == 0 && dY == 0 {
					dY = -1 // Move up if no last move was made
				}
			}
			// Apply
			newHead.X = move.Snake.Body[0].X + dX
			newHead.Y = move.Snake.Body[0].Y + dY
		}

		// Append new head, pop old tail
		move.Snake.Body = append([]*Point{newHead}, move.Snake.Body[:len(move.Snake.Body)-1]...)
	}
	return nil
}

func (r *StandardRuleset) reduceSnakeHealth(gs *GameState) error {
	for _, snake := range gs.Snakes {
		snake.Health = snake.Health - 1
	}
	return nil
}

func (r *StandardRuleset) eliminateSnakes(gs *GameState, boardWidth int32, boardHeight int32) error {
	for _, snake := range gs.Snakes {
		if r.snakeHasStarved(snake) {
			snake.EliminatedCause = ELIMINATED_STARVATION
		} else if r.snakeIsOutOfBounds(snake, boardWidth, boardHeight) {
			snake.EliminatedCause = ELIMINATED_OUT_OF_BOUNDS
		} else {
			for _, other := range gs.Snakes {
				if r.snakeHasBodyCollided(snake, other) {
					if snake.ID == other.ID {
						snake.EliminatedCause = ELIMINATED_SELF_COLLISION
					} else {
						snake.EliminatedCause = ELIMINATED_COLLISION
					}
					break
				} else if r.snakeHasLostHeadToHead(snake, other) {
					snake.EliminatedCause = ELIMINATED_HEAD_TO_HEAD
					break
				}
			}
		}
	}
	return nil
}

func (r *StandardRuleset) snakeHasStarved(s *Snake) bool {
	return s.Health <= 0
}

func (r *StandardRuleset) snakeIsOutOfBounds(s *Snake, boardWidth int32, boardHeight int32) bool {
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

func (r *StandardRuleset) snakeHasBodyCollided(s *Snake, other *Snake) bool {
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

func (r *StandardRuleset) snakeHasLostHeadToHead(s *Snake, other *Snake) bool {
	if s.Body[0].X == other.Body[0].X && s.Body[0].Y == other.Body[0].Y {
		return len(s.Body) <= len(other.Body)
	}
	return false
}

func (r *StandardRuleset) feedSnakes(gs *GameState) error {
	var newFood []*Point
	var tail *Point

	for _, food := range gs.Food {
		foodHasBeenEaten := false
		for _, snake := range gs.Snakes {
			if snake.Body[0].X == food.X && snake.Body[0].Y == food.Y {
				foodHasBeenEaten = true
				// Update snake
				snake.Health = SNAKE_MAX_HEALTH
				tail = snake.Body[len(snake.Body)-1]
				snake.Body = append(snake.Body, &Point{X: tail.X, Y: tail.Y})
			}
		}
		// Persist food to next GameState if not eaten
		if !foodHasBeenEaten {
			newFood = append(newFood, food)
		}
	}

	gs.Food = newFood
	return nil
}

func (r *StandardRuleset) maybeSpawnFood(gs *GameState) error {
	// TODO
	return nil
}
