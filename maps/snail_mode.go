package maps

import (
	"github.com/BattlesnakeOfficial/rules"
)

type SnailModeMap struct{}

// init registers this map in the global registry.
func init() {
	globalRegistry.RegisterMap("snail_mode", SnailModeMap{})
}

// ID returns a unique identifier for this map.
func (m SnailModeMap) ID() string {
	return "snail_mode"
}

// Meta returns the non-functional metadata about this map.
func (m SnailModeMap) Meta() Metadata {
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
func (m SnailModeMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	// Use StandardMap to populate snakes and food
	return StandardMap{}.SetupBoard(initialBoardState, settings, editor)
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

func (m SnailModeMap) PreUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return nil
}

// PostUpdateBoard does the work of placing the hazards along the 'snail tail' of snakes
// This is responsible for saving the current tail location off the board
// and restoring the previous tail position. This also handles removing one hazards from
// the current stacks so the hazards tails fade as the snake moves away.
func (m SnailModeMap) PostUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	err := StandardMap{}.PostUpdateBoard(lastBoardState, settings, editor)
	if err != nil {
		return err
	}

	// This map decrements the stack of hazards on a point each turn, so they
	// need to be cleared first.
	editor.ClearHazards()

	// This is a list of all the hazards we want to add for the previous tails
	// These were stored off board in the previous turn as a way to save state
	// When we add the locations to this list we have already converted the off-board
	// points to on-board points
	tailLocations := make([]rules.Point, 0, len(lastBoardState.Snakes))

	// Count the number of hazards for a given position
	// Add non-double tail locations to a slice
	hazardCounts := map[rules.Point]int{}
	for _, hazard := range lastBoardState.Hazards {

		// discard out of bound
		if outOfBounds(hazard, lastBoardState.Width, lastBoardState.Height) {
			onBoardTail := getPrevTailLocation(hazard, lastBoardState.Height)
			tailLocations = append(tailLocations, onBoardTail)
		} else {
			hazardCounts[hazard]++
		}
	}

	// Add back existing hazards, but with a stack of 1 less than before.
	// This has the effect of making the snail-trail disappear over time.
	for hazard, count := range hazardCounts {

		for i := 0; i < count-1; i++ {
			editor.AddHazard(hazard)
		}
	}

	// Store a stack of hazards for the tail of each snake.  This is stored out
	// of bounds and then applied on the next turn.  The stack count is equal
	// the lenght of the snake.
	for _, snake := range lastBoardState.Snakes {
		if isEliminated(&snake) {
			continue
		}

		// Double tail means that the tail will stay on the same square for more
		// than one turn, so we don't want to spawn hazards
		if doubleTail(&snake) {
			continue
		}

		tail := snake.Body[len(snake.Body)-1]
		offBoardTail := storeTailLocation(tail, lastBoardState.Height)
		for i := 0; i < len(snake.Body); i++ {
			editor.AddHazard(offBoardTail)
		}
	}

	// Read offboard tails and move them to the board. The offboard tails are
	// stacked based on the length of the snake
	for _, p := range tailLocations {

		// Skip position if a snakes head occupies it.
		// Otherwise hazard shows up in the viewer on top of a snake head, but
		// does not damage the snake, which is visually confusing.
		isHead := false
		for _, snake := range lastBoardState.Snakes {
			if isEliminated(&snake) {
				continue
			}
			head := snake.Body[0]
			if p.X == head.X && p.Y == head.Y {
				isHead = true
				break
			}
		}
		if isHead {
			continue
		}

		editor.AddHazard(p)
	}

	return nil
}
