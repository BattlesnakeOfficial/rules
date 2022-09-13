package maps

import (
	"github.com/BattlesnakeOfficial/rules"
)

func init() {
	globalRegistry.RegisterMap("hz_castle_wall", CastleWallMediumHazardsMap{})
	globalRegistry.RegisterMap("hz_castle_wall_lg", CastleWallLargeHazardsMap{})
	globalRegistry.RegisterMap("hz_castle_wall_xl", CastleWallExtraLargeHazardsMap{})
}

func setupCastleWallBoard(maxPlayers int, startingPositions []rules.Point, hazards []rules.Point, initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	rand := settings.GetRand(initialBoardState.Turn)

	if len(initialBoardState.Snakes) > int(maxPlayers) {
		return rules.ErrorTooManySnakes
	}

	// place snakes
	rand.Shuffle(len(startingPositions), func(i int, j int) {
		startingPositions[i], startingPositions[j] = startingPositions[j], startingPositions[i]
	})
	for index, snake := range initialBoardState.Snakes {
		head := startingPositions[index]
		editor.PlaceSnake(snake.ID, []rules.Point{head, head, head}, rules.SnakeMaxHealth)
	}

	// place hazards
	for _, h := range hazards {
		editor.AddHazard(h)
	}
	return nil
}

func updateCastleWallBoard(maxFood int, food []rules.Point, lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	// no food spawning for first 10 turns
	if lastBoardState.Turn < 10 {
		return nil
	}

	// skip food spawn when max food present
	if len(lastBoardState.Food) == maxFood {
		return nil
	}

	rand := settings.GetRand(lastBoardState.Turn)

	rand.Shuffle(len(food), func(i int, j int) {
		food[i], food[j] = food[j], food[i]
	})

foodPlacementLoop:
	for _, f := range food {
		for _, snake := range lastBoardState.Snakes {
			for i, point := range snake.Body {
				if point.X == f.X && point.Y == f.Y {
					continue foodPlacementLoop
				}

				// also avoid spawning food next to a snake head
				if i == 0 {
					if point.X+1 == f.X && point.Y == f.Y {
						continue foodPlacementLoop
					}
					if point.X-1 == f.X && point.Y == f.Y {
						continue foodPlacementLoop
					}
					if point.X == f.X && point.Y+1 == f.Y {
						continue foodPlacementLoop
					}
					if point.X == f.X && point.Y-1 == f.Y {
						continue foodPlacementLoop
					}
				}
			}
		}

		for _, existingFood := range lastBoardState.Food {
			if existingFood.X == f.X && existingFood.Y == f.Y {
				continue foodPlacementLoop
			}

			// also avoid spawning food in same passage as existing food
			if existingFood.X+1 == f.X && existingFood.Y == f.Y {
				continue foodPlacementLoop
			}
			if existingFood.X-1 == f.X && existingFood.Y == f.Y {
				continue foodPlacementLoop
			}
			if existingFood.X == f.X && existingFood.Y+1 == f.Y {
				continue foodPlacementLoop
			}
			if existingFood.X == f.X && existingFood.Y-1 == f.Y {
				continue foodPlacementLoop
			}
		}

		editor.AddFood(f)
		break
	}

	return nil
}

type CastleWallMediumHazardsMap struct{}

func (m CastleWallMediumHazardsMap) ID() string {
	return "hz_castle_wall"
}

func (m CastleWallMediumHazardsMap) Meta() Metadata {
	return Metadata{
		Name:        "hz_castle_wall",
		Description: "Wall of hazards around the board with dangerous bridges",
		Author:      "bcambl",
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  8,
		BoardSizes:  FixedSizes(Dimensions{11, 11}),
		Tags:        []string{TAG_FOOD_PLACEMENT, TAG_HAZARD_PLACEMENT, TAG_SNAKE_PLACEMENT},
	}
}

func (m CastleWallMediumHazardsMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	if !m.Meta().BoardSizes.IsAllowable(initialBoardState.Width, initialBoardState.Height) {
		return rules.RulesetError("This map can only be played on a 11x11 board")
	}

	var startPositions []rules.Point
	startPositions = append(startPositions, castleWallMediumStartPositions[0]...)
	if len(initialBoardState.Snakes) >= 5 {
		startPositions = append(startPositions, castleWallMediumStartPositions[1]...)
	}

	return setupCastleWallBoard(m.Meta().MaxPlayers, startPositions, castleWallMediumHazards, initialBoardState, settings, editor)
}

func (m CastleWallMediumHazardsMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	maxFood := 2
	return updateCastleWallBoard(maxFood, castleWallMediumFood, lastBoardState, settings, editor)
}

var castleWallMediumStartPositions = [][]rules.Point{
	{
		{X: 1, Y: 1},
		{X: 1, Y: 9},
		{X: 9, Y: 1},
		{X: 9, Y: 9},
	},
	{
		{X: 1, Y: 5},
		{X: 5, Y: 1},
		{X: 5, Y: 9},
		{X: 9, Y: 5},
	},
}

var castleWallMediumFood = []rules.Point{
	{X: 2, Y: 5},
	{X: 5, Y: 2},
	{X: 5, Y: 8},
	{X: 8, Y: 5},
}

var castleWallMediumHazards = []rules.Point{
	{X: 2, Y: 2},
	{X: 2, Y: 3},
	{X: 2, Y: 4},
	{X: 2, Y: 6},
	{X: 2, Y: 7},
	{X: 2, Y: 8},
	{X: 3, Y: 2},
	{X: 3, Y: 8},
	{X: 4, Y: 2},
	{X: 4, Y: 8},
	{X: 6, Y: 2},
	{X: 6, Y: 8},
	{X: 7, Y: 2},
	{X: 7, Y: 8},
	{X: 8, Y: 2},
	{X: 8, Y: 3},
	{X: 8, Y: 4},
	{X: 8, Y: 6},
	{X: 8, Y: 7},
	{X: 8, Y: 8},
	// double hazards near passages:
	{X: 2, Y: 4},
	{X: 2, Y: 6},
	{X: 4, Y: 2},
	{X: 4, Y: 8},
	{X: 6, Y: 2},
	{X: 6, Y: 8},
	{X: 8, Y: 4},
	{X: 8, Y: 6},
}

type CastleWallLargeHazardsMap struct{}

func (m CastleWallLargeHazardsMap) ID() string {
	return "hz_castle_wall_lg"
}

func (m CastleWallLargeHazardsMap) Meta() Metadata {
	return Metadata{
		Name:        "hz_castle_wall_lg",
		Description: "Wall of hazards around the board with dangerous bridges",
		Author:      "bcambl",
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  8,
		BoardSizes:  FixedSizes(Dimensions{19, 19}),
		Tags:        []string{TAG_FOOD_PLACEMENT, TAG_HAZARD_PLACEMENT, TAG_SNAKE_PLACEMENT},
	}
}

func (m CastleWallLargeHazardsMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	if !m.Meta().BoardSizes.IsAllowable(initialBoardState.Width, initialBoardState.Height) {
		return rules.RulesetError("This map can only be played on a 19x19 board")
	}

	var startPositions []rules.Point
	startPositions = append(startPositions, castleWallLargeStartPositions[0]...)
	if len(initialBoardState.Snakes) >= 5 {
		startPositions = append(startPositions, castleWallLargeStartPositions[1]...)
	}

	return setupCastleWallBoard(m.Meta().MaxPlayers, startPositions, castleWallLargeHazards, initialBoardState, settings, editor)
}

func (m CastleWallLargeHazardsMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	maxFood := 2
	return updateCastleWallBoard(maxFood, castleWallLargeFood, lastBoardState, settings, editor)
}

var castleWallLargeStartPositions = [][]rules.Point{
	{
		{X: 1, Y: 1},
		{X: 1, Y: 17},
		{X: 17, Y: 1},
		{X: 17, Y: 17},
	},
	{
		{X: 1, Y: 9},
		{X: 9, Y: 1},
		{X: 9, Y: 17},
		{X: 17, Y: 9},
	},
}

var castleWallLargeFood = []rules.Point{
	{X: 2, Y: 8},
	{X: 2, Y: 10},
	{X: 3, Y: 8},
	{X: 3, Y: 10},
	{X: 8, Y: 2},
	{X: 8, Y: 3},
	{X: 8, Y: 15},
	{X: 8, Y: 16},
	{X: 10, Y: 2},
	{X: 10, Y: 3},
	{X: 10, Y: 15},
	{X: 10, Y: 16},
	{X: 15, Y: 8},
	{X: 15, Y: 10},
	{X: 16, Y: 8},
	{X: 16, Y: 10},
}

