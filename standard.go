package rulesets

import (
	"errors"
	"math/rand"
)

type StandardRuleset struct{}

const (
	BOARD_SIZE_SMALL  = 7
	BOARD_SIZE_MEDIUM = 11
	BOARD_SIZE_LARGE  = 19
	FOOD_SPAWN_CHANCE = 0.1
	SNAKE_MAX_HEALTH  = 100

	// bvanvugt - TODO: Just return formatted strings instead of codes?
	ELIMINATED_COLLISION      = "snake-collision"
	ELIMINATED_SELF_COLLISION = "snake-self-collision"
	ELIMINATED_STARVATION     = "starvation"
	ELIMINATED_HEAD_TO_HEAD   = "head-collision"
	ELIMINATED_OUT_OF_BOUNDS  = "wall-collision"
)

func (r *StandardRuleset) CreateInitialBoard(width int32, height int32, snakeIDs []string) (*BoardState, error) {
	var err error

	snakes := []*Snake{}
	for _, id := range snakeIDs {
		snakes = append(snakes,
			&Snake{
				ID:     id,
				Health: SNAKE_MAX_HEALTH,
			},
		)
	}

	initialBoardState := &BoardState{
		Height: height,
		Width:  width,
		Snakes: snakes,
	}

	// Place Snakes
	if r.isKnownBoardSize(initialBoardState) {
		err = r.placeSnakesFixed(initialBoardState)
	} else {
		err = r.placeSnakesRandomly(initialBoardState)
	}
	if err != nil {
		return nil, err
	}

	// Place Food
	err = r.placeInitialFood(initialBoardState)
	if err != nil {
		return nil, err
	}

	return initialBoardState, nil
}

func (r *StandardRuleset) placeSnakesFixed(b *BoardState) error {
	// Sanity check
	if len(b.Snakes) >= 8 {
		return errors.New("too many snakes for fixed start positions")
	}

	// Create start points
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

	// Randomly order them
	rand.Shuffle(len(startPoints), func(i int, j int) {
		startPoints[i], startPoints[j] = startPoints[j], startPoints[i]
	})

	// Assign to snakes in order given
	for i, snake := range b.Snakes {
		p := startPoints[i]
		for j := 0; j < 3; j++ {
			snake.Body = append(snake.Body, &Point{p.X, p.Y})
		}

	}
	return nil
}

func (r *StandardRuleset) placeSnakesRandomly(b *BoardState) error {
	for _, snake := range b.Snakes {
		unoccupiedPoints := r.getUnoccupiedPoints(b)
		p := unoccupiedPoints[rand.Intn(len(unoccupiedPoints))]
		for j := 0; j < 3; j++ {
			snake.Body = append(snake.Body, &Point{p.X, p.Y})
		}
	}
	return nil
}

func (r *StandardRuleset) isKnownBoardSize(b *BoardState) bool {
	if b.Height == BOARD_SIZE_SMALL && b.Width == BOARD_SIZE_SMALL {
		return true
	}
	if b.Height == BOARD_SIZE_MEDIUM && b.Width == BOARD_SIZE_MEDIUM {
		return true
	}
	if b.Height == BOARD_SIZE_LARGE && b.Width == BOARD_SIZE_LARGE {
		return true
	}
	return false
}

func (r *StandardRuleset) placeInitialFood(b *BoardState) error {
	r.spawnFood(b, len(b.Snakes))
	return nil
}

