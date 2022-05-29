package maps

import (
	"errors"

	"github.com/BattlesnakeOfficial/rules"
)

type RoyaleHazardsMap struct{}

func init() {
	globalRegistry.RegisterMap(RoyaleHazardsMap{})
}

func (m RoyaleHazardsMap) ID() string {
	return "royale"
}

func (m RoyaleHazardsMap) Meta() Metadata {
	return Metadata{
		Name:        "Royale",
		Description: "A map where hazards are generated every N turns",
		Author:      "Battlesnake",
	}
}

func (m RoyaleHazardsMap) SetupBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return StandardMap{}.SetupBoard(lastBoardState, settings, editor)
}

func (m RoyaleHazardsMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	// Use StandardMap to populate food
	if err := (StandardMap{}).UpdateBoard(lastBoardState, settings, editor); err != nil {
		return err
	}

	// Royale uses the current turn to generate hazards, not the previous turn that's in the board state
	turn := lastBoardState.Turn + 1

	if settings.RoyaleSettings.ShrinkEveryNTurns < 1 {
		return errors.New("royale game can't shrink more frequently than every turn")
	}

	if turn < settings.RoyaleSettings.ShrinkEveryNTurns {
		return nil
	}

	// Reset hazards every turn and re-generate them
	editor.ClearHazards()

	// Get random generator for turn zero, because we're regenerating all hazards every time.
	randGenerator := settings.GetRand(0)

	numShrinks := turn / settings.RoyaleSettings.ShrinkEveryNTurns
	minX, maxX := 0, lastBoardState.Width-1
	minY, maxY := 0, lastBoardState.Height-1
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

	for x := 0; x < lastBoardState.Width; x++ {
		for y := 0; y < lastBoardState.Height; y++ {
			if x < minX || x > maxX || y < minY || y > maxY {
				editor.AddHazard(rules.Point{X: x, Y: y})
			}
		}
	}

	return nil
}
