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
