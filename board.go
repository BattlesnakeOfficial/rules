package rules

import "math/rand"

type BoardState struct {
	Turn    int32
	Height  int32
	Width   int32
	Food    []Point
	Snakes  []Snake
	Hazards []Point
}

// NewBoardState returns an empty but fully initialized BoardState
func NewBoardState(width, height int32) *BoardState {
	return &BoardState{
		Turn:    0,
		Height:  height,
		Width:   width,
		Food:    []Point{},
		Snakes:  []Snake{},
		Hazards: []Point{},
	}
}

// Clone returns a deep copy of prevState that can be safely modified inside Ruleset.CreateNextBoardState
func (prevState *BoardState) Clone() *BoardState {
	nextState := &BoardState{
		Turn:    prevState.Turn,
		Height:  prevState.Height,
		Width:   prevState.Width,
		Food:    append([]Point{}, prevState.Food...),
		Snakes:  make([]Snake, len(prevState.Snakes)),
		Hazards: append([]Point{}, prevState.Hazards...),
	}
	for i := 0; i < len(prevState.Snakes); i++ {
		nextState.Snakes[i].ID = prevState.Snakes[i].ID
		nextState.Snakes[i].Health = prevState.Snakes[i].Health
		nextState.Snakes[i].Body = append([]Point{}, prevState.Snakes[i].Body...)
		nextState.Snakes[i].EliminatedCause = prevState.Snakes[i].EliminatedCause
		nextState.Snakes[i].EliminatedBy = prevState.Snakes[i].EliminatedBy
	}
	return nextState
}

// CreateDefaultBoardState is a convenience function for fully initializing a
// "default" board state with snakes and food.
// In a real game, the engine may generate the board without calling this
// function, or customize the results based on game-specific settings.
func CreateDefaultBoardState(width int32, height int32, snakeIDs []string) (*BoardState, error) {
	initialBoardState := NewBoardState(width, height)

	err := PlaceSnakesAutomatically(initialBoardState, snakeIDs)
	if err != nil {
		return nil, err
	}

	err = PlaceFoodAutomatically(initialBoardState)
	if err != nil {
		return nil, err
	}

	return initialBoardState, nil
}

// PlaceSnakesAutomatically initializes the array of snakes based on the provided snake IDs and the size of the board.
func PlaceSnakesAutomatically(b *BoardState, snakeIDs []string) error {
	if isKnownBoardSize(b) {
		return PlaceSnakesFixed(b, snakeIDs)
	}
	return PlaceSnakesRandomly(b, snakeIDs)
}

func PlaceSnakesFixed(b *BoardState, snakeIDs []string) error {
	b.Snakes = make([]Snake, len(snakeIDs))

	for i := 0; i < len(snakeIDs); i++ {
		b.Snakes[i] = Snake{
			ID:     snakeIDs[i],
			Health: SnakeMaxHealth,
		}
	}

	// Create start 8 points
	mn, md, mx := int32(1), (b.Width-1)/2, b.Width-2
	startPoints := []Point{
		{mn, mn},
		{mn, md},
		{mn, mx},
		{md, mn},
		{md, mx},
		{mx, mn},
		{mx, md},
		{mx, mx},
	}

	// Sanity check
	if len(b.Snakes) > len(startPoints) {
		return ErrorTooManySnakes
	}

	// Randomly order them
	rand.Shuffle(len(startPoints), func(i int, j int) {
		startPoints[i], startPoints[j] = startPoints[j], startPoints[i]
	})

	// Assign to snakes in order given
	for i := 0; i < len(b.Snakes); i++ {
		for j := 0; j < SnakeStartSize; j++ {
			b.Snakes[i].Body = append(b.Snakes[i].Body, startPoints[i])
		}

	}
	return nil
}

func PlaceSnakesRandomly(b *BoardState, snakeIDs []string) error {
	b.Snakes = make([]Snake, len(snakeIDs))

	for i := 0; i < len(snakeIDs); i++ {
		b.Snakes[i] = Snake{
			ID:     snakeIDs[i],
			Health: SnakeMaxHealth,
		}
	}

	for i := 0; i < len(b.Snakes); i++ {
		unoccupiedPoints := getEvenUnoccupiedPoints(b)
		if len(unoccupiedPoints) <= 0 {
			return ErrorNoRoomForSnake
		}
		p := unoccupiedPoints[rand.Intn(len(unoccupiedPoints))]
		for j := 0; j < SnakeStartSize; j++ {
			b.Snakes[i].Body = append(b.Snakes[i].Body, p)
		}
	}
	return nil
}

// PlaceSnake adds a snake to the board with the given ID and body coordinates.
func PlaceSnake(b *BoardState, snakeID string, body []Point) error {
	b.Snakes = append(b.Snakes, Snake{
		ID:     snakeID,
		Health: SnakeMaxHealth,
		Body:   body,
	})
	return nil
}

