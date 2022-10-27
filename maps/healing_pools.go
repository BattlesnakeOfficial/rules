package maps

import (
	"math/rand"

	"github.com/BattlesnakeOfficial/rules"
)

type HealingPoolsMap struct{}

func init() {
	globalRegistry.RegisterMap("healing_pools", HealingPoolsMap{})
}

func (m HealingPoolsMap) ID() string {
	return "healing_pools"
}

func (m HealingPoolsMap) Meta() Metadata {
	return Metadata{
		Name:        "Healing Pools",
		Description: "A simple map that spawns fixed single cell hazard areas based on the map size.",
		Author:      "Battlesnake",
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  8,
		BoardSizes:  FixedSizes(Dimensions{7, 7}, Dimensions{11, 11}, Dimensions{19, 19}),
		Tags:        []string{TAG_HAZARD_PLACEMENT},
	}
}

func (m HealingPoolsMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	if err := (StandardMap{}).SetupBoard(initialBoardState, settings, editor); err != nil {
		return err
	}

	rand := settings.GetRand(0)

	options, ok := poolLocationOptions[rules.Point{X: initialBoardState.Width, Y: initialBoardState.Height}]
	if !ok {
		return rules.RulesetError("board size is not supported by this map")
	}

	i := rand.Intn(len(options))

	for _, p := range options[i] {
		editor.AddHazard(p)
	}

	return nil
}

func (m HealingPoolsMap) PreUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return nil
}

func (m HealingPoolsMap) PostUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	if err := (StandardMap{}).PostUpdateBoard(lastBoardState, settings, editor); err != nil {
		return err
	}

	shrinkEveryNTurns := settings.Int(rules.ParamShrinkEveryNTurns, 0)
	if lastBoardState.Turn > 0 && shrinkEveryNTurns > 0 && len(lastBoardState.Hazards) > 0 && lastBoardState.Turn%shrinkEveryNTurns == 0 {
		// Attempt to remove a healing pool every ShrinkEveryNTurns until there are none remaining
		i := rand.Intn(len(lastBoardState.Hazards))
		editor.RemoveHazard(lastBoardState.Hazards[i])
	}

	return nil
}

var poolLocationOptions = map[rules.Point][][]rules.Point{
	{X: 7, Y: 7}: {
		{
			{X: 3, Y: 3},
		},
	},
	{X: 11, Y: 11}: {
		{
			{X: 3, Y: 3},
			{X: 7, Y: 7},
		},
		{
			{X: 3, Y: 7},
			{X: 7, Y: 3},
		},
		{
			{X: 3, Y: 5},
			{X: 7, Y: 5},
		},
		{
			{X: 5, Y: 7},
			{X: 5, Y: 3},
		},
	},
	{X: 19, Y: 19}: {
		{
			{X: 5, Y: 5},
			{X: 13, Y: 13},
			{X: 5, Y: 13},
			{X: 13, Y: 5},
		},
		{
			{X: 5, Y: 10},
			{X: 13, Y: 10},
			{X: 10, Y: 13},
			{X: 10, Y: 5},
		},
	},
}
