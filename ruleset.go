package rules

const (
	MoveUp    = "up"
	MoveDown  = "down"
	MoveRight = "right"
	MoveLeft  = "left"
)

type Point struct {
	X int32
	Y int32
}

type Snake struct {
	ID              string
	Body            []Point
	Health          int32
	EliminatedCause string
	EliminatedBy    string
}

type BoardState struct {
	Height int32
	Width  int32
	Food   []Point
	Snakes []Snake
}

type SnakeMove struct {
	ID   string
	Move string
}

type Ruleset interface {
	CreateInitialBoardState(width int32, height int32, snakeIDs []string) (*BoardState, error)
	CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error)
	IsGameOver(state *BoardState) (bool, error)
}

type RulesetError string

func (e RulesetError) Error() string { return string(e) }
