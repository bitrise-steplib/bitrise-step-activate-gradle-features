package step_test

import (
	"testing"

	"github.com/bitrise-io/bitrise-plugins-annotations/service"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	utilsMocks "github.com/bitrise-io/go-utils/v2/mocks"
	"github.com/bitrise-steplib/bitrise-step-activate-gradle-features/step"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Step(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		logger := &utilsMocks.Logger{}
		logger.On("EnableDebugLog", true).Return().Once()
		logger.On("Println", mock.Anything).Return().Once()
		logger.On("Debugf", mock.Anything, mock.Anything).Return()
		logger.On("Infof", step.NoFeaturesEnabledMsg).Return().Once()
		logger.On("Infof", step.GradleFeaturesActivatedMsg).Return().Once()

		envRepo := NewMockEnvRepo()
		envRepo.Set("build_cache_enabled", "true")           //nolint: errcheck
		envRepo.Set("build_cache_push", "true")              //nolint: errcheck
		envRepo.Set("build_cache_validation_level", "error") //nolint: errcheck
		envRepo.Set("BITRISEIO_BUILD_CACHE_ENABLED", "true") //nolint: errcheck
		envRepo.Set("test_distribution_enabled", "true")     //nolint: errcheck
		envRepo.Set("test_distribution_shard_size", "50")    //nolint: errcheck
		envRepo.Set("verbose", "true")

		command := &MockCommand{}

		sut := step.New(
			logger,
			stepconf.NewInputParser(envRepo),
			envRepo,
			func(annotation service.Annotation) error { return nil },
			command,
		)

		err := sut.Run()

		assert.Nil(t, err)
		assert.Equal(t, 1, command.Executed)
		assert.Equal(t, command.Args, []string{
			"activate",
			"gradle",
			"--debug",
			"--analytics",
			"--cache",
			"--cache-push=true",
			"--cache-validation=error",
			"--test-distribution",
			"--test-distribution-shard-size=50",
		})
	})

	t.Run("Failed to parse input", func(t *testing.T) {
		logger := &utilsMocks.Logger{}
		command := &MockCommand{}
		envRepo := NewMockEnvRepo()

		sut := step.New(
			logger,
			stepconf.NewInputParser(envRepo),
			envRepo,
			func(annotation service.Annotation) error { return nil },
			command,
		)

		err := sut.Run()
		assert.ErrorContains(t, err, step.FailedToParseInputsMsg)
		assert.Equal(t, 0, command.Executed)
	})

	t.Run("No features enabled", func(t *testing.T) {
		envRepo := NewMockEnvRepo()
		envRepo.Set("verbose", "false")                   //nolint: errcheck
		envRepo.Set("BITRISE_ANALYTICS_DISABLED", "true") //nolint: errcheck

		logger := &utilsMocks.Logger{}
		logger.On("EnableDebugLog", false).Return().Once()
		logger.On("Println", mock.Anything).Return().Once()
		logger.On("Debugf", mock.Anything, mock.Anything).Return()
		logger.On("Infof", step.NoFeaturesEnabledMsg).Return().Once()

		command := &MockCommand{}

		sut := step.New(
			logger,
			stepconf.NewInputParser(envRepo),
			envRepo,
			func(annotation service.Annotation) error { return nil },
			command,
		)

		err := sut.Run()
		assert.Nil(t, err)
		assert.Equal(t, 0, command.Executed)
	})

	t.Run("Failed to activate", func(t *testing.T) {
		envRepo := NewMockEnvRepo()
		envRepo.Set("test_distribution_enabled", "true")  //nolint: errcheck
		envRepo.Set("test_distribution_shard_size", "50") //nolint: errcheck
		envRepo.Set("verbose", "false")                   //nolint: errcheck

		logger := &utilsMocks.Logger{}
		logger.On("EnableDebugLog", false).Return().Once()
		logger.On("Println", mock.Anything).Return().Once()
		logger.On("Debugf", mock.Anything, mock.Anything).Return()
		logger.On("Infof", step.NoFeaturesEnabledMsg).Return().Once()

		command := &MockCommand{
			ExecutionError: assert.AnError,
		}

		sut := step.New(
			logger,
			stepconf.NewInputParser(envRepo),
			envRepo,
			func(annotation service.Annotation) error { return nil },
			command,
		)

		err := sut.Run()
		assert.EqualError(t, err, step.FailedToActivateMsg+": "+assert.AnError.Error())
		assert.Equal(t, 1, command.Executed)
	})
}
