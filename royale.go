package rules

import (
	"errors"
	"math/rand"
)

type RoyaleRuleset struct {
	StandardRuleset

	Seed int64

	ShrinkEveryNTurns int32
}

func (r *RoyaleRuleset) Name() string { return GameTypeRoyale }

func (r RoyaleRuleset) Pipeline() (*Pipeline, error) {
	return NewPipeline(
		"snake.movement.standard",
		"health.reduce.standard",
		"hazard.damage.standard",
		"snake.eatfood.standard",
		"food.spawn.standard",
		"snake.eliminate.standard",
		"hazard.spawn.royale",
		"gameover.standard",
	)
}

func (r *RoyaleRuleset) CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
	if r.StandardRuleset.HazardDamagePerTurn < 1 {
		return nil, errors.New("royale damage per turn must be greater than zero")
	}

	p, err := r.Pipeline()
	if err != nil {
		return nil, err
	}
	_, nextState, err := p.Execute(prevState, r.Settings(), moves)

	return nextState, err
}

func PopulateHazardsRoyale(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
	if IsInitialisation(b, settings, moves) {
		return false, nil
	}
	b.Hazards = []Point{}

	// Royale uses the current turn to generate hazards, not the previous turn that's in the board state
	turn := b.Turn + 1

	if settings.RoyaleSettings.ShrinkEveryNTurns < 1 {
		return false, errors.New("royale game can't shrink more frequently than every turn")
	}

	if turn < settings.RoyaleSettings.ShrinkEveryNTurns {
		return false, nil
	}

	randGenerator := rand.New(rand.NewSource(settings.RoyaleSettings.seed))

	numShrinks := turn / settings.RoyaleSettings.ShrinkEveryNTurns
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
				b.Hazards = append(b.Hazards, Point{x, y})
			}
		}
	}

	return false, nil
}

func (r RoyaleRuleset) Settings() Settings {
	s := r.StandardRuleset.Settings()
	s.RoyaleSettings = RoyaleSettings{
		seed:              r.Seed,
		ShrinkEveryNTurns: r.ShrinkEveryNTurns,
	}
	return s
}
