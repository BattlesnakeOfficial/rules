package rules

import (
	"errors"
	"math/rand"
)

type StandardRuleset struct{}

const (
	BoardSizeSmall  = 7
	BoardSizeMedium = 11
	BoardSizeLarge  = 19

	FoodSpawnChance = 0.15

	SnakeMaxHealth = 100
	SnakeStartSize = 3

	// bvanvugt - TODO: Just return formatted strings instead of codes?
	NotEliminated                   = ""
	EliminatedByCollision           = "snake-collision"
	EliminatedBySelfCollision       = "snake-self-collision"
	EliminatedByStarvation          = "starvation"
	EliminatedByHeadToHeadCollision = "head-collision"
	EliminatedByOutOfBounds         = "wall-collision"

	// TODO - Error consts
)

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
		Point{mn, mn},
		Point{mn, md},
		Point{mn, mx},
		Point{md, mn},
		Point{md, mx},
		Point{mx, mn},
		Point{mx, md},
		Point{mx, mx},
	}

	// Sanity check
	if len(b.Snakes) > len(startPoints) {
		return errors.New("too many snakes for fixed start positions")
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
		unoccupiedPoints := r.getUnoccupiedPoints(b)
		if len(unoccupiedPoints) <= 0 {
			return errors.New("not enough space to place snake")
		}
		p := unoccupiedPoints[rand.Intn(len(unoccupiedPoints))]
		for j := 0; j < SnakeStartSize; j++ {
			b.Snakes[i].Body = append(b.Snakes[i].Body, p)
		}
	}
	return nil
}

