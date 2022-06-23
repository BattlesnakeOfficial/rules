package maps

import (
	"math"

	"github.com/BattlesnakeOfficial/rules"
)

type InnerBorderHazardsMap struct{}

func init() {
	globalRegistry.RegisterMap("hz_inner_wall", InnerBorderHazardsMap{})
	globalRegistry.RegisterMap("hz_rings", ConcentricRingsHazardsMap{})
	globalRegistry.RegisterMap("hz_columns", ColumnsHazardsMap{})
	globalRegistry.RegisterMap("hz_rivers_bridges", RiverAndBridgesHazardsMap{})
	globalRegistry.RegisterMap("hz_spiral", SpiralHazardsMap{})
	globalRegistry.RegisterMap("hz_scatter", ScatterFillMap{})
	globalRegistry.RegisterMap("hz_grow_box", DirectionalExpandingBoxMap{})
	globalRegistry.RegisterMap("hz_expand_box", ExpandingBoxMap{})
	globalRegistry.RegisterMap("hz_expand_scatter", ExpandingScatterMap{})
}

func (m InnerBorderHazardsMap) ID() string {
	return "hz_inner_wall"
}

func (m InnerBorderHazardsMap) Meta() Metadata {
	return Metadata{
		Name:        "hz_inner_wall",
		Description: "Creates a static map on turn 0 that is a 1-square wall of hazard that is inset 2 squares from the edge of the board",
		Author:      "Battlesnake",
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  8,
		BoardSizes:  AnySize(),
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
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  8,
		BoardSizes:  AnySize(),
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
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  8,
		BoardSizes:  AnySize(),
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

type SpiralHazardsMap struct{}

func (m SpiralHazardsMap) ID() string {
	return "hz_spiral"
}

func (m SpiralHazardsMap) Meta() Metadata {
	return Metadata{
		Name: "hz_spiral",
		Description: `Generates a dynamic hazard map that grows in a spiral pattern clockwise from a random point on
 the map. Each 2 turns a new hazard square is added to the map`,
		Author:     "altersaddle",
		Version:    1,
		MinPlayers: 1,
		MaxPlayers: 8,
		BoardSizes: AnySize(),
	}
}

func (m SpiralHazardsMap) SetupBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return (StandardMap{}).SetupBoard(lastBoardState, settings, editor)
}

func (m SpiralHazardsMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	err := StandardMap{}.UpdateBoard(lastBoardState, settings, editor)
	if err != nil {
		return err
	}

	currentTurn := lastBoardState.Turn + 1
	spawnEveryNTurns := 3

	// no-op if we're not on a turn that spawns hazards
	if currentTurn < spawnEveryNTurns || currentTurn%spawnEveryNTurns != 0 {
		return nil
	}

	rand := settings.GetRand(0)
	spawnArea := 0.3 // Center spiral in the middle 0.6 of the board

	// randomly choose a location between the start point and the edge of the board
	spawnOffsetX := int(math.Floor(float64(lastBoardState.Width) * spawnArea))
	maxX := lastBoardState.Width - 1 - spawnOffsetX
	startX := rand.Range(spawnOffsetX, maxX)
	spawnOffsetY := int(math.Floor(float64(lastBoardState.Height) * spawnArea))
	maxY := lastBoardState.Height - 1 - spawnOffsetY
	startY := rand.Range(spawnOffsetY, maxY)

	if currentTurn == spawnEveryNTurns {
		editor.AddHazard(rules.Point{X: startX, Y: startY})
		return nil
	}

	// determine number of rings in spiral
	numRings := maxInt(startX, startY, lastBoardState.Width-startX, lastBoardState.Height-startY)

	turnCtr := spawnEveryNTurns
	for ring := 0; ring < numRings; ring++ {
		offset := ring + 1
		x := startX - ring
		y := startY + offset

		numSquaresInRing := 8 * offset
		for i := 0; i < numSquaresInRing; i++ {
			turnCtr += spawnEveryNTurns

			if turnCtr > currentTurn {
				break
			}

			if turnCtr == currentTurn && isOnBoard(lastBoardState.Width, lastBoardState.Height, x, y) {
				editor.AddHazard(rules.Point{X: x, Y: y})
			}

			// move the "cursor"
			if y == startY+offset && x < startX+offset {
				// top line, move right
				x += 1
			} else if x == startX+offset && y > startY-offset {

				// right side, go down
				y -= 1
			} else if y == startY-offset && x > startX-offset {
				// bottom line, move left
				x -= 1
			} else if x == startX-offset && y < startY+offset {
				y += 1
			}
		}
	}

	return nil
}

type ScatterFillMap struct{}

func (m ScatterFillMap) ID() string {
	return "hz_scatter"
}

func (m ScatterFillMap) Meta() Metadata {
	return Metadata{
		Name:        "hz_scatter",
		Description: `Fills the entire board with hazard squares that are set to appear on regular turn schedule. Each square is picked at random.`,
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  8,
		BoardSizes:  AnySize(),
	}
}

func (m ScatterFillMap) SetupBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return (StandardMap{}).SetupBoard(lastBoardState, settings, editor)
}

