package rules

import (
	"errors"
	"math/rand"
)

type RoyaleRuleset struct {
	StandardRuleset

	Seed int64

	Turn              int32
	ShrinkEveryNTurns int32
	DamagePerTurn     int32

	// Output
	OutOfBounds []Point
}

func (r *RoyaleRuleset) CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
	if r.ShrinkEveryNTurns < 1 {
		return nil, errors.New("royale game must shrink at least every turn")
	}

	nextBoardState, err := r.StandardRuleset.CreateNextBoardState(prevState, moves)
	if err != nil {
		return nil, err
	}

	// Algorithm:
	// - Populate OOB for last turn
	// - Apply damage to snake heads that are OOB
	// - Re-populate OOB for this turn
	// ---> This means damage on board shrinks doesn't hit until the following turn.

	// TODO: LOG?
	err = r.populateOutOfBounds(nextBoardState, r.Turn-1)
	if err != nil {
		return nil, err
	}

	// TODO: LOG?
	err = r.damageOutOfBounds(nextBoardState)
	if err != nil {
		return nil, err
	}

	// TODO: LOG?
	err = r.populateOutOfBounds(nextBoardState, r.Turn)
	if err != nil {
		return nil, err
	}

	return nextBoardState, nil
}

func (r *RoyaleRuleset) Name() string { return "royale" }

func (r *RoyaleRuleset) Version() string { return "1.0.0" }

func (r *RoyaleRuleset) populateOutOfBounds(b *BoardState, turn int32) error {
	r.OutOfBounds = []Point{}

	if r.ShrinkEveryNTurns < 1 {
		return errors.New("royale game must shrink at least every turn")
	}

	if turn < r.ShrinkEveryNTurns {
		return nil
	}

	randGenerator := rand.New(rand.NewSource(r.Seed))

	numShrinks := turn / r.ShrinkEveryNTurns
	minX, maxX := int32(0), b.Width-1
	minY, maxY := int32(0), b.Height-1
	for i := int32(0); i < numShrinks; i++ {
		switch randGenerator.Intn(4) {
		case 0:
			minX += 1
		case 1:
			maxX -= 1
		case 2:
			minY += 1
		case 3:
			maxY -= 1
		}
	}

	for x := int32(0); x < b.Width; x++ {
		for y := int32(0); y < b.Height; y++ {
			if x < minX || x > maxX || y < minY || y > maxY {
				r.OutOfBounds = append(r.OutOfBounds, Point{x, y})
			}
		}
	}

	return nil
}

func (r *RoyaleRuleset) damageOutOfBounds(b *BoardState) error {
	if r.DamagePerTurn < 1 {
		return errors.New("royale damage per turn must be greater than zero")
	}

	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause == NotEliminated {
			head := snake.Body[0]
			for _, p := range r.OutOfBounds {
				if head == p {
					// Snake is now out of bounds, reduce health
					snake.Health = snake.Health - r.DamagePerTurn
					if snake.Health < 0 {
						snake.Health = 0
					}
					if r.StandardRuleset.snakeIsOutOfHealth(snake) {
						snake.EliminatedCause = EliminatedByOutOfHealth
					}
				}
			}
		}
	}

	return nil
}
