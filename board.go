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

// Clone returns a deep copy of prevState that can be safely modified without affecting the original
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
		nextState.Snakes[i].Body = append([]Point{}, prevState.Snakes[i].Body...)
		nextState.Snakes[i].Health = prevState.Snakes[i].Health
		nextState.Snakes[i].EliminatedCause = prevState.Snakes[i].EliminatedCause
		nextState.Snakes[i].EliminatedOnTurn = prevState.Snakes[i].EliminatedOnTurn
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

func InitializeSnakes(b *BoardState, snakeIDs []string) {
	b.Snakes = make([]Snake, len(snakeIDs))

	for i := 0; i < len(snakeIDs); i++ {
		b.Snakes[i] = Snake{
			ID:     snakeIDs[i],
			Health: SnakeMaxHealth,
			Body:   []Point{},
		}
	}
}

// PlaceSnakesAutomatically initializes the array of snakes based on the provided snake IDs and the size of the board.
func PlaceSnakesAutomatically(b *BoardState, snakeIDs []string) error {
	if isKnownBoardSize(b) {
		return PlaceSnakesFixed(b, snakeIDs)
	}
	return PlaceSnakesRandomly(b, snakeIDs)
}

func PlaceSnakesFixed(b *BoardState, snakeIDs []string) error {
	InitializeSnakes(b, snakeIDs)

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
	InitializeSnakes(b, snakeIDs)

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
	for index, snake := range b.Snakes {
		if snake.ID == snakeID {
			b.Snakes[index].Body = body
			return nil
		}
	}

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
	centerCoord := Point{(b.Width - 1) / 2, (b.Height - 1) / 2}

	// Place 1 food within exactly 2 moves of each snake, but never towards the center or in a corner
	for i := 0; i < len(b.Snakes); i++ {
		snakeHead := b.Snakes[i].Body[0]
		possibleFoodLocations := []Point{
			{snakeHead.X - 1, snakeHead.Y - 1},
			{snakeHead.X - 1, snakeHead.Y + 1},
			{snakeHead.X + 1, snakeHead.Y - 1},
			{snakeHead.X + 1, snakeHead.Y + 1},
		}

		// Remove any invalid/unwanted positions
		availableFoodLocations := []Point{}
		for _, p := range possibleFoodLocations {

			// Ignore points already occupied by food
			isOccupiedAlready := false
			for _, food := range b.Food {
				if food.X == p.X && food.Y == p.Y {
					isOccupiedAlready = true
					break
				}
			}
			if isOccupiedAlready {
				continue
			}

			// Food must be further than snake from center on at least one axis
			isAwayFromCenter := false
			if p.X < snakeHead.X && snakeHead.X < centerCoord.X {
				isAwayFromCenter = true
			} else if centerCoord.X < snakeHead.X && snakeHead.X < p.X {
				isAwayFromCenter = true
			} else if p.Y < snakeHead.Y && snakeHead.Y < centerCoord.Y {
				isAwayFromCenter = true
			} else if centerCoord.Y < snakeHead.Y && snakeHead.Y < p.Y {
				isAwayFromCenter = true
			}
			if !isAwayFromCenter {
				continue
			}

			// Don't spawn food in corners
			if (p.X == 0 || p.X == (b.Width-1)) && (p.Y == 0 || p.Y == (b.Height-1)) {
				continue
			}

			availableFoodLocations = append(availableFoodLocations, p)
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

func absInt32(n int32) int32 {
	if n < 0 {
		return -n
	}
	return n
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

func getDistanceBetweenPoints(a, b Point) int32 {
	return absInt32(a.X-b.X) + absInt32(a.Y-b.Y)
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

// Automatically place the snakes already present in the BoardState
func placeSnakesAutomaticallyWithRand(rand Rand, b *BoardState) error {
	if isKnownBoardSize(b) {
		return placeSnakesFixedWithRand(rand, b)
	}
	return placeSnakesRandomlyWithRand(rand, b)
}

// Place the snakes already present in the BoardState based on fixed locations on known board sizes.
func placeSnakesFixedWithRand(rand Rand, b *BoardState) error {
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
		b.Snakes[i].Body = []Point{}

		for j := 0; j < SnakeStartSize; j++ {
			b.Snakes[i].Body = append(b.Snakes[i].Body, startPoints[i])
		}

	}
	return nil
}

// Place the snakes already present in the BoardState randomly.
func placeSnakesRandomlyWithRand(rand Rand, b *BoardState) error {
	for i := 0; i < len(b.Snakes); i++ {
		unoccupiedPoints := getEvenUnoccupiedPoints(b)
		if len(unoccupiedPoints) <= 0 {
			return ErrorNoRoomForSnake
		}
		b.Snakes[i].Body = []Point{}
		p := unoccupiedPoints[rand.Intn(len(unoccupiedPoints))]
		for j := 0; j < SnakeStartSize; j++ {
			b.Snakes[i].Body = append(b.Snakes[i].Body, p)
		}
	}
	return nil
}

// PlaceFoodAutomatically initializes the array of food based on the size of the board and the number of snakes.
func placeFoodAutomaticallyWithRand(rand Rand, b *BoardState) error {
	if isKnownBoardSize(b) {
		return placeFoodFixedWithRand(rand, b)
	}
	return placeFoodRandomlyWithRand(rand, b, int32(len(b.Snakes)))
}

func placeFoodFixedWithRand(rand Rand, b *BoardState) error {
	centerCoord := Point{(b.Width - 1) / 2, (b.Height - 1) / 2}

	// Place 1 food within exactly 2 moves of each snake, but never towards the center or in a corner
	for i := 0; i < len(b.Snakes); i++ {
		snakeHead := b.Snakes[i].Body[0]
		possibleFoodLocations := []Point{
			{snakeHead.X - 1, snakeHead.Y - 1},
			{snakeHead.X - 1, snakeHead.Y + 1},
			{snakeHead.X + 1, snakeHead.Y - 1},
			{snakeHead.X + 1, snakeHead.Y + 1},
		}

		// Remove any invalid/unwanted positions
		availableFoodLocations := []Point{}
		for _, p := range possibleFoodLocations {

			// Ignore points already occupied by food
			isOccupiedAlready := false
			for _, food := range b.Food {
				if food.X == p.X && food.Y == p.Y {
					isOccupiedAlready = true
					break
				}
			}
			if isOccupiedAlready {
				continue
			}

			// Food must be further than snake from center on at least one axis
			isAwayFromCenter := false
			if p.X < snakeHead.X && snakeHead.X < centerCoord.X {
				isAwayFromCenter = true
			} else if centerCoord.X < snakeHead.X && snakeHead.X < p.X {
				isAwayFromCenter = true
			} else if p.Y < snakeHead.Y && snakeHead.Y < centerCoord.Y {
				isAwayFromCenter = true
			} else if centerCoord.Y < snakeHead.Y && snakeHead.Y < p.Y {
				isAwayFromCenter = true
			}
			if !isAwayFromCenter {
				continue
			}

			// Don't spawn food in corners
			if (p.X == 0 || p.X == (b.Width-1)) && (p.Y == 0 || p.Y == (b.Height-1)) {
				continue
			}

			availableFoodLocations = append(availableFoodLocations, p)
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
func placeFoodRandomlyWithRand(rand Rand, b *BoardState, n int32) error {
	for i := int32(0); i < n; i++ {
		unoccupiedPoints := getUnoccupiedPoints(b, false)
		if len(unoccupiedPoints) > 0 {
			newFood := unoccupiedPoints[rand.Intn(len(unoccupiedPoints))]
			b.Food = append(b.Food, newFood)
		}
	}
	return nil
}
