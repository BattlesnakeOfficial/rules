package rules

import "fmt"

const (
	StageSpawnFoodStandard    = "spawn_food.standard"
	StageGameOverStandard     = "game_over.standard"
	StageStarvationStandard   = "starvation.standard"
	StageFeedSnakesStandard   = "feed_snakes.standard"
	StageMovementStandard     = "movement.standard"
	StageHazardDamageStandard = "hazard_damage.standard"
	StageEliminationStandard  = "elimination.standard"

	StageGameOverSoloSnake           = "game_over.solo_snake"
	StageSpawnFoodNoFood             = "spawn_food.no_food"
	StageSpawnHazardsShrinkMap       = "spawn_hazards.shrink_map"
	StageModifySnakesAlwaysGrow      = "modify_snakes.always_grow"
	StageMovementWrapBoundaries      = "movement.wrap_boundaries"
	StageModifySnakesShareAttributes = "modify_snakes.share_attributes"
)

// globalRegistry is a global, default mapping of stage names to stage functions.
// It can be extended by plugins through the use of registration functions.
// Plugins that wish to extend the available game stages should call RegisterPipelineStageError
// to add additional stages.
var globalRegistry = StageRegistry{
	StageSpawnFoodNoFood:        RemoveFoodConstrictor,
	StageSpawnFoodStandard:      SpawnFoodStandard,
	StageGameOverSoloSnake:      GameOverSolo,
	StageGameOverStandard:       GameOverStandard,
	StageHazardDamageStandard:   DamageHazardsStandard,
	StageSpawnHazardsShrinkMap:  PopulateHazardsRoyale,
	StageStarvationStandard:     ReduceSnakeHealthStandard,
	StageFeedSnakesStandard:     FeedSnakesStandard,
	StageEliminationStandard:    EliminateSnakesStandard,
	StageModifySnakesAlwaysGrow: GrowSnakesConstrictor,
	StageMovementStandard:       MoveSnakesStandard,
	StageMovementWrapBoundaries: MoveSnakesWrapped,
}

// Pipeline is an ordered sequences of game stages which are executed to produce the
// next game state.
//
// If a stage produces an error or an ended game state, the pipeline is halted at that stage.
type Pipeline interface {
	// Execute runs the pipeline stages and produces a next game state.
	//
	// If any stage produces an error or an ended game state, the pipeline
	// immediately stops at that stage.
	//
	// Errors should be checked and the other results ignored if error is non-nil.
	//
	// If the pipeline is already in an error state (this can be checked by calling Err()),
	// this error will be immediately returned and the pipeline will not run.
	//
	// After the pipeline runs, the results will be the result of the last stage that was executed.
	Execute(*BoardState, Settings, []SnakeMove) (bool, *BoardState, error)

	// Err provides a way to check for errors before/without calling Execute.
	// Err returns an error if the Pipeline is in an error state.
	// If this error is not nil, this error will also be returned from Execute, so it is
	// optional to call Err.
	// The idea is to reduce error-checking verbosity for the majority of cases where a
	// Pipeline is immediately executed after construction (i.e. NewPipeline(...).Execute(...)).
	Err() error
}

// StageFunc represents a single stage of an ordered pipeline and applies custom logic to the board state each turn.
// It is expected to modify the boardState directly.
// The return values are a boolean (to indicate whether the game has ended as a result of the stage)
// and an error if any errors occurred during the stage.
//
// Errors should be treated as meaning the stage failed and the board state is now invalid.
type StageFunc func(*BoardState, Settings, []SnakeMove) (bool, error)

// IsInitialization checks whether the current state means the game is initialising (turn zero).
// Useful for StageFuncs that need to apply different behaviour on initialisation.
func IsInitialization(b *BoardState, settings Settings, moves []SnakeMove) bool {
	// We can safely assume that the game state is in the initialisation phase when
	// the turn hasn't advanced and the moves are empty
	return b.Turn <= 0 && len(moves) == 0
}

// StageRegistry is a mapping of stage names to stage functions
type StageRegistry map[string]StageFunc

// RegisterPipelineStage adds a stage to the registry.
// If a stage has already been mapped it will be overwritten by the newly
// registered function.
func (sr StageRegistry) RegisterPipelineStage(s string, fn StageFunc) {
	sr[s] = fn
}

// RegisterPipelineStageError adds a stage to the registry.
// If a stage has already been mapped an error will be returned.
func (sr StageRegistry) RegisterPipelineStageError(s string, fn StageFunc) error {
	if _, ok := sr[s]; ok {
		return RulesetError(fmt.Sprintf("stage '%s' has already been registered", s))
	}

	sr.RegisterPipelineStage(s, fn)
	return nil
}

// RegisterPipelineStage adds a stage to the global stage registry.
// It will panic if the a stage has already been registered with the same name.
func RegisterPipelineStage(s string, fn StageFunc) {
	err := globalRegistry.RegisterPipelineStageError(s, fn)
	if err != nil {
		panic(err)
	}
}

// pipeline is an implementation of Pipeline
type pipeline struct {
	// stages is a list of stages that should be executed from slice start to end
	stages []StageFunc
	// if the pipeline has an error
	err error
}

// NewPipeline constructs an instance of Pipeline using the global registry.
// It is a convenience wrapper for NewPipelineFromRegistry when you want
// to use the default, global registry.
func NewPipeline(stageNames ...string) Pipeline {
	return NewPipelineFromRegistry(globalRegistry, stageNames...)
}

// NewPipelineFromRegistry constructs an instance of Pipeline, using the specified registry
// of pipeline stage functions.
//
// The order of execution for the pipeline stages will correspond to the order that
// the stage names are provided.
//
// Example:
//
//	NewPipelineFromRegistry(r, s, "stage1", "stage2")
//
// ... will result in stage "stage1" running first, then stage "stage2" running after.
//
// An error will be returned if an unregistered stage name is used (a name that is not
// mapped in the registry).
func NewPipelineFromRegistry(registry map[string]StageFunc, stageNames ...string) Pipeline {
	// this can't be useful and probably indicates a problem
	if len(registry) == 0 {
		return &pipeline{err: ErrorEmptyRegistry}
	}

	// this also can't be useful and probably indicates a problem
	if len(stageNames) == 0 {
		return &pipeline{err: ErrorNoStages}
	}

	p := pipeline{}
	for _, s := range stageNames {
		fn, ok := registry[s]
		if !ok {
			return pipeline{err: ErrorStageNotFound}
		}

		p.stages = append(p.stages, fn)
	}

	return &p
}

// impl
func (p pipeline) Err() error {
	return p.err
}

// impl
func (p pipeline) Execute(state *BoardState, settings Settings, moves []SnakeMove) (bool, *BoardState, error) {
	// Design Detail
	//
	// If the pipeline is in an error state, Execute must return that error
	// because the pipeline is invalid and cannot execute.
	//
	// This is done for API use convenience to satisfy the common pattern
	// of wanting to write NewPipeline().Execute(...).
	//
	// This way you can do that without having to do 2 error checks.
	// It defers errors from construction to being checked on execution.
	if p.err != nil {
		return false, nil, p.err
	}

	// Actually execute
	var ended bool
	var err error
	state = state.Clone()
	for _, fn := range p.stages {
		// execute current stage
		ended, err = fn(state, settings, moves)

		// stop if we hit any errors or if the game is ended
		if err != nil || ended {
			return ended, state, err
		}
	}

	// return the result of the last stage as the final pipeline result
	return ended, state, err
}
