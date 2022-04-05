package rules

import (
	"errors"
)

// StageRegistry is a mapping of stage names to stage functions
type StageRegistry map[string]StageFunc

// globalRegistry is a global mapping of stage names to stage functions
var globalRegistry = StageRegistry{
	"reducehealth.standard":   ReduceSnakeHealthStandard,
	"hazarddamage.standard":   DamageHazardsStandard,
	"eatfood.standard":        FeedSnakesStandard,
	"placefood.standard":      SpawnFoodStandard,
	"placehazard.royale":      PopulateHazardsRoyale,
	"removefood.constrictor":  RemoveFoodConstrictor,
	"growsnake.constrictor":   GrowSnakesConstrictor,
	"eliminatesnake.standard": EliminateSnakesStandard,
	"movement.standard":       MoveSnakesStandard,
	"movement.wrapped":        MoveSnakesWrapped,
	"gameover.standard":       GameOverStandard,
	"gameover.solo":           GameOverSolo,
	"gameover.squad":          GameOverSquad,
	"collision.squad":         ResurrectSnakesSquad,
	"sharedattr.squad":        ShareAttributesSquad,
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
// If a stage has already been mapped it will be overwritten by the newly
// registered function.
func RegisterPipelineStage(s string, fn StageFunc) {
	globalRegistry.RegisterPipelineStage(s, fn)
}

// RegisterPipelineStageError adds a stage to the global stage registry.
// If a stage has already been mapped an error will be returned.
func RegisterPipelineStageError(s string, fn StageFunc) error {
	return globalRegistry.RegisterPipelineStageError(s, fn)
}

// Pipeline is an ordered sequences of game stages which are executed to produce the
// next game state.
//
// If a stage produces an error or an ended game state, the pipeline is halted at that stage.
type Pipeline struct {
	// stages is a list of stages that should be executed from slice start to end
	stages []StageFunc
}

// NewPipeline constructs an instance of Pipeline using the global registry.
// It is a convenience wrapper for NewPipelineFromRegistry when you want
// to use the default, global registry.
func NewPipeline(stageNames ...string) (*Pipeline, error) {
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
func NewPipelineFromRegistry(registry map[string]StageFunc, stageNames ...string) (*Pipeline, error) {
	// this can't be useful and probably indicates a problem
	if len(registry) == 0 {
		return nil, errors.New("empty registry")
	}

	// this also can't be useful and probably indicates a problem
	if len(stageNames) == 0 {
		return nil, errors.New("no stages")
	}

	p := &Pipeline{}
	for _, s := range stageNames {
		fn, ok := registry[s]
		if !ok {
			return nil, errors.New("stage not found")
		}

		p.stages = append(p.stages, fn)
	}

	return p, nil
}

// Append appends the stages from another pipeline to the end of this pipeline.
// It is similar to the standard append in that the backing slice of stages will be re-used, when possible.
func (p *Pipeline) Append(p2 *Pipeline) *Pipeline {
	return &Pipeline{
		stages: append(p.stages, p2.stages...),
	}
}

// Execute runs all of the pipeline stages and produces a next game state.
// If any stage produces an error or an ended game state, the pipeline
// immediately stops at that stage.
// The result is always the result of the last stage that was executed.
func (p *Pipeline) Execute(state *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
	var ended bool
	var err error
	for _, fn := range p.stages {
		// execute current stage
		ended, err = fn(state, settings, moves)

		// stop if we hit any errors or if the game is ended
		if err != nil || ended {
			return ended, err
		}
	}

	// return the result of the last stage as the final pipeline result
	return ended, err
}