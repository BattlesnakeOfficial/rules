package rules

import (
	"errors"
)

// StageRegistry is a mapping of stage names to stage functions
type StageRegistry map[string]StageFunc

const (
	StageSpawnFoodStandard    = "spawn_food.standard"
	StageGameOverStandard     = "game_over.standard"
	StageStarvationStandard   = "starvation.standard"
	StageFeedSnakesStandard   = "feed_snakes.standard"
	StageMovementStandard     = "movement.standard"
	StageHazardDamageStandard = "hazard_damage.standard"
	StageEliminationStandard  = "elimination.standard"

	StageGameOverSoloSnake                   = "game_over.solo_snake"
	StageGameOverBySquad                     = "game_over.by_squad"
	StageSpawnFoodNoFood                     = "spawn_food.no_food"
	StageSpawnHazardsShrinkMap               = "spawn_hazards.shrink_map"
	StageEliminationResurrectSquadCollisions = "elimination.resurrect_squad_collisions"
	StageModifySnakesAlwaysGrow              = "modify_snakes.always_grow"
	StageMovementWrapBoundaries              = "movement.wrap_boundaries"
	StageModifySnakesShareAttributes         = "modify_snakes.share_attributes"
)

// globalRegistry is a global, default mapping of stage names to stage functions.
// It can be extended by plugins through the use of registration functions.
// Plugins that wish to extend the available game stages should call RegisterPipelineStageError
// to add additional stages.
var globalRegistry = StageRegistry{
	StageSpawnFoodNoFood:                     RemoveFoodConstrictor,
	StageSpawnFoodStandard:                   SpawnFoodStandard,
	StageGameOverSoloSnake:                   GameOverSolo,
	StageGameOverBySquad:                     GameOverSquad,
	StageGameOverStandard:                    GameOverStandard,
	StageHazardDamageStandard:                DamageHazardsStandard,
	StageSpawnHazardsShrinkMap:               PopulateHazardsRoyale,
	StageStarvationStandard:                  ReduceSnakeHealthStandard,
	StageEliminationResurrectSquadCollisions: ResurrectSnakesSquad,
	StageFeedSnakesStandard:                  FeedSnakesStandard,
	StageEliminationStandard:                 EliminateSnakesStandard,
	StageModifySnakesAlwaysGrow:              GrowSnakesConstrictor,
	StageMovementStandard:                    MoveSnakesStandard,
	StageMovementWrapBoundaries:              MoveSnakesWrapped,
	StageModifySnakesShareAttributes:         ShareAttributesSquad,
}

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
		return errors.New("stage has already been registered")
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

// Pipeline is an ordered sequences of game stages which are executed to produce the
// next game state.
//
// If a stage produces an error or an ended game state, the pipeline is halted at that stage.
type Pipeline interface {
	// Execute runs a sequence of stages and produces a next game state.
	//
	// If any stage produces an error or an ended game state, the pipeline
	// immediately stops at that stage.
	//
	// The result is the result of the last stage that was executed.
	//
	Execute(*BoardState, Settings, []SnakeMove) (bool, *BoardState, error)
	// Error can be called to check if the pipeline is in an error state.
	Error() error
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
// 	NewPipelineFromRegistry(r, s, "stage1", "stage2")
// ... will result in stage "stage1" running first, then stage "stage2" running after.
//
// An error will be returned if an unregistered stage name is used (a name that is not
// mapped in the registry).
func NewPipelineFromRegistry(registry map[string]StageFunc, stageNames ...string) Pipeline {
	// this can't be useful and probably indicates a problem
	if len(registry) == 0 {
		return &pipeline{err: errors.New("empty registry")}
	}

	// this also can't be useful and probably indicates a problem
	if len(stageNames) == 0 {
		return &pipeline{err: errors.New("no stages")}
	}

	p := pipeline{}
	for _, s := range stageNames {
		fn, ok := registry[s]
		if !ok {
			return pipeline{err: errors.New("stage not found")}
		}

		p.stages = append(p.stages, fn)
	}

	return &p
}

// impl
func (p pipeline) Error() error {
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
