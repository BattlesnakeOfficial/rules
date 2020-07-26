package rules

import (
	"errors"
)

type RoyaleRuleset struct {
	StandardRuleset

	Turn              int32
	ShrinkEveryNTurns int32

	// Output
	OutOfBounds []Point
}

func (r *RoyaleRuleset) CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
	if r.ShrinkEveryNTurns < 1 {
		return nil, errors.New("royale game must shrink at least every 1 turn")
	}

	nextBoardState, err := r.StandardRuleset.CreateNextBoardState(prevState, moves)
	if err != nil {
		return nil, err
	}

	// TODO: LOG?
	err = r.populateOutOfBounds(nextBoardState)
	if err != nil {
		return nil, err
	}

	// TODO: LOG?
	err = r.eliminateOutOfBounds(nextBoardState)
	if err != nil {
		return nil, err
	}

	return nextBoardState, nil
}

func (r *RoyaleRuleset) populateOutOfBounds(b *BoardState) error {
	r.OutOfBounds = []Point{}

	if r.ShrinkEveryNTurns < 1 {
		return errors.New("royale game must shrink at least every 1 turn")
	}

	if r.Turn < r.ShrinkEveryNTurns {
		return nil
	}

	numShrinks := r.Turn / r.ShrinkEveryNTurns
	minX, maxX := numShrinks, b.Width-1-numShrinks
	minY, maxY := numShrinks, b.Height-1-numShrinks
	for x := int32(0); x < b.Width; x++ {
		for y := int32(0); y < b.Height; y++ {
			if x < minX || x > maxX || y < minY || y > maxY {
				r.OutOfBounds = append(r.OutOfBounds, Point{x, y})
			}
		}
	}

	return nil
}

func (r *RoyaleRuleset) eliminateOutOfBounds(b *BoardState) error {
	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause == NotEliminated {
			head := snake.Body[0]
			for _, p := range r.OutOfBounds {
				if head == p {
					// Snake is now out of bounds, eliminate it
					snake.EliminatedCause = EliminatedByOutOfBounds
				}
			}
		}
	}

	return nil
}