// PlaceFoodAutomatically initializes the array of food based on the size of the board and the number of snakes.
func PlaceFoodAutomatically(b *BoardState) error {
	if isKnownBoardSize(b) {
		return PlaceFoodFixed(b)
	}
	return PlaceFoodRandomly(b, int32(len(b.Snakes)))
}

func PlaceFoodFixed(b *BoardState) error {
	// Place 1 food within exactly 2 moves of each snake
	for i := 0; i < len(b.Snakes); i++ {
		snakeHead := b.Snakes[i].Body[0]
		possibleFoodLocations := []Point{
			{snakeHead.X - 1, snakeHead.Y - 1},
			{snakeHead.X - 1, snakeHead.Y + 1},
			{snakeHead.X + 1, snakeHead.Y - 1},
			{snakeHead.X + 1, snakeHead.Y + 1},
		}
		availableFoodLocations := []Point{}

		for _, p := range possibleFoodLocations {
			isOccupiedAlready := false
			for _, food := range b.Food {
				if food.X == p.X && food.Y == p.Y {
					isOccupiedAlready = true
					break
				}
			}

			if !isOccupiedAlready {
				availableFoodLocations = append(availableFoodLocations, p)
			}
		}

		if len(availableFoodLocations) <= 0 {
			return ErrorNoRoomForFood
		}

		// Select randomly from available locations
		placedFood := availableFoodLocations[rand.Intn(len(availableFoodLocations))]
		b.Food = append(b.Food, placedFood)
	}

	// Finally, always place 1 food in center of board for dramatic purposes
	isCenterOccupied := true
	centerCoord := Point{(b.Width - 1) / 2, (b.Height - 1) / 2}
	unoccupiedPoints := getUnoccupiedPoints(b, true)
	for _, point := range unoccupiedPoints {
		if point == centerCoord {
			isCenterOccupied = false
			break
		}
	}
	if isCenterOccupied {
		return ErrorNoRoomForFood
	}
	b.Food = append(b.Food, centerCoord)

	return nil
}

// PlaceFoodRandomly adds up to n new food to the board in random unoccupied squares
func PlaceFoodRandomly(b *BoardState, n int32) error {
	for i := int32(0); i < n; i++ {
		unoccupiedPoints := getUnoccupiedPoints(b, false)
		if len(unoccupiedPoints) > 0 {
			newFood := unoccupiedPoints[rand.Intn(len(unoccupiedPoints))]
			b.Food = append(b.Food, newFood)
		}
	}
	return nil
}

func isKnownBoardSize(b *BoardState) bool {
	if b.Height == BoardSizeSmall && b.Width == BoardSizeSmall {
		return true
	}
	if b.Height == BoardSizeMedium && b.Width == BoardSizeMedium {
		return true
	}
	if b.Height == BoardSizeLarge && b.Width == BoardSizeLarge {
		return true
	}
	return false
}

func getUnoccupiedPoints(b *BoardState, includePossibleMoves bool) []Point {
	pointIsOccupied := map[int32]map[int32]bool{}
	for _, p := range b.Food {
		if _, xExists := pointIsOccupied[p.X]; !xExists {
			pointIsOccupied[p.X] = map[int32]bool{}
		}
		pointIsOccupied[p.X][p.Y] = true
	}
	for _, snake := range b.Snakes {
		if snake.EliminatedCause != NotEliminated {
			continue
		}
		for i, p := range snake.Body {
			if _, xExists := pointIsOccupied[p.X]; !xExists {
				pointIsOccupied[p.X] = map[int32]bool{}
			}
			pointIsOccupied[p.X][p.Y] = true

			if i == 0 && !includePossibleMoves {
				nextMovePoints := []Point{
					{X: p.X - 1, Y: p.Y},
					{X: p.X + 1, Y: p.Y},
					{X: p.X, Y: p.Y - 1},
					{X: p.X, Y: p.Y + 1},
				}
				for _, nextP := range nextMovePoints {
					if _, xExists := pointIsOccupied[nextP.X]; !xExists {
						pointIsOccupied[nextP.X] = map[int32]bool{}
					}
					pointIsOccupied[nextP.X][nextP.Y] = true
				}
			}
		}
	}

	unoccupiedPoints := []Point{}
	for x := int32(0); x < b.Width; x++ {
		for y := int32(0); y < b.Height; y++ {
			if _, xExists := pointIsOccupied[x]; xExists {
				if isOccupied, yExists := pointIsOccupied[x][y]; yExists {
					if isOccupied {
						continue
					}
				}
			}
			unoccupiedPoints = append(unoccupiedPoints, Point{X: x, Y: y})
		}
	}
	return unoccupiedPoints
}

func getEvenUnoccupiedPoints(b *BoardState) []Point {
	// Start by getting unoccupied points
	unoccupiedPoints := getUnoccupiedPoints(b, true)

	// Create a new array to hold points that are  even
	evenUnoccupiedPoints := []Point{}

	for _, point := range unoccupiedPoints {
		if ((point.X + point.Y) % 2) == 0 {
			evenUnoccupiedPoints = append(evenUnoccupiedPoints, point)
		}
	}
	return evenUnoccupiedPoints
}
