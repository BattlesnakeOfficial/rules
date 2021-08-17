package rules

import (
	"math/rand"
	"sort"
)

type StandardRuleset struct {
	FoodSpawnChance     int32 // [0, 100]
	MinimumFood         int32
	HazardDamagePerTurn int32
}

func (r *StandardRuleset) Name() string { return "standard" }

func (r *StandardRuleset) CreateInitialBoardState(width int32, height int32, snakeIDs []string) (*BoardState, error) {
	initialBoardState := &BoardState{
		Height: height,
		Width:  width,
		Snakes: make([]Snake, len(snakeIDs)),
	}

	for i := 0; i < len(snakeIDs); i++ {
		initialBoardState.Snakes[i] = Snake{
			ID:     snakeIDs[i],
			Health: SnakeMaxHealth,
		}
	}

	err := r.placeSnakes(initialBoardState)
	if err != nil {
		return nil, err
	}

	err = r.placeFood(initialBoardState)
	if err != nil {
		return nil, err
	}

	return initialBoardState, nil
}

func (r *StandardRuleset) placeSnakes(b *BoardState) error {
	if r.isKnownBoardSize(b) {
		return r.placeSnakesFixed(b)
	}
	return r.placeSnakesRandomly(b)
}

func (r *StandardRuleset) placeSnakesFixed(b *BoardState) error {
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

func (r *StandardRuleset) placeSnakesRandomly(b *BoardState) error {

	for i := 0; i < len(b.Snakes); i++ {
		unoccupiedPoints := r.getEvenUnoccupiedPoints(b)
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

func (r *StandardRuleset) placeFood(b *BoardState) error {
	if r.isKnownBoardSize(b) {
		return r.placeFoodFixed(b)
	}
	return r.placeFoodRandomly(b)
}

func (r *StandardRuleset) placeFoodFixed(b *BoardState) error {
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
	unoccupiedPoints := r.getUnoccupiedPoints(b, true)
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

func (r *StandardRuleset) placeFoodRandomly(b *BoardState) error {
	return r.spawnFood(b, int32(len(b.Snakes)))
}

func (r *StandardRuleset) isKnownBoardSize(b *BoardState) bool {
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

func (r *StandardRuleset) CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
	// We specifically want to copy prevState, so as not to alter it directly.
	nextState := &BoardState{
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

	// TODO: Gut check the BoardState?

	// TODO: LOG?
	err := r.moveSnakes(nextState, moves)
	if err != nil {
		return nil, err
	}

	// TODO: LOG?
	err = r.reduceSnakeHealth(nextState)
	if err != nil {
		return nil, err
	}

	err = r.maybeDamageHazards(nextState)
	if err != nil {
		return nil, err
	}

	// TODO: LOG?
	// bvanvugt: We specifically want this to happen before elimination for two reasons:
	// 1) We want snakes to be able to eat on their very last turn and still survive.
	// 2) So that head-to-head collisions on food still remove the food.
	//    This does create an artifact though, where head-to-head collisions
	//    of equal length actually show length + 1 and full health, as if both snakes ate.
	err = r.maybeFeedSnakes(nextState)
	if err != nil {
		return nil, err
	}

	// TODO: LOG?
	err = r.maybeSpawnFood(nextState)
	if err != nil {
		return nil, err
	}

	// TODO: LOG?
	err = r.maybeEliminateSnakes(nextState)
	if err != nil {
		return nil, err
	}

	return nextState, nil
}

func (r *StandardRuleset) moveSnakes(b *BoardState, moves []SnakeMove) error {
	// Sanity check that all non-eliminated snakes have moves and bodies.
	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause != NotEliminated {
			continue
		}

		if len(snake.Body) == 0 {
			return ErrorZeroLengthSnake
		}
		moveFound := false
		for _, move := range moves {
			if snake.ID == move.ID {
				moveFound = true
				break
			}
		}
		if !moveFound {
			return ErrorNoMoveFound
		}
	}

	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause != NotEliminated {
			continue
		}

		for _, move := range moves {
			if move.ID == snake.ID {
				var newHead = Point{}
				switch move.Move {
				case MoveDown:
					newHead.X = snake.Body[0].X
					newHead.Y = snake.Body[0].Y - 1
				case MoveLeft:
					newHead.X = snake.Body[0].X - 1
					newHead.Y = snake.Body[0].Y
				case MoveRight:
					newHead.X = snake.Body[0].X + 1
					newHead.Y = snake.Body[0].Y
				case MoveUp:
					newHead.X = snake.Body[0].X
					newHead.Y = snake.Body[0].Y + 1
				default:
					// Default to UP
					var dX int32 = 0
					var dY int32 = 1
					// If neck is available, use neck to determine last direction
					if len(snake.Body) >= 2 {
						dX = snake.Body[0].X - snake.Body[1].X
						dY = snake.Body[0].Y - snake.Body[1].Y
						if dX == 0 && dY == 0 {
							dY = 1 // Move up if no last move was made
						}
					}
					// Apply
					newHead.X = snake.Body[0].X + dX
					newHead.Y = snake.Body[0].Y + dY
				}

				// Append new head, pop old tail
				snake.Body = append([]Point{newHead}, snake.Body[:len(snake.Body)-1]...)
			}
		}
	}
	return nil
}

func (r *StandardRuleset) reduceSnakeHealth(b *BoardState) error {
	for i := 0; i < len(b.Snakes); i++ {
		if b.Snakes[i].EliminatedCause == NotEliminated {
			b.Snakes[i].Health = b.Snakes[i].Health - 1
		}
	}
	return nil
}

func (r *StandardRuleset) maybeDamageHazards(b *BoardState) error {
	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause != NotEliminated {
			continue
		}
		head := snake.Body[0]
		for _, p := range b.Hazards {
			if head == p {
				// If there's a food in this square, don't reduce health
				foundFood := false
				for _, food := range b.Food {
					if p == food {
						foundFood = true
					}
				}
				if foundFood {
					continue
				}

				// Snake is in a hazard, reduce health
				snake.Health = snake.Health - r.HazardDamagePerTurn
				if snake.Health < 0 {
					snake.Health = 0
				}
				if r.snakeIsOutOfHealth(snake) {
					snake.EliminatedCause = EliminatedByOutOfHealth
				}
			}
		}
	}

	return nil
}

func (r *StandardRuleset) maybeEliminateSnakes(b *BoardState) error {
	// First order snake indices by length.
	// In multi-collision scenarios we want to always attribute elimination to the longest snake.
	snakeIndicesByLength := make([]int, len(b.Snakes))
	for i := 0; i < len(b.Snakes); i++ {
		snakeIndicesByLength[i] = i
	}
	sort.Slice(snakeIndicesByLength, func(i int, j int) bool {
		lenI := len(b.Snakes[snakeIndicesByLength[i]].Body)
		lenJ := len(b.Snakes[snakeIndicesByLength[j]].Body)
		return lenI > lenJ
	})

	// First, iterate over all non-eliminated snakes and eliminate the ones
	// that are out of health or have moved out of bounds.
	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause != NotEliminated {
			continue
		}
		if len(snake.Body) <= 0 {
			return ErrorZeroLengthSnake
		}

		if r.snakeIsOutOfHealth(snake) {
			snake.EliminatedCause = EliminatedByOutOfHealth
			continue
		}

		if r.snakeIsOutOfBounds(snake, b.Width, b.Height) {
			snake.EliminatedCause = EliminatedByOutOfBounds
			continue
		}
	}

	// Next, look for any collisions. Note we apply collision eliminations
	// after this check so that snakes can collide with each other and be properly eliminated.
	type CollisionElimination struct {
		ID    string
		Cause string
		By    string
	}
	collisionEliminations := []CollisionElimination{}
	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if snake.EliminatedCause != NotEliminated {
			continue
		}
		if len(snake.Body) <= 0 {
			return ErrorZeroLengthSnake
		}

		// Check for self-collisions first
		if r.snakeHasBodyCollided(snake, snake) {
			collisionEliminations = append(collisionEliminations, CollisionElimination{
				ID:    snake.ID,
				Cause: EliminatedBySelfCollision,
				By:    snake.ID,
			})
			continue
		}

		// Check for body collisions with other snakes second
		hasBodyCollided := false
		for _, otherIndex := range snakeIndicesByLength {
			other := &b.Snakes[otherIndex]
			if other.EliminatedCause != NotEliminated {
				continue
			}
			if snake.ID != other.ID && r.snakeHasBodyCollided(snake, other) {
				collisionEliminations = append(collisionEliminations, CollisionElimination{
					ID:    snake.ID,
					Cause: EliminatedByCollision,
					By:    other.ID,
				})
				hasBodyCollided = true
				break
			}
		}
		if hasBodyCollided {
			continue
		}

		// Check for head-to-heads last
		hasHeadCollided := false
		for _, otherIndex := range snakeIndicesByLength {
			other := &b.Snakes[otherIndex]
			if other.EliminatedCause != NotEliminated {
				continue
			}
			if snake.ID != other.ID && r.snakeHasLostHeadToHead(snake, other) {
				collisionEliminations = append(collisionEliminations, CollisionElimination{
					ID:    snake.ID,
					Cause: EliminatedByHeadToHeadCollision,
					By:    other.ID,
				})
				hasHeadCollided = true
				break
			}
		}
		if hasHeadCollided {
			continue
		}
	}

	// Apply collision eliminations
	for _, elimination := range collisionEliminations {
		for i := 0; i < len(b.Snakes); i++ {
			snake := &b.Snakes[i]
			if snake.ID == elimination.ID {
				snake.EliminatedCause = elimination.Cause
				snake.EliminatedBy = elimination.By
				break
			}
		}
	}

	return nil
}

