package maps

import (
	"github.com/BattlesnakeOfficial/rules"
)

type CastleWallMap struct{}

func init() {
	globalRegistry.RegisterMap("castle_wall", CastleWallMap{})
}

func (m CastleWallMap) ID() string {
	return "castle_wall"
}

func (m CastleWallMap) Meta() Metadata {
	return Metadata{
		Name:        "Castle Wall",
		Description: "Wall of hazards around the board with dangerous bridges",
		Author:      "bcambl",
		Version:     1,
		MinPlayers:  2,
		MaxPlayers:  8,
		BoardSizes: FixedSizes(
			Dimensions{11, 11},
			Dimensions{19, 19},
			Dimensions{25, 25}),
	}
}

func (m CastleWallMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {

	rand := settings.GetRand(initialBoardState.Turn)

	// Place snakes
	snakePositions, ok := CastleWallPositions.Snakes[rules.Point{X: initialBoardState.Width, Y: initialBoardState.Height}]
	if !ok {
		return rules.RulesetError("board size is not supported by this map")
	}

	var snakes []rules.Point
	// always support up to 8x snakes on all supported board sizes
	snakes = append(snakes, snakePositions[0]...)
	if len(initialBoardState.Snakes) >= 5 {
		snakes = append(snakes, snakePositions[1]...)
	}

	// only support 8 or less snakes on boards smaller than XLarge
	if (len(initialBoardState.Snakes) > 8) && (initialBoardState.Width < rules.BoardSizeXLarge) {
		return rules.ErrorTooManySnakes
	}

	// only support 12 or less snakes for XLarge and XXLarge board sizes
	if (initialBoardState.Width >= rules.BoardSizeXLarge) && (len(initialBoardState.Snakes) > 12) {
		return rules.ErrorTooManySnakes
	}

	// add positions 9-12 for XLarge and XXLarge board sizes
	if (initialBoardState.Width >= rules.BoardSizeXLarge) && (len(initialBoardState.Snakes) > 8) {
		snakes = append(snakes, snakePositions[2]...)
	}

	rand.Shuffle(len(snakes), func(i int, j int) {
		snakes[i], snakes[j] = snakes[j], snakes[i]
	})

	for index, snake := range initialBoardState.Snakes {
		head := snakes[index]
		editor.PlaceSnake(snake.ID, []rules.Point{head, head, head}, snake.Health)
	}

	hazards, ok := CastleWallPositions.Hazards[rules.Point{X: initialBoardState.Width, Y: initialBoardState.Height}]
	if !ok {
		return rules.RulesetError("board size is not supported by this map")
	}

	for _, h := range hazards {
		editor.AddHazard(h)
	}
	return nil
}

func (m CastleWallMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {

	// no food spawning for first 10 turns to reduce favoring snakes spawning closer to bridges
	if lastBoardState.Turn < 10 {
		return nil
	}

	// max of 2 food on small and medium boards
	if len(lastBoardState.Food) > 1 && lastBoardState.Width < rules.BoardSizeXLarge {
		return nil
	}

	// max of 4 food on XLarge and XXLarge boards
	if len(lastBoardState.Food) > 3 && lastBoardState.Width >= rules.BoardSizeXLarge {
		return nil
	}

	rand := settings.GetRand(lastBoardState.Turn)

	food, ok := CastleWallPositions.Food[rules.Point{X: lastBoardState.Width, Y: lastBoardState.Height}]
	if !ok {
		return rules.RulesetError("board size is not supported by this map")
	}

	rand.Shuffle(len(food), func(i int, j int) {
		food[i], food[j] = food[j], food[i]
	})

	for _, f := range food {
		tileIsOccupied := false

	snakeLoop:
		for _, snake := range lastBoardState.Snakes {
			for _, point := range snake.Body {
				if point.X == f.X && point.Y == f.Y {
					tileIsOccupied = true
					break snakeLoop
				}
			}
		}

		for _, existingFood := range lastBoardState.Food {
			if existingFood.X == f.X && existingFood.Y == f.Y {
				tileIsOccupied = true
				break
			}

			// also avoid spawning food in same passage as existing food
			if existingFood.X+1 == f.X && existingFood.Y == f.Y {
				tileIsOccupied = true
				break
			}
			if existingFood.X-1 == f.X && existingFood.Y == f.Y {
				tileIsOccupied = true
				break
			}
			if existingFood.X == f.X && existingFood.Y+1 == f.Y {
				tileIsOccupied = true
				break
			}
			if existingFood.X == f.X && existingFood.Y-1 == f.Y {
				tileIsOccupied = true
				break
			}
		}

		if !tileIsOccupied {
			editor.AddFood(f)
			break
		}
	}

	return nil
}

type CastleWall struct {
	Snakes  map[rules.Point][][]rules.Point
	Food    map[rules.Point][]rules.Point
	Hazards map[rules.Point][]rules.Point
}

// CastleWallPositions contains all starting snake positions, food spawns, and hazard locations
var CastleWallPositions = CastleWall{
	Snakes: map[rules.Point][][]rules.Point{
		{X: 11, Y: 11}: {
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
		},
		{X: 19, Y: 19}: {
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
		},
		{X: 25, Y: 25}: {
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
		},
	},
	Food: map[rules.Point][]rules.Point{
		{X: 11, Y: 11}: {
			{X: 2, Y: 5},
			{X: 5, Y: 2},
			{X: 5, Y: 8},
			{X: 8, Y: 5},
		},
		{X: 19, Y: 19}: {
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
		},
		{X: 25, Y: 25}: {
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
		},
	},
	Hazards: map[rules.Point][]rules.Point{
		{X: 11, Y: 11}: {
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
		},
		{X: 19, Y: 19}: {
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
		},
		{X: 25, Y: 25}: {
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
		},
	},
}
