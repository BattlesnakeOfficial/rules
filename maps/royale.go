package maps

import (
	"errors"

	"github.com/BattlesnakeOfficial/rules"
)

type RoyaleHazardsMap struct{}

func init() {
	globalRegistry.RegisterMap("royale", RoyaleHazardsMap{})
}

func (m RoyaleHazardsMap) ID() string {
	return "royale"
}

func (m RoyaleHazardsMap) Meta() Metadata {
	return Metadata{
		Name:        "Royale",
		Description: "A map where hazards are generated every N turns",
		Author:      "Battlesnake",
		Version:     2,
		MinPlayers:  1,
		MaxPlayers:  16,
		BoardSizes:  OddSizes(rules.BoardSizeSmall, rules.BoardSizeXXLarge),
		Tags:        []string{TAG_HAZARD_PLACEMENT},
	}
}

func (m RoyaleHazardsMap) SetupBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return StandardMap{}.SetupBoard(lastBoardState, settings, editor)
}

func (m RoyaleHazardsMap) PreUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return nil
}

func (m RoyaleHazardsMap) PostUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	// Use StandardMap to populate food
	if err := (StandardMap{}).PostUpdateBoard(lastBoardState, settings, editor); err != nil {
		return err
	}

	// Royale uses the current turn to generate hazards, not the previous turn that's in the board state
	turn := lastBoardState.Turn + 1

	shrinkEveryNTurns := settings.Int(rules.ParamShrinkEveryNTurns, 20)
	if shrinkEveryNTurns < 1 {
		return errors.New("royale game can't shrink more frequently than every turn")
	}

	if turn < shrinkEveryNTurns {
		return nil
	}

	// Reset hazards every turn and re-generate them
	editor.ClearHazards()

	// Get random generator for turn zero, because we're regenerating all hazards every time.
	randGenerator := settings.GetRand(0)

	numShrinks := turn / shrinkEveryNTurns
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
