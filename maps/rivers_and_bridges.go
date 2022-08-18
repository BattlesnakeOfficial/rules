package maps

import (
	"github.com/BattlesnakeOfficial/rules"
)

func init() {
	globalRegistry.RegisterMap("hz_rivers_bridges", RiverAndBridgesMediumHazardsMap{})
	globalRegistry.RegisterMap("hz_rivers_bridges_lg", RiverAndBridgesLargeHazardsMap{})
	globalRegistry.RegisterMap("hz_rivers_bridges_xl", RiverAndBridgesExtraLargeHazardsMap{})
}

func setupRiverAndBridgesBoard(startingPositions [][]rules.Point, hazards []rules.Point, lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	width := lastBoardState.Width
	height := lastBoardState.Height

	numSnakes := len(lastBoardState.Snakes)
	if numSnakes == 0 {
		return rules.RulesetError("too few snakes - at least one snake must be present")
	}

	rand := settings.GetRand(0)

	snakeIDs := make([]string, 0, len(lastBoardState.Snakes))
	for _, snake := range lastBoardState.Snakes {
		snakeIDs = append(snakeIDs, snake.ID)
	}

	tempBoardState := rules.NewBoardState(width, height)
	tempBoardState.Snakes = make([]rules.Snake, len(snakeIDs))

	for i := 0; i < len(snakeIDs); i++ {
		tempBoardState.Snakes[i] = rules.Snake{
			ID:     snakeIDs[i],
			Health: rules.SnakeMaxHealth,
		}
	}
	err := rules.PlaceSnakesInQuadrants(rand, tempBoardState, startingPositions)
	if err != nil {
		return err
	}

	err = rules.PlaceFoodFixed(rand, tempBoardState)
	if err != nil {
		return err
	}

	// Copy food from temp board state
	for _, f := range tempBoardState.Food {
		// skip the center food
		if f.X == lastBoardState.Width/2 && f.Y == lastBoardState.Height/2 {
			continue
		}
		editor.AddFood(f)
	}

	// Copy snakes from temp board state
	for _, snake := range tempBoardState.Snakes {
		editor.PlaceSnake(snake.ID, snake.Body, snake.Health)
	}

	for _, p := range hazards {
		editor.AddHazard(p)
	}

	return nil
}

func placeRiverAndBridgesFood(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	rand := settings.GetRand(lastBoardState.Turn)

	foodNeeded := checkFoodNeedingPlacement(rand, settings, lastBoardState)
	if foodNeeded > 0 {
		pts := rules.GetUnoccupiedPoints(lastBoardState, false, true)
		placeFoodRandomlyAtPositions(rand, lastBoardState, editor, foodNeeded, pts)
	}

	return nil
}

type RiverAndBridgesMediumHazardsMap struct{}

func (m RiverAndBridgesMediumHazardsMap) ID() string {
	return "hz_rivers_bridges"
}

func (m RiverAndBridgesMediumHazardsMap) Meta() Metadata {
	return Metadata{
		Name: "hz_rivers_bridges",
		Description: `Creates fixed maps that have a lake of hazard in the middle with rivers going in the cardinal directions.
Each river has one or two 1-square "bridges" over them`,
		Author:     "Battlesnake",
		Version:    1,
		MinPlayers: 1,
		MaxPlayers: 8,
		BoardSizes: FixedSizes(Dimensions{11, 11}),
		Tags:       []string{TAG_FOOD_PLACEMENT, TAG_HAZARD_PLACEMENT, TAG_SNAKE_PLACEMENT},
	}
}

func (m RiverAndBridgesMediumHazardsMap) SetupBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	if !m.Meta().BoardSizes.IsAllowable(uint(lastBoardState.Width), uint(lastBoardState.Height)) {
		return rules.RulesetError("This map can only be played on a 11x11 board")
	}

	err := setupRiverAndBridgesBoard(riversAndBridgesMediumStartPositions, riversAndBridgesMediumHazards, lastBoardState, settings, editor)
	if err != nil {
		return err
	}
	return nil
}

func (m RiverAndBridgesMediumHazardsMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return placeRiverAndBridgesFood(lastBoardState, settings, editor)
}

var riversAndBridgesMediumStartPositions = [][]rules.Point{
	{
		{X: 1, Y: 1},
		{X: 3, Y: 3},
	},
	{
		{X: 9, Y: 9},
		{X: 7, Y: 7},
	},
	{
		{X: 1, Y: 9},
		{X: 3, Y: 9},
	},
	{
		{X: 9, Y: 1},
		{X: 7, Y: 3},
	},
}

