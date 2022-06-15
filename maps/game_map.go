package maps

import (
	"github.com/BattlesnakeOfficial/rules"
	"github.com/blang/semver/v4"
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

// version is an internal type used to represent the version of a game map.
// A valid  version string conforms to semantic versioning (https://semver.org/).
// Examples: ["0.0.1", "1.2.34", "0.1.0"].
type version string

// IsValid returns whether the version is a valid semantic version string
func (v version) IsValid() bool {
	_, err := semver.Parse(string(v))
	return err == nil
}

type Metadata struct {
	Name        string
	Author      string
	Description string
	// Version is the current, semver version of the map.
	// Must be a valid semver (https://semver.org/) string.
	Version version
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