var castleWallLargeHazards = []rules.Point{
	{X: 2, Y: 11},
	{X: 2, Y: 12},
	{X: 2, Y: 13},
	{X: 2, Y: 14},
	{X: 2, Y: 15},
	{X: 2, Y: 16},
	{X: 2, Y: 2},
	{X: 2, Y: 3},
	{X: 2, Y: 4},
	{X: 2, Y: 5},
	{X: 2, Y: 6},
	{X: 2, Y: 7},
	{X: 2, Y: 9},
	{X: 3, Y: 11},
	{X: 3, Y: 12},
	{X: 3, Y: 13},
	{X: 3, Y: 14},
	{X: 3, Y: 15},
	{X: 3, Y: 16},
	{X: 3, Y: 2},
	{X: 3, Y: 3},
	{X: 3, Y: 4},
	{X: 3, Y: 5},
	{X: 3, Y: 6},
	{X: 3, Y: 7},
	{X: 3, Y: 9},
	{X: 4, Y: 15},
	{X: 4, Y: 16},
	{X: 4, Y: 2},
	{X: 4, Y: 3},
	{X: 5, Y: 15},
	{X: 5, Y: 16},
	{X: 5, Y: 2},
	{X: 5, Y: 3},
	{X: 6, Y: 15},
	{X: 6, Y: 16},
	{X: 6, Y: 2},
	{X: 6, Y: 3},
	{X: 7, Y: 15},
	{X: 7, Y: 16},
	{X: 7, Y: 2},
	{X: 7, Y: 3},
	{X: 9, Y: 15},
	{X: 9, Y: 16},
	{X: 9, Y: 2},
	{X: 9, Y: 3},
	{X: 11, Y: 15},
	{X: 11, Y: 16},
	{X: 11, Y: 2},
	{X: 11, Y: 3},
	{X: 12, Y: 15},
	{X: 12, Y: 16},
	{X: 12, Y: 2},
	{X: 12, Y: 3},
	{X: 13, Y: 15},
	{X: 13, Y: 16},
	{X: 13, Y: 2},
	{X: 13, Y: 3},
	{X: 14, Y: 15},
	{X: 14, Y: 16},
	{X: 14, Y: 2},
	{X: 14, Y: 3},
	{X: 15, Y: 11},
	{X: 15, Y: 12},
	{X: 15, Y: 13},
	{X: 15, Y: 14},
	{X: 15, Y: 15},
	{X: 15, Y: 16},
	{X: 15, Y: 2},
	{X: 15, Y: 3},
	{X: 15, Y: 4},
	{X: 15, Y: 5},
	{X: 15, Y: 6},
	{X: 15, Y: 7},
	{X: 15, Y: 9},
	{X: 16, Y: 11},
	{X: 16, Y: 12},
	{X: 16, Y: 13},
	{X: 16, Y: 14},
	{X: 16, Y: 15},
	{X: 16, Y: 16},
	{X: 16, Y: 2},
	{X: 16, Y: 3},
	{X: 16, Y: 4},
	{X: 16, Y: 5},
	{X: 16, Y: 6},
	{X: 16, Y: 7},
	{X: 16, Y: 9},
	// double hazards near passages:
	{X: 2, Y: 7},
	{X: 2, Y: 9},
	{X: 2, Y: 11},
	{X: 3, Y: 7},
	{X: 3, Y: 9},
	{X: 3, Y: 11},
	{X: 7, Y: 2},
	{X: 7, Y: 3},
	{X: 7, Y: 15},
	{X: 7, Y: 16},
	{X: 9, Y: 2},
	{X: 9, Y: 3},
	{X: 9, Y: 15},
	{X: 9, Y: 16},
	{X: 11, Y: 2},
	{X: 11, Y: 3},
	{X: 11, Y: 15},
	{X: 11, Y: 16},
	{X: 15, Y: 7},
	{X: 15, Y: 9},
	{X: 15, Y: 11},
	{X: 16, Y: 7},
	{X: 16, Y: 9},
	{X: 16, Y: 11},
}

type CastleWallExtraLargeHazardsMap struct{}

func (m CastleWallExtraLargeHazardsMap) ID() string {
	return "hz_castle_wall_xl"
}

