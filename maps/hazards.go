package maps

import (
	"errors"
	"fmt"

	"github.com/BattlesnakeOfficial/rules"
)

type InnerBorderHazardsMap struct{}

func init() {
	globalRegistry.RegisterMap("hz_inner_wall", InnerBorderHazardsMap{})
	globalRegistry.RegisterMap("hz_rings", ConcentricRingsHazardsMap{})
	globalRegistry.RegisterMap("hz_columns", ColumnsHazardsMap{})
	globalRegistry.RegisterMap("hz_rivers_bridges", RiverAndBridgesHazardsMap{})
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
	hazards, err := drawRing(lastBoardState.Width, lastBoardState.Height, 2, 2)
	if err != nil {
		return err
	}

	for _, p := range hazards {
		editor.AddHazard(p)
	}

	return nil
}

func (m InnerBorderHazardsMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return StandardMap{}.UpdateBoard(lastBoardState, settings, editor)
}

type ConcentricRingsHazardsMap struct{}

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
		hazards, err := drawRing(lastBoardState.Width, lastBoardState.Height, offset, offset)
		if err != nil {
			return err
		}
		for _, p := range hazards {
			editor.AddHazard(p)
		}
	}

	return nil
}

func (m ConcentricRingsHazardsMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return StandardMap{}.UpdateBoard(lastBoardState, settings, editor)
}

type ColumnsHazardsMap struct{}

func (m ColumnsHazardsMap) ID() string {
	return "hz_columns"
}

func (m ColumnsHazardsMap) Meta() Metadata {
	return Metadata{
		Name:        "hz_columns",
		Description: "Creates a static map on turn 0 that fills in odd squares, i.e. (1,1), (1,3), (3,3) ... with hazard sauce",
		Author:      "Battlesnake",
	}
}

func (m ColumnsHazardsMap) SetupBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	if err := (StandardMap{}).SetupBoard(lastBoardState, settings, editor); err != nil {
		return err
	}

	for x := 0; x < lastBoardState.Width; x++ {
		for y := 0; y < lastBoardState.Height; y++ {
			if x%2 == 1 && y%2 == 1 {
				editor.AddHazard(rules.Point{X: x, Y: y})
			}
		}
	}

	return nil
}

func (m ColumnsHazardsMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return StandardMap{}.UpdateBoard(lastBoardState, settings, editor)
}

type RiverAndBridgesHazardsMap struct{}

func (m RiverAndBridgesHazardsMap) ID() string {
	return "hz_rivers_bridges"
}

func (m RiverAndBridgesHazardsMap) Meta() Metadata {
	return Metadata{
		Name: "hz_rivers_bridges",
		Description: `Creates fixed maps that have a lake of hazard in the middle with rivers going in the cardinal directions.
Each river has one or two 1-square "bridges" over them`,
		Author: "Battlesnake",
	}
}

func (m RiverAndBridgesHazardsMap) SetupBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	if err := (StandardMap{}).SetupBoard(lastBoardState, settings, editor); err != nil {
		return err
	}

	key := fmt.Sprintf("%dx%d", lastBoardState.Width, lastBoardState.Height)
	hazards, ok := riversAndBridgesMaps[key]
	if !ok {
		return errors.New("Board size is not supported by this map")
	}
	for _, p := range hazards {
		editor.AddHazard(p)
	}

	return nil
}

func (m RiverAndBridgesHazardsMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return StandardMap{}.UpdateBoard(lastBoardState, settings, editor)
}

var riversAndBridgesMaps = map[string][]rules.Point{
	"11x11": {
		{X: 5, Y: 10},
		{X: 5, Y: 9},
		{X: 5, Y: 7},
		{X: 5, Y: 6},
		{X: 5, Y: 5},
		{X: 5, Y: 4},
		{X: 5, Y: 3},
		{X: 5, Y: 0},
		{X: 5, Y: 1},
		{X: 6, Y: 5},
		{X: 7, Y: 5},
		{X: 9, Y: 5},
		{X: 10, Y: 5},
		{X: 4, Y: 5},
		{X: 3, Y: 5},
		{X: 1, Y: 5},
		{X: 0, Y: 5},
	},
	"19x19": {
		{X: 9, Y: 0},
		{X: 9, Y: 1},
		{X: 9, Y: 2},
		{X: 9, Y: 5},
		{X: 9, Y: 6},
		{X: 9, Y: 7},
		{X: 9, Y: 9},
		{X: 9, Y: 8},
		{X: 9, Y: 10},
		{X: 9, Y: 12},
		{X: 9, Y: 11},
		{X: 9, Y: 13},
		{X: 9, Y: 14},
		{X: 9, Y: 16},
		{X: 9, Y: 17},
		{X: 9, Y: 18},
		{X: 0, Y: 9},
		{X: 2, Y: 9},
		{X: 1, Y: 9},
		{X: 3, Y: 9},
		{X: 5, Y: 9},
		{X: 6, Y: 9},
		{X: 7, Y: 9},
		{X: 8, Y: 9},
		{X: 10, Y: 9},
		{X: 13, Y: 9},
		{X: 12, Y: 9},
		{X: 11, Y: 9},
		{X: 15, Y: 9},
		{X: 16, Y: 9},
		{X: 17, Y: 9},
		{X: 18, Y: 9},
		{X: 9, Y: 4},
		{X: 8, Y: 10},
		{X: 8, Y: 8},
		{X: 10, Y: 8},
		{X: 10, Y: 10},
	},
	"25x25": {
		{X: 12, Y: 24},
		{X: 12, Y: 21},
		{X: 12, Y: 20},
		{X: 12, Y: 19},
		{X: 12, Y: 18},
		{X: 12, Y: 15},
		{X: 12, Y: 14},
		{X: 12, Y: 13},
		{X: 12, Y: 12},
		{X: 12, Y: 11},
		{X: 12, Y: 10},
		{X: 12, Y: 9},
		{X: 12, Y: 5},
		{X: 12, Y: 4},
		{X: 12, Y: 3},
		{X: 12, Y: 0},
		{X: 0, Y: 12},
		{X: 3, Y: 12},
		{X: 4, Y: 12},
		{X: 5, Y: 12},
		{X: 6, Y: 12},
		{X: 9, Y: 12},
		{X: 10, Y: 12},
		{X: 11, Y: 12},
		{X: 13, Y: 12},
		{X: 14, Y: 12},
		{X: 15, Y: 12},
		{X: 18, Y: 12},
		{X: 20, Y: 12},
		{X: 19, Y: 12},
		{X: 21, Y: 12},
		{X: 24, Y: 12},
		{X: 11, Y: 14},
		{X: 10, Y: 13},
		{X: 11, Y: 13},
		{X: 10, Y: 11},
		{X: 11, Y: 11},
		{X: 11, Y: 10},
		{X: 13, Y: 10},
		{X: 14, Y: 11},
		{X: 13, Y: 11},
		{X: 13, Y: 13},
		{X: 14, Y: 13},
		{X: 13, Y: 14},
		{X: 12, Y: 6},
		{X: 12, Y: 2},
		{X: 2, Y: 12},
		{X: 22, Y: 12},
		{X: 12, Y: 22},
		{X: 16, Y: 12},
		{X: 12, Y: 8},
		{X: 8, Y: 12},
		{X: 12, Y: 16},
	},
}
