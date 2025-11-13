package features_test

import (
	"testing"

	"github.com/bitrise-io/bitrise-plugins-annotations/service"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	utilsMocks "github.com/bitrise-io/go-utils/v2/mocks"
	"github.com/bitrise-steplib/bitrise-step-activate-gradle-features/step/features"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_BuildCacheFeature(t *testing.T) {
	makeLogger := func() *utilsMocks.Logger {
		logger := &utilsMocks.Logger{}
		logger.On("Debugf", mock.Anything, mock.Anything).Return()
		logger.On("Errorf", mock.Anything, mock.Anything).Return()
		logger.On("Infof", mock.Anything).Return()
		return logger
	}

	t.Run("Happy path", func(t *testing.T) {
		envRepo := NewMockEnvRepo()
		envRepo.Set("build_cache_enabled", "true")
		envRepo.Set("build_cache_push", "true")
		envRepo.Set("build_cache_validation_level", "error")
		envRepo.Set("BITRISEIO_BUILD_CACHE_ENABLED", "true")

		actual := features.BuildCacheFeature(
			stepconf.NewInputParser(envRepo),
			envRepo,
			makeLogger(),
			func(annotation service.Annotation) error { return nil },
		)

		assert.Equal(t, features.BuildCache{
			Enabled:         true,
			Push:            true,
			ValidationLevel: "error",
		}, *actual)
	})

	t.Run("Missing envs", func(t *testing.T) {
		envRepo := NewMockEnvRepo()
		// missing envs

		logger := &utilsMocks.Logger{}
		logger.On("Debugf", features.BuildCacheCheckMsg).Return().Once()
		logger.On("Debugf", features.BuildCacheParsingFailed, mock.Anything).Return().Once()

		actual := features.BuildCacheFeature(
			stepconf.NewInputParser(envRepo),
			envRepo,
			makeLogger(),
			func(annotation service.Annotation) error { return nil },
		)

		assert.Nil(t, actual)
	})

	t.Run("Disabled", func(t *testing.T) {
		envRepo := NewMockEnvRepo()
		envRepo.Set("build_cache_enabled", "false")
		envRepo.Set("build_cache_push", "true")
		envRepo.Set("build_cache_validation_level", "error")

		logger := &utilsMocks.Logger{}
		logger.On("Debugf", features.BuildCacheCheckMsg).Return().Once()
		logger.On("Debugf", features.BuildCacheDisabledMsg).Return().Once()

		actual := features.BuildCacheFeature(
			stepconf.NewInputParser(envRepo),
			envRepo,
			makeLogger(),
			func(annotation service.Annotation) error { return nil },
		)

		assert.Nil(t, actual)
	})

	t.Run("Inactive", func(t *testing.T) {
		envRepo := NewMockEnvRepo()
		envRepo.Set("build_cache_enabled", "true")
		envRepo.Set("build_cache_push", "true")
		envRepo.Set("build_cache_validation_level", "error")
		// missing BITRISEIO_BUILD_CACHE_ENABLED

		logger := &utilsMocks.Logger{}
		logger.On("Debugf", features.BuildCacheCheckMsg).Return().Once()
		logger.On("Errorf", features.BuildCacheNotActivatedMsg).Return().Once()

		actual := features.BuildCacheFeature(
			stepconf.NewInputParser(envRepo),
			envRepo,
			makeLogger(),
			func(annotation service.Annotation) error { return nil },
		)

		assert.Nil(t, actual)
	})
}

func Test_BuildCacheCLIFlags(t *testing.T) {
	t.Run("Happy Path", func(t *testing.T) {
		bc := features.BuildCache{
			Enabled:         true,
			Push:            true,
			ValidationLevel: "error",
		}

		actual := bc.CLIFlags()
		expected := []string{
			"--cache",
			"--cache-push=true",
			"--cache-validation=error",
		}

		assert.Equal(t, expected, actual)
	})

	t.Run("Disabled", func(t *testing.T) {
		bc := features.BuildCache{
			Enabled:         false,
			Push:            true,
			ValidationLevel: "error",
		}

		actual := bc.CLIFlags()
		expected := []string{}

		assert.Equal(t, expected, actual)
	})
}