func (m CastleWallExtraLargeHazardsMap) Meta() Metadata {
	return Metadata{
		Name:        "hz_castle_wall_xl",
		Description: "Wall of hazards around the board with dangerous bridges",
		Author:      "bcambl",
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  12,
		BoardSizes:  FixedSizes(Dimensions{25, 25}),
		Tags:        []string{TAG_FOOD_PLACEMENT, TAG_HAZARD_PLACEMENT, TAG_SNAKE_PLACEMENT},
	}
}

func (m CastleWallExtraLargeHazardsMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	if !m.Meta().BoardSizes.IsAllowable(initialBoardState.Width, initialBoardState.Height) {
		return rules.RulesetError("This map can only be played on a 25x25 board")
	}

	var startPositions []rules.Point
	startPositions = append(startPositions, castleWallExtraLargeStartPositions[0]...)
	if len(initialBoardState.Snakes) >= 5 {
		startPositions = append(startPositions, castleWallExtraLargeStartPositions[1]...)
	}
	// add positions 9-12 when required
	if len(initialBoardState.Snakes) > 8 {
		startPositions = append(startPositions, castleWallExtraLargeStartPositions[2]...)
	}

	return setupCastleWallBoard(m.Meta().MaxPlayers, startPositions, castleWallExtraLargeHazards, initialBoardState, settings, editor)
}

func (m CastleWallExtraLargeHazardsMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	maxFood := 4
	return updateCastleWallBoard(maxFood, castleWallExtraLargeFood, lastBoardState, settings, editor)
}

var castleWallExtraLargeStartPositions = [][]rules.Point{
	{
		{X: 1, Y: 5},
		{X: 1, Y: 19},
		{X: 23, Y: 5},
		{X: 23, Y: 19},
	},
	{
		{X: 5, Y: 1},
		{X: 5, Y: 23},
		{X: 19, Y: 1},
		{X: 19, Y: 23},
	},
	{
		{X: 1, Y: 12},
		{X: 12, Y: 23},
		{X: 12, Y: 1},
		{X: 23, Y: 12},
	},
}

var castleWallExtraLargeFood = []rules.Point{
	{X: 3, Y: 8},
	{X: 3, Y: 12},
	{X: 3, Y: 16},
	{X: 4, Y: 8},
	{X: 4, Y: 12},
	{X: 4, Y: 16},
	{X: 8, Y: 3},
	{X: 8, Y: 4},
	{X: 8, Y: 20},
	{X: 8, Y: 21},
	{X: 12, Y: 3},
	{X: 12, Y: 4},
	{X: 12, Y: 20},
	{X: 12, Y: 21},
	{X: 16, Y: 3},
	{X: 16, Y: 4},
	{X: 16, Y: 20},
	{X: 16, Y: 21},
	{X: 20, Y: 8},
	{X: 20, Y: 12},
	{X: 20, Y: 16},
	{X: 21, Y: 8},
	{X: 21, Y: 12},
	{X: 21, Y: 16},
}