func (m ScatterFillMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	err := StandardMap{}.UpdateBoard(lastBoardState, settings, editor)
	if err != nil {
		return err
	}

	currentTurn := lastBoardState.Turn + 1
	spawnEveryNTurns := 2

	// no-op if we're not on a turn that spawns hazards
	if currentTurn < spawnEveryNTurns || currentTurn%spawnEveryNTurns != 0 {
		return nil
	}

	positions := make([]rules.Point, 0, lastBoardState.Width*lastBoardState.Height)
	for x := 0; x < lastBoardState.Width; x++ {
		for y := 0; y < lastBoardState.Height; y++ {
			positions = append(positions, rules.Point{X: x, Y: y})
		}
	}

	rand := settings.GetRand(0)
	rand.Shuffle(len(positions), func(i, j int) {
		positions[i], positions[j] = positions[j], positions[i]
	})

	editor.AddHazard(positions[(currentTurn-2)/2])
	return nil
}

type DirectionalExpandingBoxMap struct{}

func (m DirectionalExpandingBoxMap) ID() string {
	return "hz_grow_box"
}

func (m DirectionalExpandingBoxMap) Meta() Metadata {
	return Metadata{
		Name:        "hz_grow_box",
		Description: `Creates an area of hazard that expands from a point with one random side growing on a turn schedule.`,
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  8,
		BoardSizes:  AnySize(),
	}
}

func (m DirectionalExpandingBoxMap) SetupBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return (StandardMap{}).SetupBoard(lastBoardState, settings, editor)
}

func (m DirectionalExpandingBoxMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	err := StandardMap{}.UpdateBoard(lastBoardState, settings, editor)
	if err != nil {
		return err
	}

	currentTurn := lastBoardState.Turn + 1
	startTurn := 1
	spawnEveryNTurns := 15

	// no-op if we're not on a turn that spawns hazards
	if (currentTurn-startTurn)%spawnEveryNTurns != 0 {
		return nil
	}

	// no-op if we have spawned the entire board already
	if len(lastBoardState.Hazards) == lastBoardState.Width*lastBoardState.Height {
		return nil
	}

	rand := settings.GetRand(0)
	startX := rand.Range(2, lastBoardState.Width-2)
	startY := rand.Range(2, lastBoardState.Height-2)

	if currentTurn == 1 {
		editor.AddHazard(rules.Point{X: startX, Y: startY})
		return nil
	}

	topLeft := rules.Point{X: startX, Y: startY}
	bottomRight := rules.Point{X: startX, Y: startY}

	// var growthDirection string
	maxTurns := (currentTurn - startTurn) / spawnEveryNTurns
	for i := 0; i < maxTurns; i++ {
		directions := []string{}
		if topLeft.X > 0 {
			directions = append(directions, "left")
		}
		if topLeft.Y < lastBoardState.Height-1 {
			directions = append(directions, "up")
		}
		if bottomRight.X < lastBoardState.Width-1 {
			directions = append(directions, "right")
		}
		if bottomRight.Y > 0 {
			directions = append(directions, "down")
		}
		if len(directions) == 0 {
			return nil
		}
		choice := rand.Intn(len(directions))
		growthDirection := directions[choice]

		addHazards := i == maxTurns-1

		if growthDirection == "left" {
			x := topLeft.X - 1
			if addHazards {
				for y := bottomRight.Y; y < topLeft.Y+1; y++ {
					editor.AddHazard(rules.Point{X: x, Y: y})
				}
			}
			topLeft.X = x
		} else if growthDirection == "right" {
			x := bottomRight.X + 1
			if addHazards {
				for y := bottomRight.Y; y < topLeft.Y+1; y++ {
					editor.AddHazard(rules.Point{X: x, Y: y})
				}
			}
			bottomRight.X = x
		} else if growthDirection == "up" {
			y := topLeft.Y + 1
			if addHazards {
				for x := topLeft.X; x < bottomRight.X+1; x++ {
					editor.AddHazard(rules.Point{X: x, Y: y})
				}
			}
			topLeft.Y = y
		} else if growthDirection == "down" {
			y := bottomRight.Y - 1
			if addHazards {
				for x := topLeft.X; x < bottomRight.X+1; x++ {
					editor.AddHazard(rules.Point{X: x, Y: y})
				}
			}
			bottomRight.Y = y
		}
	}
	return nil
}

