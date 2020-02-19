package rules

import (
	"errors"
)

type TeamRuleset struct {
	StandardRuleset

	TeamMap map[string]string

	AllowBodyCollisions bool
	SharedElimination   bool
	SharedLength        bool
	SharedHealth        bool
}

const EliminatedByTeam = "team-is-eliminated"

func (r *TeamRuleset) CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
	nextBoardState, err := r.StandardRuleset.CreateNextBoardState(prevState, moves)
	if err != nil {
		return nil, err
	}

	// TODO: LOG?
	err = r.resurrectTeamBodyCollisions(nextBoardState)
	if err != nil {
		return nil, err
	}

	// TODO: LOG?
	err = r.shareTeamAttributes(nextBoardState)
	if err != nil {
		return nil, err
	}

	return nextBoardState, nil
}

func (r *TeamRuleset) areSnakesOnSameTeam(snake *Snake, other *Snake) bool {
	return r.areSnakeIDsOnSameTeam(snake.ID, other.ID)
}

func (r *TeamRuleset) areSnakeIDsOnSameTeam(snakeID string, otherID string) bool {
	return snakeID != otherID && r.TeamMap[snakeID] == r.TeamMap[otherID]
}

func (r *TeamRuleset) resurrectTeamBodyCollisions(b *BoardState) error {
	if !(r.AllowBodyCollisions) {
		return nil
	}

	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause == EliminatedByCollision {
			if snake.EliminatedBy == "" {
				return errors.New("snake eliminated by collision and eliminatedby is not set")
			}
			if r.areSnakeIDsOnSameTeam(snake.ID, snake.EliminatedBy) {
				snake.EliminatedCause = NotEliminated
				snake.EliminatedBy = ""
			}
		}
	}

	return nil
}

func (r *TeamRuleset) shareTeamAttributes(b *BoardState) error {
	if !(r.SharedElimination || r.SharedLength || r.SharedHealth) {
		return nil
	}

	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause != NotEliminated {
			continue
		}

		for j := 0; j < len(b.Snakes); j++ {
			other := &b.Snakes[j]
			if r.areSnakesOnSameTeam(snake, other) {

				if r.SharedHealth {
					if snake.Health < other.Health {
						snake.Health = other.Health
					}
				}

				if r.SharedLength {
					if len(snake.Body) == 0 || len(other.Body) == 0 {
						return errors.New("found snake of zero length")
					}
					for len(snake.Body) < len(other.Body) {
						r.growSnake(snake)
					}
				}

				if r.SharedElimination {
					if snake.EliminatedCause == NotEliminated && other.EliminatedCause != NotEliminated {
						snake.EliminatedCause = EliminatedByTeam
						// We intentionally do not set snake.EliminatedBy to not place blame,
						// especially when there might be multiple snakes eliminated
					}
				}
			}
		}
	}

	return nil
}
