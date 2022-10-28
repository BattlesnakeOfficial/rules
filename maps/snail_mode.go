package maps

import (
	"github.com/BattlesnakeOfficial/rules"
)

type SnailModeMap struct {
	lastTailPositions map[rules.Point]int // local state is preserved during the turn
}

// init registers this map in the global registry.
func init() {
	globalRegistry.RegisterMap("snail_mode", &SnailModeMap{lastTailPositions: nil})
}

// ID returns a unique identifier for this map.
func (m *SnailModeMap) ID() string {
	return "snail_mode"
}

// Meta returns the non-functional metadata about this map.
func (m *SnailModeMap) Meta() Metadata {
	return Metadata{
		Name:        "Snail Mode",
		Description: "Snakes leave behind a trail of hazards",
		Author:      "coreyja and jlafayette",
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  16,
		BoardSizes:  OddSizes(rules.BoardSizeSmall, rules.BoardSizeXXLarge),
		Tags:        []string{TAG_EXPERIMENTAL, TAG_HAZARD_PLACEMENT},
	}
}

// SetupBoard here is pretty 'standard' and doesn't do any special setup for this game mode
func (m *SnailModeMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
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

// storeTailLocation returns an offboard point that corresponds to the given point.
// This is useful for storing state that can be accessed next turn.
func storeTailLocation(point rules.Point, height int) rules.Point {
	return rules.Point{X: point.X, Y: point.Y + height}
}

// getPrevTailLocation returns the onboard point that corresponds to an offboard point.
// This is useful for restoring state that was stored last turn.
func getPrevTailLocation(point rules.Point, height int) rules.Point {
	return rules.Point{X: point.X, Y: point.Y - height}
}

// outOfBounds determines if the given point is out of bounds for the current board size
func outOfBounds(p rules.Point, w, h int) bool {
	return p.X < 0 || p.Y < 0 || p.X >= w || p.Y >= h
}

// doubleTail determine if the snake has a double stacked tail currently
func doubleTail(snake *rules.Snake) bool {
	almostTail := snake.Body[len(snake.Body)-2]
	tail := snake.Body[len(snake.Body)-1]
	return almostTail.X == tail.X && almostTail.Y == tail.Y
}

// isEliminated determines if the snake is already eliminated
func isEliminated(s *rules.Snake) bool {
	return s.EliminatedCause != rules.NotEliminated
}

// PreUpdateBoard stores the tail position of each snake in memory, to be
// able to place hazards there after the snakes move.
func (m *SnailModeMap) PreUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	m.lastTailPositions = make(map[rules.Point]int)
	for _, snake := range lastBoardState.Snakes {
		if isEliminated(&snake) {
			continue
		}
		// Double tail means that the tail will stay on the same square for more
		// than one turn, so we don't want to spawn hazards
		if doubleTail(&snake) {
			continue
		}
		m.lastTailPositions[snake.Body[len(snake.Body)-1]] = len(snake.Body)
	}
	return nil
}

// PostUpdateBoard does the work of placing the hazards along the 'snail tail' of snakes
// This also handles removing one hazards from the current stacks so the hazards tails fade as the snake moves away.
func (m *SnailModeMap) PostUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	err := StandardMap{}.PostUpdateBoard(lastBoardState, settings, editor)
	if err != nil {
		return err
	}

	// This map decrements the stack of hazards on a point each turn, so they
	// need to be cleared first.
	editor.ClearHazards()

	// Count the number of hazards for a given position
	hazardCounts := map[rules.Point]int{}
	for _, hazard := range lastBoardState.Hazards {
		hazardCounts[hazard]++
	}

	// Add back existing hazards, but with a stack of 1 less than before.
	// This has the effect of making the snail-trail disappear over time.
	for hazard, count := range hazardCounts {
		for i := 0; i < count-1; i++ {
			editor.AddHazard(hazard)
		}
	}

	// Place a new stack of hazards where each snake's tail used to be
NewHazardLoop:
	for location, count := range m.lastTailPositions {
		for _, snake := range lastBoardState.Snakes {
			if isEliminated(&snake) {
				continue
			}
			head := snake.Body[0]
			if location.X == head.X && location.Y == head.Y {
				// Skip position if a snakes head occupies it.
				// Otherwise hazard shows up in the viewer on top of a snake head, but
				// does not damage the snake, which is visually confusing.
				continue NewHazardLoop
			}
		}
		for i := 0; i < count; i++ {
			editor.AddHazard(location)
		}
	}

	return nil
}
