package rules

import (
	"errors"
)

var royaleRulesetStages = []string{
	StageGameOverStandard,
	StageMovementStandard,
	StageStarvationStandard,
	StageHazardDamageStandard,
	StageFeedSnakesStandard,
	StageEliminationStandard,
	StageSpawnHazardsShrinkMap,
}

func PopulateHazardsRoyale(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
	if IsInitialization(b, settings, moves) {
		return false, nil
	}
	b.Hazards = []Point{}

	// Royale uses the current turn to generate hazards, not the previous turn that's in the board state
	turn := b.Turn + 1

	shrinkEveryNTurns := settings.Int(ParamShrinkEveryNTurns, 20)
	if shrinkEveryNTurns < 1 {
		return false, errors.New("royale game can't shrink more frequently than every turn")
	}

	if turn < shrinkEveryNTurns {
		return false, nil
	}

	randGenerator := settings.GetRand(0)

	numShrinks := turn / shrinkEveryNTurns
	minX, maxX := 0, b.Width-1
	minY, maxY := 0, b.Height-1
	for i := 0; i < numShrinks; i++ {
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

	for x := 0; x < b.Width; x++ {
		for y := 0; y < b.Height; y++ {
			if x < minX || x > maxX || y < minY || y > maxY {
				b.Hazards = append(b.Hazards, Point{X: x, Y: y})
			}
		}
	}

	return false, nil
}
