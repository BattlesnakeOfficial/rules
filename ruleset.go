package rules

type RulesetError string

func (err RulesetError) Error() string { return string(err) }

const (
	MoveUp    = "up"
	MoveDown  = "down"
	MoveRight = "right"
	MoveLeft  = "left"

	BoardSizeSmall  = 7
	BoardSizeMedium = 11
	BoardSizeLarge  = 19

	SnakeMaxHealth = 100
	SnakeStartSize = 3

	// bvanvugt - TODO: Just return formatted strings instead of codes?
	NotEliminated                   = ""
	EliminatedByCollision           = "snake-collision"
	EliminatedBySelfCollision       = "snake-self-collision"
	EliminatedByOutOfHealth         = "out-of-health"
	EliminatedByHeadToHeadCollision = "head-collision"
	EliminatedByOutOfBounds         = "wall-collision"

	// TODO - Error consts
	ErrorTooManySnakes   = RulesetError("too many snakes for fixed start positions")
	ErrorNoRoomForSnake  = RulesetError("not enough space to place snake")
	ErrorNoRoomForFood   = RulesetError("not enough space to place food")
	ErrorNoMoveFound     = RulesetError("move not provided for snake")
	ErrorZeroLengthSnake = RulesetError("snake is length zero")
)

type Point struct {
	X int32
	Y int32
}

type Snake struct {
	ID               string
	Body             []Point
	Health           int32
	EliminatedCause  string
	EliminatedOnTurn int32
	EliminatedBy     string
}

type SnakeMove struct {
	ID   string
	Move string
}

type Ruleset interface {
	Name() string
	ModifyInitialBoardState(initialState *BoardState) (*BoardState, error)
	CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error)
	IsGameOver(state *BoardState) (bool, error)
}

type RulesetSettings struct {
	FoodSpawnChance     int32          `json:"foodSpawnChance"`
	MinimumFood         int32          `json:"minimumFood"`
	HazardDamagePerTurn int32          `json:"hazardDamagePerTurn"`
	RoyaleSettings      RoyaleSettings `json:"royale"`
	SquadSettings       SquadSettings  `json:"squad"`
}

type RoyaleSettings struct {
	seed              int64
	ShrinkEveryNTurns int32 `json:"shrinkEveryNTurns"`
}

type SquadSettings struct {
	AllowBodyCollisions bool `json:"allowBodyCollisions"`
	SharedElimination   bool `json:"sharedElimination"`
	SharedHealth        bool `json:"sharedHealth"`
	SharedLength        bool `json:"sharedLength"`
}

// Represents a single stage of an ordered pipeline and applies custom logic to the board state each turn.
// modifyBoardState is expected to modify the boardState directly, not copy it.
type Stage interface {
	ModifyBoardState(boardState *BoardState, settings RulesetSettings, snakeIDs []string, moves []SnakeMove) (gameOver bool, err error)
}

// Allows converting a plain function to a RulesStage
type StageFunc func(*BoardState, RulesetSettings, []string, []SnakeMove) (bool, error)

func (f StageFunc) ModifyBoardState(boardState *BoardState, settings RulesetSettings, snakeIDs []string, moves []SnakeMove) (bool, error) {
	return f(boardState, settings, snakeIDs, moves)
}
