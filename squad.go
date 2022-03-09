package rules

import (
	"errors"

	"github.com/tidwall/sjson"
)

type SquadRuleset struct {
	StandardRuleset

	SquadMap map[string]string

	// These are intentionally designed so that they default to a standard game.
	AllowBodyCollisions bool
	SharedElimination   bool
	SharedHealth        bool
	SharedLength        bool
}

const EliminatedBySquad = "squad-eliminated"

func (r *SquadRuleset) Name() string { return "squad" }

func (r *SquadRuleset) CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
	nextBoardState, err := r.StandardRuleset.CreateNextBoardState(prevState, moves)
	if err != nil {
		return nil, err
	}

	err = r.resurrectSquadBodyCollisions(nextBoardState)
	if err != nil {
		return nil, err
	}

	err = r.shareSquadAttributes(nextBoardState)
	if err != nil {
		return nil, err
	}

	return nextBoardState, nil
}

func areSnakesOnSameSquad(settings SettingsJSON, snake *Snake, other *Snake) bool {
	return areSnakeIDsOnSameSquad(settings, snake.ID, other.ID)
}

func areSnakeIDsOnSameSquad(settings SettingsJSON, snakeID string, otherID string) bool {
	s1 := settings.GetString("squad", "squadMap", snakeID)
	if s1 == "" {
		// snake not on a squad, so they can't be on the same squad
		return false
	}

	s2 := settings.GetString("squad", "squadMap", otherID)

	return s1 == s2
}

func (r *SquadRuleset) resurrectSquadBodyCollisions(b *BoardState) error {
	_, err := r.callStageFunc(ResurrectSnakesSquad, b, []SnakeMove{})
	return err
}

func ResurrectSnakesSquad(b *BoardState, settings SettingsJSON, moves []SnakeMove) (bool, error) {
	if !settings.GetBool("squad", "allowBodyCollisions") {
		return false, nil
	}

	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause == EliminatedByCollision {
			if snake.EliminatedBy == "" {
				return false, errors.New("snake eliminated by collision and eliminatedby is not set")
			}
			if snake.ID != snake.EliminatedBy && areSnakeIDsOnSameSquad(settings, snake.ID, snake.EliminatedBy) {
				snake.EliminatedCause = NotEliminated
				snake.EliminatedBy = ""
			}
		}
	}

	return false, nil
}

func (r *SquadRuleset) shareSquadAttributes(b *BoardState) error {
	_, err := r.callStageFunc(ShareAttributesSquad, b, []SnakeMove{})
	return err
}

func ShareAttributesSquad(b *BoardState, settings SettingsJSON, moves []SnakeMove) (bool, error) {
	sharedElimination := settings.GetBool("squad", "sharedElimination")
	sharedLength := settings.GetBool("squad", "sharedLength")
	sharedHealth := settings.GetBool("squad", "sharedHealth")

	if !(sharedElimination || sharedLength || sharedHealth) {
		return false, nil
	}

	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause != NotEliminated {
			continue
		}

		for j := 0; j < len(b.Snakes); j++ {
			other := &b.Snakes[j]
			if areSnakesOnSameSquad(settings, snake, other) {
				if sharedHealth {
					if snake.Health < other.Health {
						snake.Health = other.Health
					}
				}
				if sharedLength {
					if len(snake.Body) == 0 || len(other.Body) == 0 {
						return false, errors.New("found snake of zero length")
					}
					for len(snake.Body) < len(other.Body) {
						growSnake(snake)
					}
				}
				if sharedElimination {
					if snake.EliminatedCause == NotEliminated && other.EliminatedCause != NotEliminated {
						snake.EliminatedCause = EliminatedBySquad
						// We intentionally do not set snake.EliminatedBy because there might be multiple culprits.
						snake.EliminatedBy = ""
					}
				}
			}
		}
	}

	return false, nil
}

func (r *SquadRuleset) IsGameOver(b *BoardState) (bool, error) {
	return r.callStageFunc(GameOverSquad, b, []SnakeMove{})
}

func GameOverSquad(b *BoardState, settings SettingsJSON, moves []SnakeMove) (bool, error) {
	snakesRemaining := []*Snake{}
	for i := 0; i < len(b.Snakes); i++ {
		if b.Snakes[i].EliminatedCause == NotEliminated {
			snakesRemaining = append(snakesRemaining, &b.Snakes[i])
		}
	}

	for i := 0; i < len(snakesRemaining); i++ {
		if !areSnakesOnSameSquad(settings, snakesRemaining[i], snakesRemaining[0]) {
			// There are multiple squads remaining
			return false, nil
		}
	}
	// no snakes or single squad remaining
	return true, nil
}

func (r *SquadRuleset) getSettingsJSON() (SettingsJSON, error) {
	j, err := r.StandardRuleset.getSettingsJSON()
	if err != nil {
		return nil, err
	}

	// add all the public API squad settings
	j, err = sjson.SetBytes(j, "squad", SquadSettings{
		AllowBodyCollisions: r.AllowBodyCollisions,
		SharedElimination:   r.SharedElimination,
		SharedHealth:        r.SharedHealth,
		SharedLength:        r.SharedLength,
	})
	if err != nil {
		return nil, err
	}

	// patch in the squad map
	j, err = sjson.SetBytes(j, "squad.squadMap", r.SquadMap)

	return j, err
}

// Adaptor for integrating stages into SquadRuleset
func (r *SquadRuleset) callStageFunc(stage StageFunc, boardState *BoardState, moves []SnakeMove) (bool, error) {
	settings, err := r.getSettingsJSON()
	if err != nil {
		return false, err
	}
	return stage(boardState, settings, moves)
}
