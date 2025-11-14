package main

import (
	"os"

	CLI "github.com/bitrise-io/bitrise-build-cache-cli/cmd/common"
	gradleCLI "github.com/bitrise-io/bitrise-build-cache-cli/cmd/gradle"
	"github.com/bitrise-io/bitrise-plugins-annotations/service"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/exitcode"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/bitrise-step-activate-gradle-features/step"
)

func main() {
	exitCode := run()
	os.Exit(int(exitCode))
}

func run() exitcode.ExitCode {
	// make the gradle commands registered (have to make the init called)
	_ = gradleCLI.ActivateGradleCmd.Commands()

	logger := log.NewLogger()
	envRepo := env.NewRepository()
	inputParser := stepconf.NewInputParser(envRepo)

	stepInstance := step.New(
		logger,
		inputParser,
		envRepo,
		service.Annotate,
		CLI.RootCmd,
	)
	if err := stepInstance.Run(); err != nil {
		logger.Errorf(err.Error())
		return exitcode.Failure
	}

	return exitcode.Success
}
