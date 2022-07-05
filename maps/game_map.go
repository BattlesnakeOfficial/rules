package maps

import (
	"github.com/BattlesnakeOfficial/rules"
)

type GameMap interface {
	// Return a unique identifier for this map.
	ID() string

	// Return non-functional metadata about this map.
	Meta() Metadata

	// Called to generate a new board. The map is responsible for placing all snakes, food, and hazards.
	SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error

	// Called every turn to optionally update the board.
	UpdateBoard(previousBoardState *rules.BoardState, settings rules.Settings, editor Editor) error
}

// Dimensions describes the size of a Battlesnake board.
type Dimensions struct {
	// Width is the width, in number of board squares, of the board.
	// The value 0 has a special meaning to mean unlimited.
	Width uint
	// Height is the height, in number of board squares, of the board.
	// The value 0 has a special meaning to mean unlimited.
	Height uint
}

// sizes is a list of board sizes that a map supports.
type sizes []Dimensions

// IsUnlimited reports whether the supported sizes are unlimited.
// Note that even for unlimited sizes, there will be an upper bound that can actually be run and visualised.
func (d sizes) IsUnlimited() bool {
	return len(d) == 1 && d[0].Width == 0
}

// AnySize creates sizes for a board that has no fixed sizes (supports unlimited sizes).
func AnySize() sizes {
	return sizes{Dimensions{Width: 0, Height: 0}}
}

// OddSquareSizes generates square (width = height) board sizes with an odd number of positions
// in the vertical and horizontal directions.
// Examples:
//  - OddSquareSizes(11,21) produces [(11,11), (13,13), (15,15), (17,17), (19,19), (21,21)]
func OddSquareSizes(min, max uint) sizes {
	var s sizes
	for i := min; i <= max; i += 2 {
		s = append(s, Dimensions{Width: i, Height: i})
	}

	return s
}

// FixedSizes creates dimensions for a board that has 1 or more fixed sizes.
// Examples:
// - FixedSizes(Dimension{9,11}) supports only a width of 9 and a height of 11.
// - FixedSizes(Dimensions{11,11},Dimensions{19,19}) supports sizes 11x11 and 19x19
func FixedSizes(a Dimensions, b ...Dimensions) sizes {
	s := make(sizes, 0, 1+len(b))
	s = append(s, a)
	s = append(s, b...)
	return s
}

type Metadata struct {
	Name        string
	Author      string
	Description string
	// Version is the current version of the game map.
	// Each time a map is changed, the version number should be incremented by 1.
	Version uint
	// MinPlayers is the minimum number of players that the map supports.
	MinPlayers uint
	// MaxPlayers is the maximum number of players that the map supports.
	MaxPlayers uint
	// BoardSizes is a list of supported board sizes. Board sizes can fall into one of 3 categories:
	//   1. one fixed size (i.e. [11x11])
	//   2. multiple, fixed sizes (i.e. [11x11, 19x19, 25x25])
	//   3. "unlimited" sizes (the board is not fixed and can scale to any reasonable size)
	BoardSizes sizes
}

// Editor is used by GameMap implementations to modify the board state.
type Editor interface {
	// Clears all food from the board.
	ClearFood()

	// Clears all hazards from the board.
	ClearHazards()

	// Adds a food to the board. Does not check for duplicates.
	AddFood(rules.Point)

	// Adds a hazard to the board. Does not check for duplicates.
	AddHazard(rules.Point)

	// Removes all food from a specific tile on the board.
	RemoveFood(rules.Point)

	// Removes all hazards from a specific tile on the board.
	RemoveHazard(rules.Point)

	// Updates the body and health of a snake.
	PlaceSnake(id string, body []rules.Point, health int)
}

// An Editor backed by a BoardState.
type BoardStateEditor struct {
	*rules.BoardState
}

func NewBoardStateEditor(boardState *rules.BoardState) *BoardStateEditor {
	return &BoardStateEditor{
		BoardState: boardState,
	}
}

func (editor *BoardStateEditor) ClearFood() {
	editor.Food = []rules.Point{}
}

func (editor *BoardStateEditor) ClearHazards() {
	editor.Hazards = []rules.Point{}
}

func (editor *BoardStateEditor) AddFood(p rules.Point) {
	editor.Food = append(editor.Food, rules.Point{X: p.X, Y: p.Y})
}

func (editor *BoardStateEditor) AddHazard(p rules.Point) {
	editor.Hazards = append(editor.Hazards, rules.Point{X: p.X, Y: p.Y})
}

func (editor *BoardStateEditor) RemoveFood(p rules.Point) {
	for index, food := range editor.Food {
		if food.X == p.X && food.Y == p.Y {
			editor.Food[index] = editor.Food[len(editor.Food)-1]
			editor.Food = editor.Food[:len(editor.Food)-1]
		}
	}
}

func (editor *BoardStateEditor) RemoveHazard(p rules.Point) {
	for index, food := range editor.Hazards {
		if food.X == p.X && food.Y == p.Y {
			editor.Hazards[index] = editor.Hazards[len(editor.Hazards)-1]
			editor.Hazards = editor.Hazards[:len(editor.Hazards)-1]
		}
	}
}

func (editor *BoardStateEditor) PlaceSnake(id string, body []rules.Point, health int) {
	for index, snake := range editor.Snakes {
		if snake.ID == id {
			editor.Snakes[index].Body = body
			editor.Snakes[index].Health = health
			return
		}
	}

	editor.Snakes = append(editor.Snakes, rules.Snake{
		ID:     id,
		Health: health,
		Body:   body,
	})
}
