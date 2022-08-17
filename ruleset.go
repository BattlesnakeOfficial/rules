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
	FoodSpawnChance     int            `json:"foodSpawnChance"`
	MinimumFood         int            `json:"minimumFood"`
	HazardDamagePerTurn int            `json:"hazardDamagePerTurn"`
	HazardMap           string         `json:"hazardMap"`
	HazardMapAuthor     string         `json:"hazardMapAuthor"`
	RoyaleSettings      RoyaleSettings `json:"royale"`
	SquadSettings       SquadSettings  `json:"squad"` // Deprecated, provided with default fields for API compatibility

	rand Rand
	seed int64
}

// Get a random number generator initialized based on the seed and current turn.
func (settings Settings) GetRand(turn int) Rand {
	// Allow overriding the random generator for testing
	if settings.rand != nil {
		return settings.rand
	}

	if settings.seed != 0 {
		return NewSeedRand(settings.seed + int64(turn))
	}

	// Default to global random number generator if neither seed or rand are set.
	return GlobalRand
}

func (settings Settings) WithRand(rand Rand) Settings {
	settings.rand = rand
	return settings
}

func (settings Settings) Seed() int64 {
	return settings.seed
}

func (settings Settings) WithSeed(seed int64) Settings {
	settings.seed = seed
	return settings
}

// RoyaleSettings contains settings that are specific to the "royale" game mode
type RoyaleSettings struct {
	ShrinkEveryNTurns int `json:"shrinkEveryNTurns"`
}

// SquadSettings contains settings that are specific to the "squad" game mode
type SquadSettings struct {
	AllowBodyCollisions bool `json:"allowBodyCollisions"`
	SharedElimination   bool `json:"sharedElimination"`
	SharedHealth        bool `json:"sharedHealth"`
	SharedLength        bool `json:"sharedLength"`
}

type rulesetBuilder struct {
	params map[string]string // game customisation parameters
	seed   int64             // used for random events in games
	rand   Rand              // used for random number generation
	solo   bool              // if true, only 1 alive snake is required to keep the game from ending
}

// NewRulesetBuilder returns an instance of a builder for the Ruleset types.
func NewRulesetBuilder() *rulesetBuilder {
	return &rulesetBuilder{
		params: map[string]string{},
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

// WithSeed sets the seed used for randomisation by certain game modes.
func (rb *rulesetBuilder) WithSeed(seed int64) *rulesetBuilder {
	rb.seed = seed
	return rb
}

// WithRandom overrides the random number generator with a specific instance
// instead of a Rand initialized from the seed.
func (rb *rulesetBuilder) WithRand(rand Rand) *rulesetBuilder {
	rb.rand = rand
	return rb
}

// WithSolo sets whether the ruleset is a solo game.
func (rb *rulesetBuilder) WithSolo(value bool) *rulesetBuilder {
	rb.solo = value
	return rb
}

// Ruleset constructs a customised ruleset using the parameters passed to the builder.
func (rb rulesetBuilder) Ruleset() PipelineRuleset {
	name, ok := rb.params[ParamGameType]
	if !ok {
		name = GameTypeStandard
	}

	var stages []string
	if rb.solo {
		stages = append(stages, StageGameOverSoloSnake)
	} else {
		stages = append(stages, StageGameOverStandard)
	}

	switch name {
	case GameTypeStandard:
		stages = append(stages, standardRulesetStages[1:]...)
	case GameTypeConstrictor:
		stages = append(stages, constrictorRulesetStages[1:]...)
	case GameTypeWrappedConstrictor:
		stages = append(stages, wrappedConstrictorRulesetStages[1:]...)
	case GameTypeRoyale:
		stages = append(stages, royaleRulesetStages[1:]...)
	case GameTypeSolo:
		stages = soloRulesetStages
	case GameTypeWrapped:
		stages = append(stages, wrappedRulesetStages[1:]...)
	default:
		stages = append(stages, standardRulesetStages[1:]...)
	}
	return rb.PipelineRuleset(name, NewPipeline(stages...))
}

// PipelineRuleset provides an implementation of the Ruleset using a pipeline with a name.
// It is intended to facilitate transitioning away from legacy Ruleset implementations to Pipeline
// implementations.
func (rb rulesetBuilder) PipelineRuleset(name string, p Pipeline) PipelineRuleset {
	return &pipelineRuleset{
		name:     name,
		pipeline: p,
		settings: Settings{
			FoodSpawnChance:     paramsInt(rb.params, ParamFoodSpawnChance, 0),
			MinimumFood:         paramsInt(rb.params, ParamMinimumFood, 0),
			HazardDamagePerTurn: paramsInt(rb.params, ParamHazardDamagePerTurn, 0),
			HazardMap:           rb.params[ParamHazardMap],
			HazardMapAuthor:     rb.params[ParamHazardMapAuthor],
			RoyaleSettings: RoyaleSettings{
				ShrinkEveryNTurns: paramsInt(rb.params, ParamShrinkEveryNTurns, 0),
			},
			rand: rb.rand,
			seed: rb.seed,
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

// paramsInt returns the int value for the specified parameter.
// If the parameter doesn't exist, the default value will be returned.
// If the parameter does exist, but is not a valid int, the default value will be returned.
func paramsInt(params map[string]string, paramName string, defaultValue int) int {
	if val, ok := params[paramName]; ok {
		i, err := strconv.Atoi(val)
		if err == nil {
			return i
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
