package features_test

import (
	"testing"

	utilsMocks "github.com/bitrise-io/go-utils/v2/mocks"
	"github.com/bitrise-steplib/bitrise-step-activate-gradle-features/step/features"
	"github.com/stretchr/testify/assert"
)

func Test_AnalyticsFeature(t *testing.T) {

	t.Run("Happy path", func(t *testing.T) {
		logger := &utilsMocks.Logger{}
		logger.On("Debugf", features.AnalyticsEnabledMsg).Return().Once()

		envRepo := NewMockEnvRepo()

		actual := features.AnalyticsFeature(
			envRepo,
			logger,
		)

		assert.Equal(t, features.Analytics{}, *actual)
	})

	t.Run("Disabled", func(t *testing.T) {
		logger := &utilsMocks.Logger{}

		envRepo := NewMockEnvRepo()
		envRepo.Set("BITRISE_ANALYTICS_DISABLED", "true")

		actual := features.AnalyticsFeature(
			envRepo,
			logger,
		)

		assert.Nil(t, actual)
	})
}

func Test_AnalyticsCLIFlags(t *testing.T) {
	t.Run("Happy Path", func(t *testing.T) {
		td := features.Analytics{}

		actual := td.CLIFlags()
		expected := []string{
			"--analytics",
		}

		assert.Equal(t, expected, actual)
	})
}
