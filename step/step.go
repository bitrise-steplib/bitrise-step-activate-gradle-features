package step

import (
	"fmt"

	"github.com/bitrise-io/bitrise-plugins-annotations/service"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/bitrise-step-activate-gradle-features/step/features"
)

const (
	FailedToParseInputsMsg     = "failed to parse inputs"
	NoFeaturesEnabledMsg       = "No features enabled"
	FailedToActivateMsg        = "failed to activate Gradle features"
	GradleFeaturesActivatedMsg = "Gradle features activated successfully"
)

type Input struct {
	Verbose bool `env:"verbose,required"`
}

type Step struct {
	logger      log.Logger
	inputParser stepconf.InputParser
	envRepo     env.Repository
	annotator   func(annotation service.Annotation) error
	command     Command
}

func New(
	logger log.Logger,
	inputParser stepconf.InputParser,
	envRepo env.Repository,
	annotator func(annotation service.Annotation) error,
	command Command,
) Step {
	return Step{
		logger:      logger,
		inputParser: inputParser,
		envRepo:     envRepo,
		annotator:   annotator,
		command:     command,
	}
}

func (step Step) Run() error {
	var input Input
	if err := step.inputParser.Parse(&input); err != nil {
		return fmt.Errorf(FailedToParseInputsMsg+": %w", err)
	}
	step.logger.EnableDebugLog(input.Verbose)

	collectedFeatures := step.collectFeatures()
	var hasEnabledFeatures bool
	for _, feature := range collectedFeatures {
		stepconf.Print(feature)
		hasEnabledFeatures = true
	}

	stepconf.Print(input)
	step.logger.Println()

	if !hasEnabledFeatures {
		step.logger.Infof(NoFeaturesEnabledMsg)
		return nil
	}

	if err := step.activate(input, collectedFeatures); err != nil {
		return fmt.Errorf(FailedToActivateMsg+": %w", err)
	}

	step.logger.Infof(GradleFeaturesActivatedMsg)

	return nil
}

func (step Step) collectFeatures() []Feature {
	collected := []Feature{}

	analyticsFeature := features.AnalyticsFeature(step.envRepo, step.logger)
	if analyticsFeature != nil {
		collected = append(collected, analyticsFeature)
	}

	buildCacheFeature := features.BuildCacheFeature(step.inputParser, step.envRepo, step.logger, step.annotator)
	if buildCacheFeature != nil {
		collected = append(collected, buildCacheFeature)
	}

	testDistributionFeature := features.TestDistributionFeature(step.inputParser, step.envRepo, step.logger)
	if testDistributionFeature != nil {
		collected = append(collected, testDistributionFeature)
	}

	return collected
}

func (step Step) activate(input Input, collectedFeatures []Feature) error {
	args := []string{"activate", "gradle"}
	if input.Verbose {
		args = append(args, "--debug")
	}
	for _, feature := range collectedFeatures {
		args = append(args, feature.CLIFlags()...)
	}

	step.command.SetArgs(args)
	return step.command.Execute()
}
