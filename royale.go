package rules

import (
	"errors"
	"math/rand"

	"github.com/tidwall/sjson"
)

type RoyaleRuleset struct {
	StandardRuleset

	Seed int64

	ShrinkEveryNTurns int32
}

func (r *RoyaleRuleset) Name() string { return "royale" }

func (r *RoyaleRuleset) CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
	if r.StandardRuleset.HazardDamagePerTurn < 1 {
		return nil, errors.New("royale damage per turn must be greater than zero")
	}

	nextBoardState, err := r.StandardRuleset.CreateNextBoardState(prevState, moves)
	if err != nil {
		return nil, err
	}

	// Royale's only job is now to populate the hazards for next turn - StandardRuleset takes care of applying hazard damage.
	err = r.populateHazards(nextBoardState)
	if err != nil {
		return nil, err
	}

	return nextBoardState, nil
}

func (r *RoyaleRuleset) populateHazards(b *BoardState) error {
	_, err := r.callStageFunc(PopulateHazardsRoyale, b, []SnakeMove{})
	return err
}

func PopulateHazardsRoyale(b *BoardState, settings SettingsJSON, moves []SnakeMove) (bool, error) {
	b.Hazards = []Point{}
	shrinkEveryNTurns := settings.GetInt32("royale", "shrinkEveryNTurns")

	// Royale uses the current turn to generate hazards, not the previous turn that's in the board state
	turn := b.Turn + 1

	if shrinkEveryNTurns < 1 {
		return false, errors.New("royale game can't shrink more frequently than every turn")
	}

	if turn < shrinkEveryNTurns {
		return false, nil
	}

	randGenerator := rand.New(rand.NewSource(settings.GetInt64("royale", "seed")))

	numShrinks := turn / shrinkEveryNTurns
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

func (r *RoyaleRuleset) getSettingsJSON() (SettingsJSON, error) {
	j, err := r.StandardRuleset.getSettingsJSON()
	if err != nil {
		return nil, err
	}

	// apply royale public API settings
	j, err = sjson.SetBytes(j, "royale", RoyaleSettings{
		ShrinkEveryNTurns: r.ShrinkEveryNTurns,
	})
	if err != nil {
		return nil, err
	}

	// patch in seed value
	j, err = sjson.SetBytes(j, "royale.seed", r.Seed)

	return j, err
}

// Adaptor for integrating stages into RoyaleRuleset
func (r *RoyaleRuleset) callStageFunc(stage StageFunc, boardState *BoardState, moves []SnakeMove) (bool, error) {
	settings, err := r.getSettingsJSON()
	if err != nil {
		return false, err
	}
	return stage(boardState, settings, moves)
}
