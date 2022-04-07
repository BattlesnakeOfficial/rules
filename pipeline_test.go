package rules_test

import (
	"errors"
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/stretchr/testify/require"
)

func TestPipeline(t *testing.T) {
	r := rules.StageRegistry{}

	// test empty registry error
	_, _, err := rules.NewPipelineFromRegistry(r).Execute(nil, rules.Settings{}, nil)
	require.Equal(t, errors.New("empty registry"), err)

	// test empty stages names error
	r.RegisterPipelineStage("astage", mockStageFn(false, nil))
	_, _, err = rules.NewPipelineFromRegistry(r).Execute(&rules.BoardState{}, rules.Settings{}, nil)
	require.Equal(t, errors.New("no stages"), err)

	// test that an unregistered stage name errors
	_, _, err = rules.NewPipelineFromRegistry(r, "doesntexist").Execute(&rules.BoardState{}, rules.Settings{}, nil)
	require.Equal(t, errors.New("stage not found"), err)

	// simplest case - one stage
	ended, next, err := rules.NewPipelineFromRegistry(r, "astage").Execute(&rules.BoardState{}, rules.Settings{}, nil)
	require.NoError(t, err)
	require.NoError(t, err)
	require.NotNil(t, next)
	require.False(t, ended)

	// test that the pipeline short-circuits for a stage that errors
	r.RegisterPipelineStage("errors", mockStageFn(false, errors.New("")))
	ended, next, err = rules.NewPipelineFromRegistry(r, "errors", "astage").Execute(&rules.BoardState{}, rules.Settings{}, nil)
	require.Error(t, err)
	require.NotNil(t, next)
	require.False(t, ended)

	// test that the pipeline short-circuits for a stage that ends
	r.RegisterPipelineStage("ends", mockStageFn(true, nil))
	ended, next, err = rules.NewPipelineFromRegistry(r, "ends", "astage").Execute(&rules.BoardState{}, rules.Settings{}, nil)
	require.NoError(t, err)
	require.NotNil(t, next)
	require.True(t, ended)

	// test that the pipeline runs normally for multiple stages
	ended, next, err = rules.NewPipelineFromRegistry(r, "astage", "ends").Execute(&rules.BoardState{}, rules.Settings{}, nil)
	require.NoError(t, err)
	require.NotNil(t, next)
	require.True(t, ended)
}

func TestStageRegistry(t *testing.T) {
	sr := rules.StageRegistry{}

	// register a stage without error
	require.NoError(t, sr.RegisterPipelineStageError("test", mockStageFn(false, nil)))
	require.Contains(t, sr, "test")

	// error on duplicate
	require.Error(t, sr.RegisterPipelineStageError("test", mockStageFn(false, nil)))

	// register another stage with no error
	require.NoError(t, sr.RegisterPipelineStageError("other", mockStageFn(false, nil)))
	require.Contains(t, sr, "other")

	// register stage
	sr.RegisterPipelineStage("last", mockStageFn(false, nil))
	require.Contains(t, sr, "last")

	// register existing stage (should just be okay and not panic or anything)
	sr.RegisterPipelineStage("test", mockStageFn(false, nil))
}

func mockStageFn(ended bool, err error) rules.StageFunc {
	return func(b *rules.BoardState, settings rules.Settings, moves []rules.SnakeMove) (bool, error) {
		return ended, err
	}
}