func (r *StandardRuleset) snakeIsOutOfHealth(s *Snake) bool {
	return s.Health <= 0
}

func (r *StandardRuleset) snakeIsOutOfBounds(s *Snake, boardWidth int32, boardHeight int32) bool {
	for _, point := range s.Body {
		if (point.X < 0) || (point.X >= boardWidth) {
			return true
		}
		if (point.Y < 0) || (point.Y >= boardHeight) {
			return true
		}
	}
	return false
}

func (r *StandardRuleset) snakeHasBodyCollided(s *Snake, other *Snake) bool {
	head := s.Body[0]
	for i, body := range other.Body {
		if i == 0 {
			continue
		} else if head.X == body.X && head.Y == body.Y {
			return true
		}
	}
	return false
}

func (r *StandardRuleset) snakeHasLostHeadToHead(s *Snake, other *Snake) bool {
	if s.Body[0].X == other.Body[0].X && s.Body[0].Y == other.Body[0].Y {
		return len(s.Body) <= len(other.Body)
	}
	return false
}

func (r *StandardRuleset) maybeFeedSnakes(b *BoardState) error {
	newFood := []Point{}
	for _, food := range b.Food {
		foodHasBeenEaten := false
		for i := 0; i < len(b.Snakes); i++ {
			snake := &b.Snakes[i]

			// Ignore eliminated and zero-length snakes, they can't eat.
			if snake.EliminatedCause != NotEliminated || len(snake.Body) == 0 {
				continue
			}

			if snake.Body[0].X == food.X && snake.Body[0].Y == food.Y {
				r.feedSnake(snake)
				foodHasBeenEaten = true
			}
		}
		// Persist food to next BoardState if not eaten
		if !foodHasBeenEaten {
			newFood = append(newFood, food)
		}
	}

	b.Food = newFood
	return nil
}

