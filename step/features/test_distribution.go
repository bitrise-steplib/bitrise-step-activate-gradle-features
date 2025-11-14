package features

import (
	"fmt"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

const (
	TestDistributionCheckMsg      = "Checking whether Bitrise Test Distribution is activated for this workspace ..."
	TestDistributionParsingFailed = "Test Distribution feature is not configured: %s"
	TestDistributionDisabledMsg   = "Test Distribution feature is not enabled"
)

type TestDistribution struct {
	Enabled   bool `env:"test_distribution_enabled,required"`
	ShardSize int  `env:"test_distribution_shard_size,required"`
}

func TestDistributionFeature(
	inputParser stepconf.InputParser,
	envRepo env.Repository,
	logger log.Logger,
) *TestDistribution {
	logger.Debugf(TestDistributionCheckMsg)
	var td TestDistribution
	if err := inputParser.Parse(&td); err != nil {
		logger.Debugf(TestDistributionParsingFailed, err)
		return nil
	}

	if !td.Enabled {
		logger.Debugf(TestDistributionDisabledMsg)
		return nil
	}

	return &td
}

func (td *TestDistribution) CLIFlags() []string {
	if !td.Enabled {
		return []string{}
	}

	return []string{
		"--test-distribution",
		fmt.Sprintf("--test-distribution-shard-size=%d", td.ShardSize),
	}
}
