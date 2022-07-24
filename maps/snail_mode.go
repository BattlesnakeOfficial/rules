package maps

import (
	"github.com/BattlesnakeOfficial/rules"
)

type SnailModeMap struct{}

func init() {
	globalRegistry.RegisterMap("snail_mode", SnailModeMap{})
}

func (m SnailModeMap) ID() string {
	return "snail_mode"
}

func (m SnailModeMap) Meta() Metadata {
	return Metadata{
		Name:        "snail_mode",
		Description: "Snakes leave behind a trail of hazards",
		Author:      "Corey and Josh",
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  16,
		BoardSizes:  OddSizes(rules.BoardSizeSmall, rules.BoardSizeXXLarge),
	}
}

func (m SnailModeMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	rand := settings.GetRand(0)

	if len(initialBoardState.Snakes) > int(m.Meta().MaxPlayers) {
		return rules.ErrorTooManySnakes
	}

	snakeIDs := make([]string, 0, len(initialBoardState.Snakes))
	for _, snake := range initialBoardState.Snakes {
		snakeIDs = append(snakeIDs, snake.ID)
	}

	tempBoardState := rules.NewBoardState(initialBoardState.Width, initialBoardState.Height)
	err := rules.PlaceSnakesAutomatically(rand, tempBoardState, snakeIDs)
	if err != nil {
		return err
	}

	// Copy snakes from temp board state
	for _, snake := range tempBoardState.Snakes {
		editor.PlaceSnake(snake.ID, snake.Body, snake.Health)
	}

	return nil
}

func hash(point rules.Point, height int) int {
	return point.X + point.Y*height
}

func dehash(hash int, height int) rules.Point {
	x := hash % height
	y := hash / height

	return rules.Point{X: x, Y: y}
}

func storeTailLocation(point rules.Point, height int) rules.Point {
	return rules.Point{X: point.X, Y: point.Y + height}
}

func getPrevTailLocation(point rules.Point, height int) rules.Point {
	return rules.Point{X: point.X, Y: point.Y - height}
}

func outOfBounds(p rules.Point, w, h int) bool {
	return p.X < 0 || p.Y < 0 || p.X >= w || p.Y >= h
}

func doubleTail(snake *rules.Snake) bool {
	almostTail := snake.Body[len(snake.Body)-2]
	tail := snake.Body[len(snake.Body)-1]
	return almostTail.X == tail.X && almostTail.Y == tail.Y
}

func (m SnailModeMap) UpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	err := StandardMap{}.UpdateBoard(lastBoardState, settings, editor)
	if err != nil {
		return err
	}

	// store hazards so we can decrement stack of hazards by 1
	// we can use the lastBoardState

	editor.ClearHazards()

	tailLocations := make([]rules.Point, 0, len(lastBoardState.Snakes))

	// Count the number of hazards for a given position
	// Add non-double tail locations to a slice
	hazardCounts := map[int]int{}
	for _, hazard := range lastBoardState.Hazards {

		// discard out of bound
		if outOfBounds(hazard, lastBoardState.Width, lastBoardState.Height) {
			onBoardTail := getPrevTailLocation(hazard, lastBoardState.Height)
			tailLocations = append(tailLocations, onBoardTail)
		} else {
			hazardCounts[hash(hazard, lastBoardState.Height)]++
		}
	}

	// Add back existing hazards, but with a stack of 1 less than before
	for hazardHash, count := range hazardCounts {
		hazard := dehash(hazardHash, lastBoardState.Height)

		for i := 0; i < count-1; i++ {
			editor.AddHazard(hazard)
		}
	}

	// ensure stack of 7 for tail segment of each snake
	for _, snake := range lastBoardState.Snakes {

		// Double tail means that the tail will stay on the same square for more
		// than one turn, so we don't want to spawn hazards
		if doubleTail(&snake) {
			continue
		}

		tail := snake.Body[len(snake.Body)-1]
		offBoardTail := storeTailLocation(tail, lastBoardState.Height)
		editor.AddHazard(offBoardTail)
	}

	// read offboard tails and apply 7 hazards
	for _, p := range tailLocations {
		isHead := false
		for _, snake := range lastBoardState.Snakes {
			head := snake.Body[0]
			if p.X == head.X && p.Y == head.Y {
				isHead = true
				break
			}
		}
		if isHead {
			continue
		}
		for i := 0; i < 7; i++ {
			editor.AddHazard(p)
		}
	}

	return nil
}
