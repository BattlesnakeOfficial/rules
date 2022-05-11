package rules

import (
	"strconv"
)

type Ruleset interface {
	Name() string
	ModifyInitialBoardState(initialState *BoardState) (*BoardState, error)
	CreateNextBoardState(prevState *BoardState, moves []SnakeMove) (*BoardState, error)
	IsGameOver(state *BoardState) (bool, error)
	// Settings provides the game settings that are relevant to the ruleset.
	Settings() Settings
}

type SnakeMove struct {
	ID   string
	Move string
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

	rand Rand
}

func (settings Settings) Rand() Rand {
	// Default to global random number generator if none is set.
	if settings.rand == nil {
		return GlobalRand
	}
	return settings.rand
}

func (settings Settings) WithRand(rand Rand) Settings {
	settings.rand = rand
	return settings
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

type rulesetBuilder struct {
	params map[string]string // game customisation parameters
	seed   int64             // used for random events in games
	rand   Rand              // used for random number generation
	squads map[string]string // Snake ID -> Squad Name
}

// NewRulesetBuilder returns an instance of a builder for the Ruleset types.
func NewRulesetBuilder() *rulesetBuilder {
	return &rulesetBuilder{
		params: map[string]string{},
		squads: map[string]string{},
	}
}

// WithParams accepts a map of game parameters for customizing games.
//
// Parameters are copied. If called multiple times, parameters are merged such that:
//   - existing keys in both maps get overwritten by the new ones
//   - existing keys not present in the new map will be retained
//   - non-existing keys only in the new map will be added
//
// Unrecognised parameters will be ignored and default values will be used.
// Invalid parameters (i.e. a non-numerical value where one is expected), will be ignored
// and default values will be used.
func (rb *rulesetBuilder) WithParams(params map[string]string) *rulesetBuilder {
	for k, v := range params {
		rb.params[k] = v
	}
	return rb
}

// Deprecated: WithSeed sets the seed used for randomisation by certain game modes.
func (rb *rulesetBuilder) WithSeed(seed int64) *rulesetBuilder {
	rb.seed = seed
	return rb
}

func (rb *rulesetBuilder) WithRand(rand Rand) *rulesetBuilder {
	rb.rand = rand
	return rb
}

// AddSnakeToSquad adds the specified snake (by ID) to a squad with the given name.
// This configuration may be ignored by game modes if they do not support squads.
func (rb *rulesetBuilder) AddSnakeToSquad(snakeID, squadName string) *rulesetBuilder {
	rb.squads[snakeID] = squadName
	return rb
}

// Ruleset constructs a customised ruleset using the parameters passed to the builder.
func (rb rulesetBuilder) Ruleset() PipelineRuleset {
	standardRuleset := &StandardRuleset{
		FoodSpawnChance:     paramsInt32(rb.params, ParamFoodSpawnChance, 0),
		MinimumFood:         paramsInt32(rb.params, ParamMinimumFood, 0),
		HazardDamagePerTurn: paramsInt32(rb.params, ParamHazardDamagePerTurn, 0),
		HazardMap:           rb.params[ParamHazardMap],
		HazardMapAuthor:     rb.params[ParamHazardMapAuthor],
	}

	name, ok := rb.params[ParamGameType]
	if !ok {
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
			ShrinkEveryNTurns: paramsInt32(rb.params, ParamShrinkEveryNTurns, 0),
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
		return &SquadRuleset{
			StandardRuleset:     *standardRuleset,
			SquadMap:            rb.squadMap(),
			AllowBodyCollisions: paramsBool(rb.params, ParamAllowBodyCollisions, false),
			SharedElimination:   paramsBool(rb.params, ParamSharedElimination, false),
			SharedHealth:        paramsBool(rb.params, ParamSharedHealth, false),
			SharedLength:        paramsBool(rb.params, ParamSharedLength, false),
		}
	}
	return standardRuleset
}

func (rb rulesetBuilder) squadMap() map[string]string {
	squadMap := map[string]string{}
	for id, squad := range rb.squads {
		squadMap[id] = squad
	}
	return squadMap
}

// PipelineRuleset provides an implementation of the Ruleset using a pipeline with a name.
// It is intended to facilitate transitioning away from legacy Ruleset implementations to Pipeline
// implementations.
func (rb rulesetBuilder) PipelineRuleset(name string, p Pipeline) PipelineRuleset {
	return &pipelineRuleset{
		name:     name,
		pipeline: p,
		settings: Settings{
			FoodSpawnChance:     paramsInt32(rb.params, ParamFoodSpawnChance, 0),
			MinimumFood:         paramsInt32(rb.params, ParamMinimumFood, 0),
			HazardDamagePerTurn: paramsInt32(rb.params, ParamHazardDamagePerTurn, 0),
			HazardMap:           rb.params[ParamHazardMap],
			HazardMapAuthor:     rb.params[ParamHazardMapAuthor],
			RoyaleSettings: RoyaleSettings{
				seed:              rb.seed,
				ShrinkEveryNTurns: paramsInt32(rb.params, ParamShrinkEveryNTurns, 0),
			},
			SquadSettings: SquadSettings{
				squadMap:            rb.squadMap(),
				AllowBodyCollisions: paramsBool(rb.params, ParamAllowBodyCollisions, false),
				SharedElimination:   paramsBool(rb.params, ParamSharedElimination, false),
				SharedHealth:        paramsBool(rb.params, ParamSharedHealth, false),
				SharedLength:        paramsBool(rb.params, ParamSharedLength, false),
			},
			rand: rb.rand,
		},
	}
}

// paramsBool returns the boolean value for the specified parameter.
// If the parameter doesn't exist, the default value will be returned.
// If the parameter does exist, but is not "true", false will be returned.
func paramsBool(params map[string]string, paramName string, defaultValue bool) bool {
	if val, ok := params[paramName]; ok {
		return val == "true"
	}
	return defaultValue
}

// paramsInt32 returns the int32 value for the specified parameter.
// If the parameter doesn't exist, the default value will be returned.
// If the parameter does exist, but is not a valid int, the default value will be returned.
func paramsInt32(params map[string]string, paramName string, defaultValue int32) int32 {
	if val, ok := params[paramName]; ok {
		i, err := strconv.Atoi(val)
		if err == nil {
			return int32(i)
		}
	}
	return defaultValue
}

// PipelineRuleset groups the Pipeline and Ruleset methods.
// It is intended to facilitate a transition from Ruleset legacy code to Pipeline code.
type PipelineRuleset interface {
	Ruleset
	Pipeline
}

type pipelineRuleset struct {
	pipeline Pipeline
	name     string
	settings Settings
}

// impl Ruleset
func (r pipelineRuleset) Settings() Settings {
	return r.settings
}

// impl Ruleset
func (r pipelineRuleset) Name() string { return r.name }

// impl Ruleset
// IMPORTANT: this implementation of IsGameOver deviates from the previous Ruleset implementations
// in that it checks if the *NEXT* state results in game over, not the previous state.
// This is due to the design of pipelines / stage functions not having a distinction between
// checking for game over and producing a next state.
func (r *pipelineRuleset) IsGameOver(b *BoardState) (bool, error) {
	gameover, _, err := r.Execute(b, r.Settings(), nil) // checks if next state is game over
	return gameover, err
}

// impl Ruleset
func (r pipelineRuleset) ModifyInitialBoardState(initialState *BoardState) (*BoardState, error) {
	_, nextState, err := r.Execute(initialState, r.Settings(), nil)
	return nextState, err
}

// impl Pipeline
func (r pipelineRuleset) Execute(bs *BoardState, s Settings, sm []SnakeMove) (bool, *BoardState, error) {
	return r.pipeline.Execute(bs, s, sm)
}

// impl Ruleset
func (r pipelineRuleset) CreateNextBoardState(bs *BoardState, sm []SnakeMove) (*BoardState, error) {
	_, nextState, err := r.Execute(bs, r.Settings(), sm)
	return nextState, err
}

// impl Pipeline
func (r pipelineRuleset) Err() error {
	return r.pipeline.Err()
}
