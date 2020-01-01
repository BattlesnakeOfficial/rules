package rulesets

// NOTE: IMMUTABLE THINGS HERE //

const MOVE_UP = "up"
const MOVE_DOWN = "down"
const MOVE_RIGHT = "right"
const MOVE_LEFT = "left"

type Game struct {
	Height int32
	Width  int32
}

type SnakeMove struct {
	Snake *Snake
	Move  string
}

// NOTE: MUTABLE THINGS HERE //

type Point struct {
	X int32
	Y int32
}

type Snake struct {
	ID              string
	Body            []*Point
	Health          int32
	EliminatedCause string
}

type GameState struct {
	Food   []*Point
	Snakes []*Snake
}

// RULESET API //

type Ruleset interface {
	ResolveMoves(*Game, *GameState, []*SnakeMove) (*GameState, error)
}