var castleWallExtraLargeHazards = []rules.Point{
	{X: 3, Y: 10},
	{X: 3, Y: 11},
	{X: 3, Y: 13},
	{X: 3, Y: 14},
	{X: 3, Y: 15},
	{X: 3, Y: 17},
	{X: 3, Y: 18},
	{X: 3, Y: 19},
	{X: 3, Y: 20},
	{X: 3, Y: 21},
	{X: 3, Y: 3},
	{X: 3, Y: 4},
	{X: 3, Y: 5},
	{X: 3, Y: 6},
	{X: 3, Y: 7},
	{X: 3, Y: 9},
	{X: 4, Y: 10},
	{X: 4, Y: 11},
	{X: 4, Y: 13},
	{X: 4, Y: 14},
	{X: 4, Y: 15},
	{X: 4, Y: 17},
	{X: 4, Y: 18},
	{X: 4, Y: 19},
	{X: 4, Y: 20},
	{X: 4, Y: 21},
	{X: 4, Y: 3},
	{X: 4, Y: 4},
	{X: 4, Y: 5},
	{X: 4, Y: 6},
	{X: 4, Y: 7},
	{X: 4, Y: 9},
	{X: 5, Y: 20},
	{X: 5, Y: 21},
	{X: 5, Y: 3},
	{X: 5, Y: 4},
	{X: 6, Y: 20},
	{X: 6, Y: 21},
	{X: 6, Y: 3},
	{X: 6, Y: 4},
	{X: 7, Y: 20},
	{X: 7, Y: 21},
	{X: 7, Y: 3},
	{X: 7, Y: 4},
	{X: 9, Y: 20},
	{X: 9, Y: 21},
	{X: 9, Y: 3},
	{X: 9, Y: 4},
	{X: 10, Y: 20},
	{X: 10, Y: 21},
	{X: 10, Y: 3},
	{X: 10, Y: 4},
	{X: 11, Y: 20},
	{X: 11, Y: 21},
	{X: 11, Y: 3},
	{X: 11, Y: 4},
	{X: 13, Y: 20},
	{X: 13, Y: 21},
	{X: 13, Y: 3},
	{X: 13, Y: 4},
	{X: 14, Y: 20},
	{X: 14, Y: 21},
	{X: 14, Y: 3},
	{X: 14, Y: 4},
	{X: 15, Y: 20},
	{X: 15, Y: 21},
	{X: 15, Y: 3},
	{X: 15, Y: 4},
	{X: 17, Y: 20},
	{X: 17, Y: 21},
	{X: 17, Y: 3},
	{X: 17, Y: 4},
	{X: 18, Y: 20},
	{X: 18, Y: 21},
	{X: 18, Y: 3},
	{X: 18, Y: 4},
	{X: 19, Y: 20},
	{X: 19, Y: 21},
	{X: 19, Y: 3},
	{X: 19, Y: 4},
	{X: 20, Y: 10},
	{X: 20, Y: 11},
	{X: 20, Y: 13},
	{X: 20, Y: 14},
	{X: 20, Y: 15},
	{X: 20, Y: 17},
	{X: 20, Y: 18},
	{X: 20, Y: 19},
	{X: 20, Y: 20},
	{X: 20, Y: 21},
	{X: 20, Y: 3},
	{X: 20, Y: 4},
	{X: 20, Y: 5},
	{X: 20, Y: 6},
	{X: 20, Y: 7},
	{X: 20, Y: 9},
	{X: 21, Y: 10},
	{X: 21, Y: 11},
	{X: 21, Y: 13},
	{X: 21, Y: 14},
	{X: 21, Y: 15},
	{X: 21, Y: 17},
	{X: 21, Y: 18},
	{X: 21, Y: 19},
	{X: 21, Y: 20},
	{X: 21, Y: 21},
	{X: 21, Y: 3},
	{X: 21, Y: 4},
	{X: 21, Y: 5},
	{X: 21, Y: 6},
	{X: 21, Y: 7},
	{X: 21, Y: 9},
	// double hazards near passages:
	{X: 3, Y: 7},
	{X: 3, Y: 9},
	{X: 3, Y: 11},
	{X: 3, Y: 13},
	{X: 3, Y: 15},
	{X: 3, Y: 17},
	{X: 4, Y: 7},
	{X: 4, Y: 9},
	{X: 4, Y: 11},
	{X: 4, Y: 13},
	{X: 4, Y: 15},
	{X: 4, Y: 17},
	{X: 7, Y: 3},
	{X: 7, Y: 4},
	{X: 7, Y: 20},
	{X: 7, Y: 21},
	{X: 9, Y: 3},
	{X: 9, Y: 4},
	{X: 9, Y: 20},
	{X: 9, Y: 21},
	{X: 11, Y: 3},
	{X: 11, Y: 4},
	{X: 11, Y: 20},
	{X: 11, Y: 21},
	{X: 13, Y: 3},
	{X: 13, Y: 4},
	{X: 13, Y: 20},
	{X: 13, Y: 21},
	{X: 15, Y: 3},
	{X: 15, Y: 4},
	{X: 15, Y: 20},
	{X: 15, Y: 21},
	{X: 17, Y: 3},
	{X: 17, Y: 4},
	{X: 17, Y: 20},
	{X: 17, Y: 21},
	{X: 20, Y: 7},
	{X: 20, Y: 9},
	{X: 20, Y: 11},
	{X: 20, Y: 13},
	{X: 20, Y: 15},
	{X: 20, Y: 17},
	{X: 21, Y: 7},
	{X: 21, Y: 9},
	{X: 21, Y: 11},
	{X: 21, Y: 13},
	{X: 21, Y: 15},
	{X: 21, Y: 17},
}
