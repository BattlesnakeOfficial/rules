package maps

import (
	"fmt"
	"strings"

	"github.com/BattlesnakeOfficial/rules"
)

const (
	TAG_EXPERIMENTAL     = "experimental"     // experimental map, only available via CLI
	TAG_SNAKE_PLACEMENT  = "snake-placement"  // map overrides default snake placement
	TAG_HAZARD_PLACEMENT = "hazard-placement" // map places hazards
	TAG_FOOD_PLACEMENT   = "food-placement"   // map overrides or adds to default food placement
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

type Metadata struct {
	Name        string
	Author      string
	Description string
	// Version is the current version of the game map.
	// Each time a map is changed, the version number should be incremented by 1.
	Version int
	// MinPlayers is the minimum number of players that the map supports.
	MinPlayers int
	// MaxPlayers is the maximum number of players that the map supports.
	MaxPlayers int
	// BoardSizes is a list of supported board sizes. Board sizes can fall into one of 3 categories:
	//   1. one fixed size (i.e. [11x11])
	//   2. multiple, fixed sizes (i.e. [11x11, 19x19, 25x25])
	//   3. "unlimited" sizes (the board is not fixed and can scale to any reasonable size)
	BoardSizes sizes
	// Tags is a list of strings use to categorize the map.
	Tags []string
}

func (meta Metadata) Validate(boardState *rules.BoardState) error {
	if !meta.BoardSizes.IsAllowable(boardState.Width, boardState.Height) {
		var sizesStrings []string
		for _, size := range meta.BoardSizes {
			sizesStrings = append(sizesStrings, fmt.Sprintf("%dx%d", size.Width, size.Height))
		}

		return rules.RulesetError("This map can only be played on these board sizes: " + strings.Join(sizesStrings, ", "))
	}

	if meta.MinPlayers != 0 && len(boardState.Snakes) < int(meta.MinPlayers) {
		return rules.RulesetError(fmt.Sprintf("This map can only be played with %d-%d players", meta.MinPlayers, meta.MaxPlayers))
	}

	if meta.MaxPlayers != 0 && len(boardState.Snakes) > int(meta.MaxPlayers) {
		return rules.RulesetError(fmt.Sprintf("This map can only be played with %d-%d players", meta.MinPlayers, meta.MaxPlayers))
	}

	return nil
}

// Dimensions describes the size of a Battlesnake board.
type Dimensions struct {
	// Width is the width, in number of board squares, of the board.
	// The value 0 has a special meaning to mean unlimited.
	Width int
	// Height is the height, in number of board squares, of the board.
	// The value 0 has a special meaning to mean unlimited.
	Height int
}

// sizes is a list of board sizes that a map supports.
type sizes []Dimensions

// IsUnlimited reports whether the supported sizes are unlimited.
// Note that even for unlimited sizes, there will be an upper bound that can actually be run and visualised.
func (d sizes) IsUnlimited() bool {
	return len(d) == 1 && d[0].Width == 0
}

func (d sizes) IsAllowable(Width int, Height int) bool {
	if d.IsUnlimited() {
		return true
	}

	for _, size := range d {
		if size.Width == Width && size.Height == Height {
			return true
		}
	}

	return false
}

// AnySize creates sizes for a board that has no fixed sizes (supports unlimited sizes).
func AnySize() sizes {
	return sizes{Dimensions{Width: 0, Height: 0}}
}

// OddSizes generates square (width = height) board sizes with an odd number of positions
// in the vertical and horizontal directions.
// Examples:
//  - OddSizes(11,21) produces [(11,11), (13,13), (15,15), (17,17), (19,19), (21,21)]
func OddSizes(min, max int) sizes {
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

// Editor is used by GameMap implementations to modify the board state.
type Editor interface {
	// Clears all food from the board.
	ClearFood()

	// Adds a food to the board. Does not check for duplicates.
	AddFood(rules.Point)

	// Removes all food from a specific tile on the board.
	RemoveFood(rules.Point)

	// Get the locations of food currently on the board.
	// Note: the return value is a copy and modifying it won't affect the board.
	Food() []rules.Point

	// Clears all hazards from the board.
	ClearHazards()

	// Adds a hazard to the board. Does not check for duplicates.
	AddHazard(rules.Point)

	// Removes all hazards from a specific tile on the board.
	RemoveHazard(rules.Point)

	// Get the locations of hazards currently on the board.
	// Note: the return value is a copy and modifying it won't affect the board.
	Hazards() []rules.Point

	// Updates the body and health of a snake.
	PlaceSnake(id string, body []rules.Point, health int)

	// Get the bodies of all non-eliminated snakes currently on the board, keyed by Snake ID
	// Note: the body values in the return value are a copy and modifying them won't affect the board.
	SnakeBodies() map[string][]rules.Point

	// Given a list of Snakes and a list of head coordinates, randomly place
	// the snakes on those coordinates, or return an error if placement of all
	// Snakes is impossible.
	PlaceSnakesRandomlyAtPositions(rand rules.Rand, snakes []rules.Snake, heads []rules.Point, bodyLength int) error

	// Returns true if the provided point on the board is occupied by a snake body, food, and/or hazard.
	IsOccupied(point rules.Point, snakes, hazards, food bool) bool

	// Get a set of all points on the board the are occupied by snake bodies, food, and/or hazards.
	// The value for each point will be set to true in the return value if that point is occupied by one of the selected objects.
	OccupiedPoints(snakes, hazards, food bool) map[rules.Point]bool

	// Given a list of points, return only those that are unoccupied by snake bodies, food, and/or hazards.
	FilterUnoccupiedPoints(targets []rules.Point, snakes, hazards, food bool) []rules.Point

	// Shuffle the provided slice of points randomly using the provided rules.Rand
	ShufflePoints(rules.Rand, []rules.Point)
}

// An Editor backed by a BoardState.
type BoardStateEditor struct {
	boardState *rules.BoardState
}

func NewBoardStateEditor(boardState *rules.BoardState) *BoardStateEditor {
	return &BoardStateEditor{
		boardState: boardState,
	}
}

func (editor *BoardStateEditor) ClearFood() {
	editor.boardState.Food = []rules.Point{}
}

func (editor *BoardStateEditor) AddFood(p rules.Point) {
	editor.boardState.Food = append(editor.boardState.Food, rules.Point{X: p.X, Y: p.Y})
}

func (editor *BoardStateEditor) RemoveFood(p rules.Point) {
	for index, food := range editor.boardState.Food {
		if food.X == p.X && food.Y == p.Y {
			editor.boardState.Food[index] = editor.boardState.Food[len(editor.boardState.Food)-1]
			editor.boardState.Food = editor.boardState.Food[:len(editor.boardState.Food)-1]
		}
	}
}

// Get the locations of food currently on the board.
// Note: the return value is read-only.
func (editor *BoardStateEditor) Food() []rules.Point {
	return append([]rules.Point(nil), editor.boardState.Food...)
}

func (editor *BoardStateEditor) ClearHazards() {
	editor.boardState.Hazards = []rules.Point{}
}

func (editor *BoardStateEditor) AddHazard(p rules.Point) {
	editor.boardState.Hazards = append(editor.boardState.Hazards, rules.Point{X: p.X, Y: p.Y})
}

func (editor *BoardStateEditor) RemoveHazard(p rules.Point) {
	for index, food := range editor.boardState.Hazards {
		if food.X == p.X && food.Y == p.Y {
			editor.boardState.Hazards[index] = editor.boardState.Hazards[len(editor.boardState.Hazards)-1]
			editor.boardState.Hazards = editor.boardState.Hazards[:len(editor.boardState.Hazards)-1]
		}
	}
}

// Get the locations of hazards currently on the board.
// Note: the return value is read-only.
func (editor *BoardStateEditor) Hazards() []rules.Point {
	return append([]rules.Point(nil), editor.boardState.Hazards...)
}

func (editor *BoardStateEditor) PlaceSnake(id string, body []rules.Point, health int) {
	for index, snake := range editor.boardState.Snakes {
		if snake.ID == id {
			editor.boardState.Snakes[index].Body = body
			editor.boardState.Snakes[index].Health = health
			return
		}
	}

	editor.boardState.Snakes = append(editor.boardState.Snakes, rules.Snake{
		ID:     id,
		Health: health,
		Body:   body,
	})
}

// Get the bodies of all non-eliminated snakes currently on the board.
// Note: the return value is read-only.
func (editor *BoardStateEditor) SnakeBodies() map[string][]rules.Point {
	result := make(map[string][]rules.Point, len(editor.boardState.Snakes))

	for _, snake := range editor.boardState.Snakes {
		result[snake.ID] = append([]rules.Point(nil), snake.Body...)
	}

	return result
}

// Given a list of Snakes and a list of head coordinates, randomly place
// the snakes on those coordinates, or return an error if placement of all
// Snakes is impossible.
func (editor *BoardStateEditor) PlaceSnakesRandomlyAtPositions(rand rules.Rand, snakes []rules.Snake, heads []rules.Point, bodyLength int) error {
	if len(snakes) > len(heads) {
		return rules.ErrorTooManySnakes
	}

	// Shuffle starting points
	editor.ShufflePoints(rand, heads)

	// Assign starting points to snakes in order
	for index, snake := range snakes {
		head := heads[index]
		body := make([]rules.Point, bodyLength)
		for i := 0; i < bodyLength; i++ {
			body[i] = head
		}
		editor.PlaceSnake(snake.ID, body, rules.SnakeMaxHealth)
	}

	return nil
}

// Returns true if the provided point on the board is occupied by a snake body, food, and/or hazard.
func (editor *BoardStateEditor) IsOccupied(point rules.Point, snakes, hazards, food bool) bool {
	if food {
		for _, food := range editor.boardState.Food {
			if food == point {
				return true
			}
		}
	}
	if hazards {
		for _, hazard := range editor.boardState.Hazards {
			if hazard == point {
				return true
			}
		}
	}
	if snakes {
		for _, snake := range editor.boardState.Snakes {
			for _, body := range snake.Body {
				if body == point {
					return true
				}
			}
		}
	}
	return false
}

// Get a set of all points on the board the are occupied by snake bodies, food, and/or hazards.
// The value for each point will be set to true in the return value if that point is occupied by one of the selected objects.
func (editor *BoardStateEditor) OccupiedPoints(snakes, hazards, food bool) map[rules.Point]bool {
	boardState := editor.boardState
	result := make(map[rules.Point]bool, len(boardState.Food)+len(boardState.Hazards)+len(boardState.Snakes)*3)

	if food {
		for _, food := range editor.boardState.Food {
			result[food] = true
		}
	}
	if hazards {
		for _, hazard := range editor.boardState.Hazards {
			result[hazard] = true
		}
	}
	if snakes {
		for _, snake := range editor.boardState.Snakes {
			for _, body := range snake.Body {
				result[body] = true
			}
		}
	}

	return result
}

// Given a list of points, return only those that are unoccupied by snake bodies, food, and/or hazards.
func (editor *BoardStateEditor) FilterUnoccupiedPoints(targets []rules.Point, snakes, hazards, food bool) []rules.Point {
	result := make([]rules.Point, 0, len(targets))

targetLoop:
	for _, point := range targets {
		if food {
			for _, food := range editor.boardState.Food {
				if food == point {
					continue targetLoop
				}
			}
		}
		if hazards {
			for _, hazard := range editor.boardState.Hazards {
				if hazard == point {
					continue targetLoop
				}
			}
		}
		if snakes {
			for _, snake := range editor.boardState.Snakes {
				for _, body := range snake.Body {
					if body == point {
						continue targetLoop
					}
				}
			}
		}

		result = append(result, point)
	}

	return result
}

func (editor *BoardStateEditor) ShufflePoints(rand rules.Rand, points []rules.Point) {
	rand.Shuffle(len(points), func(i int, j int) {
		points[i], points[j] = points[j], points[i]
	})
}