type ExpandingBoxMap struct{}

func (m ExpandingBoxMap) ID() string {
	return "hz_expand_box"
}

func (m ExpandingBoxMap) Meta() Metadata {
	return Metadata{
		Name:        "hz_expand_box",
		Description: `Generates an area of hazard that expands from a random point on the board outward in concentric rings on a periodic turn schedule.`,
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  8,
		BoardSizes:  AnySize(),
	}
}

func (m ExpandingBoxMap) SetupBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return (StandardMap{}).SetupBoard(lastBoardState, settings, editor)
}

func (m ExpandingBoxMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	err := StandardMap{}.UpdateBoard(lastBoardState, settings, editor)
	if err != nil {
		return err
	}

	currentTurn := lastBoardState.Turn + 1
	startTurn := 1 // first hazard appears on turn 1
	spawnEveryNTurns := 20

	// no-op if we're not on a turn that spawns hazards
	if (currentTurn-startTurn)%spawnEveryNTurns != 0 {
		return nil
	}

	// no-op if we have spawned the entire board already
	if len(lastBoardState.Hazards) == lastBoardState.Width*lastBoardState.Height {
		return nil
	}

	rand := settings.GetRand(0)

	startX := rand.Range(2, lastBoardState.Width-2)
	startY := rand.Range(2, lastBoardState.Width-2)

	if currentTurn == startTurn {
		editor.AddHazard(rules.Point{X: startX, Y: startY})
		return nil
	}

	// determine number of rings in spiral
	numRings := maxInt(startX, startY, lastBoardState.Width-startX, lastBoardState.Height-startY)

	// no-op when iterations exceed the max rings
	if currentTurn/spawnEveryNTurns > numRings {
		return nil
	}

	ring := currentTurn/spawnEveryNTurns - 1
	offset := ring + 1

	for x := startX - offset; x < startX+offset+1; x++ {
		for y := startY - offset; y < startY+offset+1; y++ {
			if isOnBoard(lastBoardState.Width, lastBoardState.Height, x, y) {
				if ((x == startX-offset || x == startX+offset) && y >= startY-offset && y <= startY+offset) || ((y == startY-offset || y == startY+offset) && x >= startX-offset && x <= startX+offset) {
					editor.AddHazard(rules.Point{X: x, Y: y})
				}
			}
		}
	}

	return nil
}

type ExpandingScatterMap struct{}

func (m ExpandingScatterMap) ID() string {
	return "hz_expand_scatter"
}

func (m ExpandingScatterMap) Meta() Metadata {
	return Metadata{
		Name:        "hz_expand_scatter",
		Description: `Builds an expanding hazard area that grows from a central point in rings that are randomly filled in on a regular turn schedule.`,
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  8,
		BoardSizes:  AnySize(),
	}
}

func (m ExpandingScatterMap) SetupBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return (StandardMap{}).SetupBoard(lastBoardState, settings, editor)
}

func (m ExpandingScatterMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	err := StandardMap{}.UpdateBoard(lastBoardState, settings, editor)
	if err != nil {
		return err
	}

	currentTurn := lastBoardState.Turn + 1
	startTurn := 1 // first hazard appears on turn 1
	spawnEveryNTurns := 2

	// no-op if we're not on a turn that spawns hazards
	if (currentTurn-startTurn)%spawnEveryNTurns != 0 {
		return nil
	}

	// no-op if we have spawned the entire board already
	if len(lastBoardState.Hazards) == lastBoardState.Width*lastBoardState.Height {
		return nil
	}

	rand := settings.GetRand(0)

	startX := rand.Range(1, lastBoardState.Width-1)
	startY := rand.Range(1, lastBoardState.Width-1)

	if currentTurn == startTurn {
		editor.AddHazard(rules.Point{X: startX, Y: startY})
		return nil
	}

	// determine number of rings in spiral
	numRings := maxInt(startX, startY, lastBoardState.Width-startX, lastBoardState.Height-startY)

	allPositions := []rules.Point{}
	for ring := 0; ring < numRings; ring++ {
		offset := ring + 1
		positions := []rules.Point{}
		for x := startX - offset; x < startX+offset+1; x++ {
			for y := startY - offset; y < startY+offset+1; y++ {
				if isOnBoard(lastBoardState.Width, lastBoardState.Height, x, y) {
					if ((x == startX-offset || x == startX+offset) && y >= startY-offset && y <= startY+offset) || ((y == startY-offset || y == startY+offset) && x >= startX-offset && x <= startX+offset) {
						positions = append(positions, rules.Point{X: x, Y: y})
					}
				}
			}
		}
		// shuffle the positions so they are added scattered/randomly
		rand.Shuffle(len(positions), func(i, j int) {
			positions[i], positions[j] = positions[j], positions[i]
		})
		allPositions = append(allPositions, positions...)
	}

	chosenPos := currentTurn/spawnEveryNTurns - 1
	editor.AddHazard(allPositions[chosenPos])

	return nil
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
		Author:     "Battlesnake",
		Version:    1,
		MinPlayers: 1,
		MaxPlayers: 12,
		BoardSizes: FixedSizes(Dimensions{11, 11}, Dimensions{19, 19}, Dimensions{25, 25}),
	}
}