func (r *StandardRuleset) ResolveMoves(prevState *BoardState, moves []*SnakeMove) (*BoardState, error) {
	// TODO: DO NOT REFERENCE prevState directly!!!!
	// we're technically altering both states
	nextState := &BoardState{
		Snakes: prevState.Snakes,
		Food:   prevState.Food,
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

	// TODO: LOG?
	err = r.eliminateSnakes(nextState)
	if err != nil {
		return nil, err
	}

	return nextState, nil
}

func (r *StandardRuleset) moveSnakes(b *BoardState, moves []*SnakeMove) error {
	for _, move := range moves {
		var newHead = &Point{}
		switch move.Move {
		case MOVE_DOWN:
			newHead.X = move.Snake.Body[0].X
			newHead.Y = move.Snake.Body[0].Y + 1
		case MOVE_LEFT:
			newHead.X = move.Snake.Body[0].X - 1
			newHead.Y = move.Snake.Body[0].Y
		case MOVE_RIGHT:
			newHead.X = move.Snake.Body[0].X + 1
			newHead.Y = move.Snake.Body[0].Y
		case MOVE_UP:
			newHead.X = move.Snake.Body[0].X
			newHead.Y = move.Snake.Body[0].Y - 1
		default:
			// Default to UP
			var dX int32 = 0
			var dY int32 = -1
			// If neck is available, use neck to determine last direction
			if len(move.Snake.Body) >= 2 {
				dX = move.Snake.Body[0].X - move.Snake.Body[1].X
				dY = move.Snake.Body[0].Y - move.Snake.Body[1].Y
				if dX == 0 && dY == 0 {
					dY = -1 // Move up if no last move was made
				}
			}
			// Apply
			newHead.X = move.Snake.Body[0].X + dX
			newHead.Y = move.Snake.Body[0].Y + dY
		}

		// Append new head, pop old tail
		move.Snake.Body = append([]*Point{newHead}, move.Snake.Body[:len(move.Snake.Body)-1]...)
	}
	return nil
}

func (r *StandardRuleset) reduceSnakeHealth(b *BoardState) error {
	for _, snake := range b.Snakes {
		snake.Health = snake.Health - 1
	}
	return nil
}

func (r *StandardRuleset) eliminateSnakes(b *BoardState) error {
	for _, snake := range b.Snakes {
		if r.snakeHasStarved(snake) {
			snake.EliminatedCause = ELIMINATED_STARVATION
		} else if r.snakeIsOutOfBounds(snake, b.Width, b.Height) {
			snake.EliminatedCause = ELIMINATED_OUT_OF_BOUNDS
		} else {
			for _, other := range b.Snakes {
				if r.snakeHasBodyCollided(snake, other) {
					if snake.ID == other.ID {
						snake.EliminatedCause = ELIMINATED_SELF_COLLISION
					} else {
						snake.EliminatedCause = ELIMINATED_COLLISION
					}
					break
				} else if r.snakeHasLostHeadToHead(snake, other) {
					snake.EliminatedCause = ELIMINATED_HEAD_TO_HEAD
					break
				}
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
	var newFood []*Point
	var tail *Point

	for _, food := range b.Food {
		foodHasBeenEaten := false
		for _, snake := range b.Snakes {
			if snake.Body[0].X == food.X && snake.Body[0].Y == food.Y {
				foodHasBeenEaten = true
				// Update snake
				snake.Health = SNAKE_MAX_HEALTH
				tail = snake.Body[len(snake.Body)-1]
				snake.Body = append(snake.Body, &Point{X: tail.X, Y: tail.Y})
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
	if rand.Float32() <= FOOD_SPAWN_CHANCE {
		r.spawnFood(b, n)
	}
	return nil
}

func (r *StandardRuleset) spawnFood(b *BoardState, n int) {
	for i := 0; i < n; i++ {
		unoccupiedPoints := r.getUnoccupiedPoints(b)
		if len(unoccupiedPoints) > 0 {
			newFood := unoccupiedPoints[rand.Intn(len(unoccupiedPoints))]
			b.Food = append(b.Food, newFood)
		}
	}
}

func (r *StandardRuleset) getUnoccupiedPoints(b *BoardState) []*Point {
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

	unoccupiedPoints := []*Point{}
	for x := int32(0); x < b.Width; x++ {
		for y := int32(0); y < b.Height; y++ {
			if _, xExists := pointIsOccupied[x]; xExists {
				if isOccupied, yExists := pointIsOccupied[x][y]; yExists {
					if isOccupied {
						continue
					}
				}
			}
			unoccupiedPoints = append(unoccupiedPoints, &Point{X: x, Y: y})
		}
	}
	return unoccupiedPoints
}
