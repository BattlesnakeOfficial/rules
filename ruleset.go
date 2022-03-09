package rules

import (
	"math"

	"github.com/buger/jsonparser"
)

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

// Settings contains all settings relevant to a game.
// It is used by game logic to take a previous game state and produce a next game state.
type Settings struct {
	FoodSpawnChance     int32          `json:"foodSpawnChance"`
	MinimumFood         int32          `json:"minimumFood"`
	HazardDamagePerTurn int32          `json:"hazardDamagePerTurn"`
	HazardMap           string         `json:"hazardMap"`
	HazardMapAuthor     string         `json:"hazardMapAuthor"`
	RoyaleSettings      RoyaleSettings `json:"royale"`
	SquadSettings       SquadSettings  `json:"squad"`
}

// RoyaleSettings contains settings that are specific to the "royale" game mode
type RoyaleSettings struct {
	seed              int64
	ShrinkEveryNTurns int32 `json:"shrinkEveryNTurns"`
}

// SquadSettings contains settings that are specific to the "squad" game mode
type SquadSettings struct {
	squadMap            map[string]string
	AllowBodyCollisions bool `json:"allowBodyCollisions"`
	SharedElimination   bool `json:"sharedElimination"`
	SharedHealth        bool `json:"sharedHealth"`
	SharedLength        bool `json:"sharedLength"`
}

type StageSettings interface {
	GetJSON() SettingsJSON
}

// SettingsJSON contains settings for game rules in JSON format
type SettingsJSON []byte

// // SettingsJSON contains settings for game rules in JSON format
// type SettingsJSON struct {
// 	Settings []byte            // JSON encoded game settings
// 	seed     int64             // seed for generating random numbers
// 	squadMap map[string]string // mapping of snake ids to squad ids
// }

// GetInt32 returns the int32 at the specified path.
// Path format is "foo.bar[0].baz" == ["foo","bar", "[0]","baz"].
func (s SettingsJSON) GetInt32(keys ...string) int32 {
	v, err := jsonparser.GetInt(s, keys...)

	// errors default to zero value
	if err != nil {
		return 0
	}

	// overflows will default to zero value
	if v < math.MinInt32 || v > math.MaxInt32 {
		return 0
	}

	return int32(v)
}

// GetInt64 returns the int64 at the specified path.
// Path format is "foo.bar[0].baz" == ["foo","bar", "[0]","baz"].
func (s SettingsJSON) GetInt64(keys ...string) int64 {
	v, err := jsonparser.GetInt(s, keys...)

	// errors default to zero value
	if err != nil {
		return 0
	}
	return v
}

// GetBool returns the bool at the specified path.
// Path format is "foo.bar[0].baz" == ["foo","bar", "[0]","baz"].
func (s SettingsJSON) GetBool(keys ...string) bool {
	v, err := jsonparser.GetBoolean(s, keys...)

	// errors default to zero value
	if err != nil {
		return false
	}
	return v
}

// GetString gets the string at the specified path.
// Path format is "foo.bar[0].baz" == ["foo","bar", "[0]","baz"].
func (s SettingsJSON) GetString(keys ...string) string {
	v, err := jsonparser.GetString(s, keys...)

	// errors default to zero value
	if err != nil {
		return ""
	}
	return v
}

// StageFunc represents a single stage of an ordered pipeline and applies custom logic to the board state each turn.
// It is expected to modify the boardState directly.
type StageFunc func(*BoardState, SettingsJSON, []SnakeMove) (bool, error)
