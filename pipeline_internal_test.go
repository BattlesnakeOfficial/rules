package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPipelineGlboals(t *testing.T) {
	oldReg := globalRegistry
	globalRegistry = StageRegistry{}

	// ensure that we can register a function without errors
	RegisterPipelineStage("test", mockStageFn(false, nil))
	require.Contains(t, globalRegistry, "test")

	// ensure that registry errors for existing stage names
	require.Error(t, RegisterPipelineStageError("test", mockStageFn(false, nil)))
	require.NoError(t, RegisterPipelineStageError("other", mockStageFn(true, nil)))

	// ensure that we can build a pipeline using the global registry
	p := NewPipeline("test", "other")
	require.NotNil(t, p)

	// ensure that it runs okay too
	ended, next, err := p.Execute(&BoardState{}, Settings{}, nil)
	require.NoError(t, err)
	require.NotNil(t, next)
	require.True(t, ended)

	globalRegistry = oldReg
}

func mockStageFn(ended bool, err error) StageFunc {
	return func(b *BoardState, settings Settings, moves []SnakeMove) (bool, error) {
		return ended, err
	}
}
