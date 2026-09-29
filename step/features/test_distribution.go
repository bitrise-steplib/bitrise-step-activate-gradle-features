package features

import (
	"errors"
	"fmt"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

const (
	TestDistributionCheckMsg       = "Checking whether Bitrise Test Distribution is activated for this workspace ..."
	TestDistributionParsingFailed  = "Test Distribution feature is not configured: %s"
	TestDistributionDisabledMsg    = "Test Distribution feature is not enabled"
	TestDistributionMissingPoolMsg = "test_distribution_pool is required when test_distribution_enabled is true"
)

type TestDistribution struct {
	Enabled   bool   `env:"test_distribution_enabled,required"`
	ShardSize int    `env:"test_distribution_shard_size,required"`
	PoolName  string `env:"test_distribution_pool"`
}

func TestDistributionFeature(
	inputParser stepconf.InputParser,
	envRepo env.Repository,
	logger log.Logger,
) (*TestDistribution, error) {
	logger.Debugf(TestDistributionCheckMsg)
	var td TestDistribution
	if err := inputParser.Parse(&td); err != nil {
		logger.Debugf(TestDistributionParsingFailed, err)

		return nil, nil
	}

	if !td.Enabled {
		logger.Debugf(TestDistributionDisabledMsg)

		return nil, nil
	}

	if td.PoolName == "" {
		return nil, errors.New(TestDistributionMissingPoolMsg)
	}

	return &td, nil
}

func (td *TestDistribution) CLIFlags() []string {
	if !td.Enabled {
		return []string{}
	}

	return []string{
		"--test-distribution",
		fmt.Sprintf("--test-distribution-shard-size=%d", td.ShardSize),
		fmt.Sprintf("--test-distribution-pool=%s", td.PoolName),
	}
}