func (r *StandardRuleset) placeFood(b *BoardState) error {
	return r.spawnFood(b, len(b.Snakes))
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

func (r *StandardRuleset) ResolveMoves(prevState *BoardState, moves []SnakeMove) (*BoardState, error) {
	// We specifically want to copy prevState, so as not to alter it directly.
	nextState := &BoardState{
		Height: prevState.Height,
		Width:  prevState.Width,
		Food:   append([]Point{}, prevState.Food...),
		Snakes: make([]Snake, len(prevState.Snakes)),
	}
	for i := 0; i < len(prevState.Snakes); i++ {
		nextState.Snakes[i].ID = prevState.Snakes[i].ID
		nextState.Snakes[i].Health = prevState.Snakes[i].Health
		nextState.Snakes[i].Body = append([]Point{}, prevState.Snakes[i].Body...)
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

	// TODO: LOG?
	err = r.eliminateSnakes(nextState)
	if err != nil {
		return nil, err
	}

	// TODO
	// bvanvugt: we specifically want this to happen before elimination
	// so that head-to-head collisions on food still remove the food.
	// It does create an artifact though, where head-to-head collisions
	// of equal length actually show length + 1

	// TODO: LOG?
	err = r.feedSnakes(nextState)
	if err != nil {
		return nil, err
	}

	// TODO: LOG?
	err = r.maybeSpawnFood(nextState, 1)
	if err != nil {
		return nil, err
	}

	return nextState, nil
}

func (r *StandardRuleset) moveSnakes(b *BoardState, moves []SnakeMove) error {
	if len(moves) < len(b.Snakes) {
		return errors.New("not enough snake moves")
	}
	if len(moves) > len(b.Snakes) {
		return errors.New("too many snake moves")
	}

	for _, move := range moves {
		var snake *Snake
		for i := 0; i < len(b.Snakes); i++ {
			if b.Snakes[i].ID == move.ID {
				snake = &b.Snakes[i]
			}
		}
		if snake == nil {
			return errors.New("snake not found for move")
		}

		// Do not move eliminated snakes
		if snake.EliminatedCause != NotEliminated {
			continue
		}

		var newHead = Point{}
		switch move.Move {
		case MoveDown:
			newHead.X = snake.Body[0].X
			newHead.Y = snake.Body[0].Y + 1
		case MoveLeft:
			newHead.X = snake.Body[0].X - 1
			newHead.Y = snake.Body[0].Y
		case MoveRight:
			newHead.X = snake.Body[0].X + 1
			newHead.Y = snake.Body[0].Y
		case MoveUp:
			newHead.X = snake.Body[0].X
			newHead.Y = snake.Body[0].Y - 1
		default:
			// Default to UP
			var dX int32 = 0
			var dY int32 = -1
			// If neck is available, use neck to determine last direction
			if len(snake.Body) >= 2 {
				dX = snake.Body[0].X - snake.Body[1].X
				dY = snake.Body[0].Y - snake.Body[1].Y
				if dX == 0 && dY == 0 {
					dY = -1 // Move up if no last move was made
				}
			}
			// Apply
			newHead.X = snake.Body[0].X + dX
			newHead.Y = snake.Body[0].Y + dY
		}

		// Append new head, pop old tail
		snake.Body = append([]Point{newHead}, snake.Body[:len(snake.Body)-1]...)
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

func (r *StandardRuleset) eliminateSnakes(b *BoardState) error {
	for i := 0; i < len(b.Snakes); i++ {
		snake := &b.Snakes[i]
		if len(snake.Body) <= 0 {
			return errors.New("snake is length zero")
		}

		if r.snakeHasStarved(snake) {
			snake.EliminatedCause = EliminatedByStarvation
			continue
		}

		if r.snakeIsOutOfBounds(snake, b.Width, b.Height) {
			snake.EliminatedCause = EliminatedByOutOfBounds
			continue
		}

		// Always check body collisions before head-to-heads
		for j := 0; j < len(b.Snakes); j++ {
			other := &b.Snakes[j]
			if r.snakeHasBodyCollided(snake, other) {
				if snake.ID == other.ID {
					snake.EliminatedCause = EliminatedBySelfCollision
				} else {
					snake.EliminatedCause = EliminatedByCollision
				}
				break
			}
		}
		if snake.EliminatedCause != NotEliminated {
			continue
		}

		// Always check body collisions before head-to-heads
		for j := 0; j < len(b.Snakes); j++ {
			other := &b.Snakes[j]
			if snake.ID != other.ID && r.snakeHasLostHeadToHead(snake, other) {
				snake.EliminatedCause = EliminatedByHeadToHeadCollision
				break
			}
		}
	}
	return nil
}

func (r *StandardRuleset) snakeHasStarved(s *Snake) bool {
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

func (r *StandardRuleset) feedSnakes(b *BoardState) error {
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
				foodHasBeenEaten = true
				// Update snake
				snake.Body = append(snake.Body, snake.Body[len(snake.Body)-1])
				snake.Health = SnakeMaxHealth
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

func (r *StandardRuleset) maybeSpawnFood(b *BoardState, n int) error {
	if rand.Float32() <= FoodSpawnChance {
		return r.spawnFood(b, n)
	}
	return nil
}

func (r *StandardRuleset) spawnFood(b *BoardState, n int) error {
	for i := 0; i < n; i++ {
		unoccupiedPoints := r.getUnoccupiedPoints(b)
		if len(unoccupiedPoints) > 0 {
			newFood := unoccupiedPoints[rand.Intn(len(unoccupiedPoints))]
			b.Food = append(b.Food, newFood)
		}
	}
	return nil
}

func (r *StandardRuleset) getUnoccupiedPoints(b *BoardState) []Point {
	pointIsOccupied := map[int32]map[int32]bool{}
	for _, p := range b.Food {
		if _, xExists := pointIsOccupied[p.X]; !xExists {
			pointIsOccupied[p.X] = map[int32]bool{}
		}
		pointIsOccupied[p.X][p.Y] = true
	}
	for _, snake := range b.Snakes {
		for _, p := range snake.Body {
			if _, xExists := pointIsOccupied[p.X]; !xExists {
				pointIsOccupied[p.X] = map[int32]bool{}
			}
			pointIsOccupied[p.X][p.Y] = true
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