func (m RiverAndBridgesHazardsMap) SetupBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	width := lastBoardState.Width
	height := lastBoardState.Height
	startPositions, ok := riversAndBridgesStartPositions[rules.Point{X: width, Y: height}]
	if !ok {
		return rules.RulesetError("board size is not supported by this map")
	}

	numSnakes := len(lastBoardState.Snakes)
	if numSnakes == 0 {
		return rules.RulesetError("too few snakes - at least one snake must be present")
	}

	maxSnakes := len(startPositions)
	if maxSnakes < numSnakes {
		return rules.ErrorTooManySnakes
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
	err := rules.PlaceSnakesAtPositions(tempBoardState, startPositions)
	if err != nil {
		return err
	}

	err = rules.PlaceFoodAutomatically(rand, tempBoardState)
	if err != nil {
		return err
	}

	// Copy food from temp board state
	for _, food := range tempBoardState.Food {
		editor.AddFood(food)
	}

	// Copy snakes from temp board state
	for _, snake := range tempBoardState.Snakes {
		editor.PlaceSnake(snake.ID, snake.Body, snake.Health)
	}

	hazards, ok := riversAndBridgesMaps[rules.Point{X: width, Y: height}]
	if !ok {
		return rules.RulesetError("board size is not supported by this map")
	}

	for _, p := range hazards {
		editor.AddHazard(p)
	}

	return nil
}

func (m RiverAndBridgesHazardsMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return StandardMap{}.UpdateBoard(lastBoardState, settings, editor)
}

var riversAndBridgesStartPositions = map[rules.Point][]rules.Point{
	{X: 11, Y: 11}: {
		{X: 1, Y: 1},
		{X: 9, Y: 9},
		{X: 1, Y: 9},
		{X: 9, Y: 1},
		{X: 3, Y: 3},
		{X: 7, Y: 7},
		{X: 3, Y: 7},
		{X: 7, Y: 3},
		{X: 1, Y: 3},
		{X: 9, Y: 7},
		{X: 9, Y: 3},
		{X: 3, Y: 9},
	},
	{X: 19, Y: 19}: {
		{X: 1, Y: 1},
		{X: 17, Y: 17},
		{X: 13, Y: 1},
		{X: 1, Y: 17},
		{X: 5, Y: 1},
		{X: 17, Y: 1},
		{X: 5, Y: 17},
		{X: 17, Y: 13},
		{X: 1, Y: 5},
		{X: 17, Y: 5},
		{X: 1, Y: 13},
		{X: 13, Y: 17},
		{X: 5, Y: 5},
		{X: 13, Y: 5},
		{X: 5, Y: 13},
		{X: 13, Y: 13},
	},
	{X: 25, Y: 25}: {
		{X: 1, Y: 1},
		{X: 23, Y: 23},
		{X: 1, Y: 23},
		{X: 23, Y: 1},
		{X: 9, Y: 9},
		{X: 15, Y: 15},
		{X: 15, Y: 9},
		{X: 9, Y: 15},
		{X: 9, Y: 1},
		{X: 23, Y: 15},
		{X: 23, Y: 9},
		{X: 1, Y: 15},
		{X: 1, Y: 9},
		{X: 15, Y: 1},
		{X: 9, Y: 23},
		{X: 15, Y: 23},
	},
}

var riversAndBridgesMaps = map[rules.Point][]rules.Point{
	{X: 11, Y: 11}: {
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
	{X: 19, Y: 19}: {
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
	{X: 25, Y: 25}: {
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
