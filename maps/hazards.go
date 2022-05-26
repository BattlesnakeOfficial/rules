package maps

import (
	"github.com/BattlesnakeOfficial/rules"
)

type InnerBorderHazardsMap struct{}

func init() {
	globalRegistry.RegisterMap("hz_inner_wall", RoyaleHazardsMap{})
}

func (m InnerBorderHazardsMap) ID() string {
	return "hz_inner_wall"
}

func (m InnerBorderHazardsMap) Meta() Metadata {
	return Metadata{
		Name:        "hz_inner_wall",
		Description: "Creates a static map on turn 0 that is a 1-square wall of hazard that is inset 2 squares from the edge of the board",
		Author:      "Battlesnake",
	}
}

func (m InnerBorderHazardsMap) SetupBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	if err := (StandardMap{}).SetupBoard(lastBoardState, settings, editor); err != nil {
		return err
	}

	// draw the initial, single ring of hazards
	hazards := drawRing(lastBoardState.Width, lastBoardState.Height, 2, 2)
	for _, p := range hazards {
		editor.AddHazard(p)
	}

	return nil
}

func (m InnerBorderHazardsMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return StandardMap{}.UpdateBoard(lastBoardState, settings, editor)
}

type ConcentricRingsHazardsMap struct{}

func init() {
	globalRegistry.RegisterMap("hz_rings", RoyaleHazardsMap{})
}

func (m ConcentricRingsHazardsMap) ID() string {
	return "hz_rings"
}

func (m ConcentricRingsHazardsMap) Meta() Metadata {
	return Metadata{
		Name:        "hz_rings",
		Description: "Creates a static map where there are rings of hazard sauce starting from the center with a 1 square space between the rings that has no sauce",
		Author:      "Battlesnake",
	}
}

func (m ConcentricRingsHazardsMap) SetupBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	if err := (StandardMap{}).SetupBoard(lastBoardState, settings, editor); err != nil {
		return err
	}

	// draw concentric rings of hazards
	for offset := 2; offset < lastBoardState.Width/2; offset += 2 {
		hazards := drawRing(lastBoardState.Width, lastBoardState.Height, offset, offset)
		for _, p := range hazards {
			editor.AddHazard(p)
		}
	}

	return nil
}

func (m ConcentricRingsHazardsMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return StandardMap{}.UpdateBoard(lastBoardState, settings, editor)
}
