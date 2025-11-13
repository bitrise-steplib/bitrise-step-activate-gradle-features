package features

import (
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

const (
	AnalyticsEnabledMsg = "Analytics feature is enabled"
)

type Analytics struct{}

func AnalyticsFeature(
	envRepo env.Repository,
	logger log.Logger,
) *Analytics {
	if envRepo.Get("BITRISE_ANALYTICS_DISABLED") == "true" {
		return nil
	}

	logger.Debugf(AnalyticsEnabledMsg)
	return &Analytics{}
}

func (a *Analytics) CLIFlags() []string {
	return []string{
		"--analytics",
	}
}
