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
	p := rules.NewPipelineFromRegistry(r)
	require.Equal(t, rules.ErrorEmptyRegistry, p.Err())
	_, _, err := p.Execute(nil, rules.Settings{}, nil)
	require.Equal(t, rules.ErrorEmptyRegistry, err)

	// test empty stages names error
	r.RegisterPipelineStage("astage", mockStageFn(false, nil))
	p = rules.NewPipelineFromRegistry(r)
	require.Equal(t, rules.ErrorNoStages, p.Err())
	_, _, err = p.Execute(rules.NewBoardState(0, 0), rules.Settings{}, nil)
	require.Equal(t, rules.ErrorNoStages, err)

	// test that an unregistered stage name errors
	p = rules.NewPipelineFromRegistry(r, "doesntexist")
	_, _, err = p.Execute(rules.NewBoardState(0, 0), rules.Settings{}, nil)
	require.Equal(t, rules.ErrorStageNotFound, p.Err())
	require.Equal(t, rules.ErrorStageNotFound, err)

	// simplest case - one stage
	ended, next, err := rules.NewPipelineFromRegistry(r, "astage").Execute(rules.NewBoardState(0, 0), rules.Settings{}, nil)
	require.NoError(t, err)
	require.NoError(t, err)
	require.NotNil(t, next)
	require.False(t, ended)

	// test that the pipeline short-circuits for a stage that errors
	r.RegisterPipelineStage("errors", mockStageFn(false, errors.New("")))
	ended, next, err = rules.NewPipelineFromRegistry(r, "errors", "astage").Execute(rules.NewBoardState(0, 0), rules.Settings{}, nil)
	require.Error(t, err)
	require.NotNil(t, next)
	require.False(t, ended)

	// test that the pipeline short-circuits for a stage that ends
	r.RegisterPipelineStage("ends", mockStageFn(true, nil))
	ended, next, err = rules.NewPipelineFromRegistry(r, "ends", "astage").Execute(rules.NewBoardState(0, 0), rules.Settings{}, nil)
	require.NoError(t, err)
	require.NotNil(t, next)
	require.True(t, ended)

	// test that the pipeline runs normally for multiple stages
	ended, next, err = rules.NewPipelineFromRegistry(r, "astage", "ends").Execute(rules.NewBoardState(0, 0), rules.Settings{}, nil)
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
	var e rules.RulesetError
	err := sr.RegisterPipelineStageError("test", mockStageFn(false, nil))
	require.Error(t, err)
	require.True(t, errors.As(err, &e), "error should be a RulesetError")
	require.Equal(t, "stage 'test' has already been registered", err.Error())

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
