package rules

import (
	"fmt"
	"strconv"
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
	EliminatedBySquad               = "squad-eliminated"

	// TODO - Error consts
	ErrorTooManySnakes   = RulesetError("too many snakes for fixed start positions")
	ErrorNoRoomForSnake  = RulesetError("not enough space to place snake")
	ErrorNoRoomForFood   = RulesetError("not enough space to place food")
	ErrorNoMoveFound     = RulesetError("move not provided for snake")
	ErrorZeroLengthSnake = RulesetError("snake is length zero")

	// Ruleset names
	GameTypeConstrictor = "constrictor"
	GameTypeRoyale      = "royale"
	GameTypeSolo        = "solo"
	GameTypeSquad       = "squad"
	GameTypeStanadard   = "standard"
	GameTypeWrapped     = "wrapped"

	// Game creation parameter names
	ParamGameType            = "name"
	ParamFoodSpawnChance     = "foodSpawnChance"
	ParamMinimumFood         = "minimumFood"
	ParamHazardDamagePerTurn = "damagePerTurn"
	ParamShrinkEveryNTurns   = "shrinkEveryNTurns"
	ParamAllowBodyCollisions = "allowBodyCollisions"
	ParamSharedElimination   = "sharedElimination"
	ParamSharedHealth        = "sharedHealth"
	ParamSharedLength        = "sharedLength"
)

type builder struct {
	params map[string]string
	seed   int64
	squads map[string]string
}

func NewBuilder() *builder {
	return &builder{
		params: map[string]string{},
		squads: map[string]string{},
	}
}

func (rb *builder) WithParams(params map[string]string) *builder {
	for k, v := range params {
		rb.params[k] = v
	}
	fmt.Printf("wp %v\n", rb)
	return rb
}

func (rb *builder) WithSeed(seed int64) *builder {
	rb.seed = seed
	fmt.Printf("ws %v\n", rb)
	return rb
}

func (rb *builder) AddSnakeToSquad(snakeID, squadName string) *builder {
	rb.squads[snakeID] = squadName
	fmt.Printf("asts %v\n", rb)
	return rb
}

// Build constructs a ruleset from the parameters passed when creating a
// new game, and returns a ruleset customised by those parameters.
func (rb builder) Ruleset() Ruleset {
	standardRuleset := &StandardRuleset{
		FoodSpawnChance:     optionFromRulesetInt(rb.params, ParamFoodSpawnChance, 0),
		MinimumFood:         optionFromRulesetInt(rb.params, ParamMinimumFood, 0),
		HazardDamagePerTurn: optionFromRulesetInt(rb.params, ParamHazardDamagePerTurn, 0),
	}

	name, ok := rb.params[ParamGameType]
	if !ok {
		fmt.Printf("%v\n", rb.params)
		return standardRuleset
	}

	switch name {
	case GameTypeConstrictor:
		return &ConstrictorRuleset{
			StandardRuleset: *standardRuleset,
		}
	case GameTypeRoyale:
		return &RoyaleRuleset{
			StandardRuleset:   *standardRuleset,
			Seed:              rb.seed,
			ShrinkEveryNTurns: optionFromRulesetInt(rb.params, ParamShrinkEveryNTurns, 0),
		}
	case GameTypeSolo:
		return &SoloRuleset{
			StandardRuleset: *standardRuleset,
		}
	case GameTypeWrapped:
		return &WrappedRuleset{
			StandardRuleset: *standardRuleset,
		}
	case GameTypeSquad:
		squadMap := map[string]string{}
		for id, squad := range rb.squads {
			squadMap[id] = squad
		}
		return &SquadRuleset{
			StandardRuleset:     *standardRuleset,
			SquadMap:            squadMap,
			AllowBodyCollisions: optionFromRulesetBool(rb.params, ParamAllowBodyCollisions, true),
			SharedElimination:   optionFromRulesetBool(rb.params, ParamSharedElimination, true),
			SharedHealth:        optionFromRulesetBool(rb.params, ParamSharedHealth, true),
			SharedLength:        optionFromRulesetBool(rb.params, ParamSharedLength, true),
		}
	}
	return standardRuleset
}

func optionFromRulesetBool(ruleset map[string]string, optionName string, defaultValue bool) bool {
	if val, ok := ruleset[optionName]; ok {
		return val == "true"
	}
	return defaultValue
}

func optionFromRulesetInt(ruleset map[string]string, optionName string, defaultValue int32) int32 {
	if val, ok := ruleset[optionName]; ok {
		i, err := strconv.Atoi(val)
		if err == nil {
			return int32(i)
		}
	}
	return defaultValue
}

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
	// Settings provides the game settings that are relevant to the ruleset.
	Settings() Settings
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

// StageFunc represents a single stage of an ordered pipeline and applies custom logic to the board state each turn.
// It is expected to modify the boardState directly.
type StageFunc func(*BoardState, Settings, []SnakeMove) (bool, error)
