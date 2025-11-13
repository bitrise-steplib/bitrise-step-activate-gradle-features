package features_test

import (
	"testing"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	utilsMocks "github.com/bitrise-io/go-utils/v2/mocks"
	"github.com/bitrise-steplib/bitrise-step-activate-gradle-features/step/features"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_TestDistributionFeature(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		envRepo := NewMockEnvRepo()
		envRepo.Set("test_distribution_enabled", "true")  //nolint: errcheck
		envRepo.Set("test_distribution_shard_size", "50") //nolint: errcheck

		logger := &utilsMocks.Logger{}
		logger.On("Debugf", mock.Anything, mock.Anything).Return()
		logger.On("Errorf", mock.Anything, mock.Anything).Return()
		logger.On("Infof", mock.Anything).Return()

		actual := features.TestDistributionFeature(
			stepconf.NewInputParser(envRepo),
			envRepo,
			logger,
		)

		assert.Equal(t, features.TestDistribution{
			Enabled:   true,
			ShardSize: 50,
		}, *actual)
	})

	t.Run("Missing envs", func(t *testing.T) {
		envRepo := NewMockEnvRepo()
		// missing envs

		logger := &utilsMocks.Logger{}
		logger.On("Debugf", features.TestDistributionCheckMsg).Return().Once()
		logger.On("Debugf", features.TestDistributionParsingFailed, mock.Anything).Return().Once()

		actual := features.TestDistributionFeature(
			stepconf.NewInputParser(envRepo),
			envRepo,
			logger,
		)

		assert.Nil(t, actual)
	})

	t.Run("Disabled", func(t *testing.T) {
		envRepo := NewMockEnvRepo()
		envRepo.Set("test_distribution_enabled", "false") //nolint: errcheck
		envRepo.Set("test_distribution_shard_size", "50") //nolint: errcheck

		logger := &utilsMocks.Logger{}
		logger.On("Debugf", features.TestDistributionCheckMsg).Return().Once()
		logger.On("Debugf", features.TestDistributionDisabledMsg).Return().Once()

		actual := features.TestDistributionFeature(
			stepconf.NewInputParser(envRepo),
			envRepo,
			logger,
		)

		assert.Nil(t, actual)
	})
}

func Test_TestDistributionCLIFlags(t *testing.T) {
	t.Run("Happy Path", func(t *testing.T) {
		td := features.TestDistribution{
			Enabled:   true,
			ShardSize: 50,
		}

		actual := td.CLIFlags()
		expected := []string{
			"--test-distribution",
			"--test-distribution-shard-size=50",
		}

		assert.Equal(t, expected, actual)
	})

	t.Run("Disabled", func(t *testing.T) {
		td := features.TestDistribution{
			Enabled:   false,
			ShardSize: 50,
		}

		actual := td.CLIFlags()
		expected := []string{}

		assert.Equal(t, expected, actual)
	})
}
