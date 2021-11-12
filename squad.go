package rules

import (
	"errors"
)

type SquadRuleset struct {
	StandardRuleset

	SquadMap map[string]string `json:"squad_map"`

	// These are intentionally designed so that they default to a standard game.
	AllowBodyCollisions bool `json:"allow_body_collisions"`
	SharedElimination   bool `json:"shared_elimination"`
	SharedHealth        bool `json:"shared_health"`
	SharedLength        bool `json:"shared_length"`
}

const EliminatedBySquad = "squad-eliminated"

func (r *SquadRuleset) Name() string { return "squad" }

func (r *SquadRuleset) CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
	nextBoardState, err := r.StandardRuleset.CreateNextBoardState(prevState, moves)
	if err != nil {
		return nil, err
	}

	// TODO: LOG?
	err = r.resurrectSquadBodyCollisions(nextBoardState)
	if err != nil {
		return nil, err
	}

	// TODO: LOG?
	err = r.shareSquadAttributes(nextBoardState)
	if err != nil {
		return nil, err
	}

	return nextBoardState, nil
}

func (r *SquadRuleset) areSnakesOnSameSquad(snake *Snake, other *Snake) bool {
	return r.areSnakeIDsOnSameSquad(snake.ID, other.ID)
}

func (r *SquadRuleset) areSnakeIDsOnSameSquad(snakeID string, otherID string) bool {
	return r.SquadMap[snakeID] == r.SquadMap[otherID]
}

func (r *SquadRuleset) resurrectSquadBodyCollisions(b *BoardState) error {
	if !r.AllowBodyCollisions {
		return nil
	}

	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause == EliminatedByCollision {
			if snake.EliminatedBy == "" {
				return errors.New("snake eliminated by collision and eliminatedby is not set")
			}
			if snake.ID != snake.EliminatedBy && r.areSnakeIDsOnSameSquad(snake.ID, snake.EliminatedBy) {
				snake.EliminatedCause = NotEliminated
				snake.EliminatedBy = ""
			}
		}
	}

	return nil
}

func (r *SquadRuleset) shareSquadAttributes(b *BoardState) error {
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
			if r.areSnakesOnSameSquad(snake, other) {
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
						snake.EliminatedCause = EliminatedBySquad
						// We intentionally do not set snake.EliminatedBy because there might be multiple culprits.
						snake.EliminatedBy = ""
					}
				}
			}
		}
	}

	return nil
}

func (r *SquadRuleset) IsGameOver(b *BoardState) (bool, error) {
	snakesRemaining := []*Snake{}
	for i := 0; i < len(b.Snakes); i++ {
		if b.Snakes[i].EliminatedCause == NotEliminated {
			snakesRemaining = append(snakesRemaining, &b.Snakes[i])
		}
	}

	for i := 0; i < len(snakesRemaining); i++ {
		if !r.areSnakesOnSameSquad(snakesRemaining[i], snakesRemaining[0]) {
			// There are multiple squads remaining
			return false, nil
		}
	}
	// no snakes or single squad remaining
	return true, nil
}
