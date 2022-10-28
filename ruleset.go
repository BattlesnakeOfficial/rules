package rules

type Ruleset interface {
	// Returns the name of the ruleset, if applicable.
	Name() string

	// Returns the settings used by the ruleset.
	Settings() Settings

	// Processes the next turn of the ruleset, returning whether the game has ended, the next BoardState, or an error.
	// For turn zero (initialization), moves will be left empty.
	Execute(prevState *BoardState, moves []SnakeMove) (gameOver bool, nextState *BoardState, err error)
}

type SnakeMove struct {
	ID   string
	Move string
}

type rulesetBuilder struct {
	params   map[string]string // game customisation parameters
	seed     int64             // used for random events in games
	rand     Rand              // used for random number generation
	solo     bool              // if true, only 1 alive snake is required to keep the game from ending
	settings *Settings         // used to set settings directly instead of via string params
}

// NewRulesetBuilder returns an instance of a builder for the Ruleset types.
func NewRulesetBuilder() *rulesetBuilder {
	return &rulesetBuilder{
		params: map[string]string{},
	}
}

// WithParams accepts a map of string parameters for customizing games.
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

// WithSettings sets the settings object for the ruleset directly.
func (rb *rulesetBuilder) WithSettings(settings Settings) *rulesetBuilder {
	rb.settings = &settings
	return rb
}

// NamedRuleset constructs a known ruleset by using name to look up a standard pipeline.
func (rb rulesetBuilder) NamedRuleset(name string) Ruleset {
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
		name = GameTypeStandard
		stages = append(stages, standardRulesetStages[1:]...)
	}
	return rb.PipelineRuleset(name, NewPipeline(stages...))
}

// PipelineRuleset constructs a ruleset with the given name and pipeline using the parameters passed to the builder.
// This can be used to create custom rulesets.
func (rb rulesetBuilder) PipelineRuleset(name string, p Pipeline) Ruleset {
	var settings Settings
	if rb.settings != nil {
		settings = *rb.settings
	} else {
		settings = NewSettings(rb.params).WithRand(rb.rand).WithSeed(rb.seed)
	}
	return &pipelineRuleset{
		name:     name,
		pipeline: p,
		settings: settings,
	}
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
func (r pipelineRuleset) Execute(bs *BoardState, sm []SnakeMove) (bool, *BoardState, error) {
	return r.pipeline.Execute(bs, r.Settings(), sm)
}

func (r pipelineRuleset) Err() error {
	return r.pipeline.Err()
}
