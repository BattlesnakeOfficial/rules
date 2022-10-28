package rules

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPipelineRuleset(t *testing.T) {
	r := StageRegistry{
		"doesnt_end": mockStageFn(false, nil),
		"ends":       mockStageFn(true, nil),
	}

	// Name/Error methods
	p := NewPipelineFromRegistry(r, "404doesntexist")
	pr := pipelineRuleset{
		name:     "test",
		pipeline: p,
	}
	require.Equal(t, "test", pr.Name())
	require.Equal(t, ErrorStageNotFound, pr.Err())

	// test game over when it does end
	p = NewPipelineFromRegistry(r, "doesnt_end", "ends")
	pr = pipelineRuleset{
		name:     "test",
		pipeline: p,
	}
	ended, _, err := pr.Execute(&BoardState{}, nil)
	require.NoError(t, err)
	require.True(t, ended)

	// Test game over when it doesn't end
	p = NewPipelineFromRegistry(r, "doesnt_end")
	pr = pipelineRuleset{
		name:     "test",
		pipeline: p,
	}
	ended, _, err = pr.Execute(&BoardState{}, nil)
	require.NoError(t, err)
	require.False(t, ended)

	// test a stage that adds food, except on initialization
	r.RegisterPipelineStage("add_food", func(bs *BoardState, s Settings, sm []SnakeMove) (bool, error) {
		if IsInitialization(bs, s, sm) {
			return false, nil
		}
		bs.Food = append(bs.Food, Point{X: 0, Y: 0})
		return false, nil
	})
	b := &BoardState{}
	p = NewPipelineFromRegistry(r, "add_food")
	pr = pipelineRuleset{
		name:     "test",
		pipeline: p,
	}
	require.Empty(t, b.Food)
	_, b, err = pr.Execute(b, nil)
	require.NoError(t, err)
	require.Empty(t, b.Food, "food should not be added on initialisation phase")
	_, b, err = pr.Execute(b, mockSnakeMoves())
	require.NoError(t, err)
	require.NotEmpty(t, b.Food, "fodo should be added now")
}

func TestPipelineGlobals(t *testing.T) {
	oldReg := globalRegistry
	globalRegistry = StageRegistry{}

	// ensure that we can register a function without errors
	RegisterPipelineStage("test", mockStageFn(false, nil))
	require.Contains(t, globalRegistry, "test")

	// ensure that the global registry panics if you register an existing stage name
	require.Panics(t, func() {
		RegisterPipelineStage("test", mockStageFn(false, nil))
	})
	RegisterPipelineStage("other", mockStageFn(true, nil)) // otherwise should not panic

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