func (r *StandardRuleset) feedSnake(snake *Snake) {
	r.growSnake(snake)
	snake.Health = SnakeMaxHealth
}

func (r *StandardRuleset) growSnake(snake *Snake) {
	if len(snake.Body) > 0 {
		snake.Body = append(snake.Body, snake.Body[len(snake.Body)-1])
	}
}

func (r *StandardRuleset) maybeSpawnFood(b *BoardState) error {
	numCurrentFood := int32(len(b.Food))
	if numCurrentFood < r.MinimumFood {
		return r.spawnFood(b, r.MinimumFood-numCurrentFood)
	} else if r.FoodSpawnChance > 0 && int32(rand.Intn(100)) < r.FoodSpawnChance {
		return r.spawnFood(b, 1)
	}
	return nil
}

func (r *StandardRuleset) spawnFood(b *BoardState, n int32) error {
	for i := int32(0); i < n; i++ {
		unoccupiedPoints := r.getUnoccupiedPoints(b, false)
		if len(unoccupiedPoints) > 0 {
			newFood := unoccupiedPoints[rand.Intn(len(unoccupiedPoints))]
			b.Food = append(b.Food, newFood)
		}
	}
	return nil
}

func (r *StandardRuleset) getUnoccupiedPoints(b *BoardState, includePossibleMoves bool) []Point {
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

func (r *StandardRuleset) getEvenUnoccupiedPoints(b *BoardState) []Point {
	// Start by getting unoccupied points
	unoccupiedPoints := r.getUnoccupiedPoints(b, true)

	// Create a new array to hold points that are  even
	evenUnoccupiedPoints := []Point{}

	for _, point := range unoccupiedPoints {
		if ((point.X + point.Y) % 2) == 0 {
			evenUnoccupiedPoints = append(evenUnoccupiedPoints, point)
		}
	}
	return evenUnoccupiedPoints
}

func (r *StandardRuleset) IsGameOver(b *BoardState) (bool, error) {
	numSnakesRemaining := 0
	for i := 0; i < len(b.Snakes); i++ {
		if b.Snakes[i].EliminatedCause == NotEliminated {
			numSnakesRemaining++
		}
	}
	return numSnakesRemaining <= 1, nil
}