var riversAndBridgesMediumHazards = []rules.Point{
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
}

type RiverAndBridgesLargeHazardsMap struct{}

func (m RiverAndBridgesLargeHazardsMap) ID() string {
	return "hz_rivers_bridges_lg"
}

func (m RiverAndBridgesLargeHazardsMap) Meta() Metadata {
	return Metadata{
		Name: "hz_rivers_bridges_lg",
		Description: `Creates fixed maps that have a lake of hazard in the middle with rivers going in the cardinal directions.
Each river has one or two 1-square "bridges" over them`,
		Author:     "Battlesnake",
		Version:    1,
		MinPlayers: 1,
		MaxPlayers: 12,
		BoardSizes: FixedSizes(Dimensions{19, 19}),
		Tags:       []string{TAG_FOOD_PLACEMENT, TAG_HAZARD_PLACEMENT, TAG_SNAKE_PLACEMENT},
	}
}

func (m RiverAndBridgesLargeHazardsMap) SetupBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	if !m.Meta().BoardSizes.IsAllowable(uint(lastBoardState.Width), uint(lastBoardState.Height)) {
		return rules.RulesetError("This map can only be played on a 19x19 board")
	}

	err := setupRiverAndBridgesBoard(riversAndBridgesLargeStartPositions, riversAndBridgesLargeHazards, lastBoardState, settings, editor)
	if err != nil {
		return err
	}
	return nil
}

func (m RiverAndBridgesLargeHazardsMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return placeRiverAndBridgesFood(lastBoardState, settings, editor)
}

var riversAndBridgesLargeStartPositions = [][]rules.Point{
	{
		{X: 1, Y: 1},
		{X: 5, Y: 1},
		{X: 1, Y: 5},
		{X: 5, Y: 5},
	},
	{
		{X: 17, Y: 1},
		{X: 17, Y: 5},
		{X: 13, Y: 5},
		{X: 13, Y: 1},
	},
	{
		{X: 1, Y: 17},
		{X: 5, Y: 17},
		{X: 1, Y: 13},
		{X: 5, Y: 13},
	},
	{
		{X: 17, Y: 17},
		{X: 17, Y: 13},
		{X: 13, Y: 17},
		{X: 13, Y: 13},
	},
}

var riversAndBridgesLargeHazards = []rules.Point{
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
}

type RiverAndBridgesExtraLargeHazardsMap struct{}

func (m RiverAndBridgesExtraLargeHazardsMap) ID() string {
	return "hz_rivers_bridges_xl"
}

func (m RiverAndBridgesExtraLargeHazardsMap) Meta() Metadata {
	return Metadata{
		Name: "hz_rivers_bridges_xl",
		Description: `Creates fixed maps that have a lake of hazard in the middle with rivers going in the cardinal directions.
Each river has one or two 1-square "bridges" over them`,
		Author:     "Battlesnake",
		Version:    1,
		MinPlayers: 1,
		MaxPlayers: 12,
		BoardSizes: FixedSizes(Dimensions{25, 25}),
		Tags:       []string{TAG_FOOD_PLACEMENT, TAG_HAZARD_PLACEMENT, TAG_SNAKE_PLACEMENT},
	}
}

func (m RiverAndBridgesExtraLargeHazardsMap) SetupBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	if !m.Meta().BoardSizes.IsAllowable(uint(lastBoardState.Width), uint(lastBoardState.Height)) {
		return rules.RulesetError("This map can only be played on a 25x25 board")
	}

	err := setupRiverAndBridgesBoard(riversAndBridgesExtraLargeStartPositions, riversAndBridgesExtraLargeHazards, lastBoardState, settings, editor)
	if err != nil {
		return err
	}
	return nil
}

func (m RiverAndBridgesExtraLargeHazardsMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return placeRiverAndBridgesFood(lastBoardState, settings, editor)
}

var riversAndBridgesExtraLargeStartPositions = [][]rules.Point{
	{
		{X: 1, Y: 1},
		{X: 9, Y: 9},
		{X: 9, Y: 1},
		{X: 1, Y: 9},
	},
	{
		{X: 23, Y: 23},
		{X: 15, Y: 15},
		{X: 23, Y: 15},
		{X: 15, Y: 23},
	},
	{
		{X: 15, Y: 1},
		{X: 15, Y: 9},
		{X: 23, Y: 9},
		{X: 23, Y: 1},
	},
	{
		{X: 9, Y: 23},
		{X: 1, Y: 23},
		{X: 9, Y: 15},
		{X: 1, Y: 15},
	},
}

var riversAndBridgesExtraLargeHazards = []rules.Point{
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
}
