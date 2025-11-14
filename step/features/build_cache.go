package features

import (
	"fmt"

	"github.com/bitrise-io/bitrise-plugins-annotations/service"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

const (
	unavailableMsg = `Bitrise Build Cache is not activated in this build.

You have added the **Activate Bitrise Build Cache for Gradle** add-on step to your workflow.

However, you don't have an activate Bitrise Build Cache Trial or Subscription for the current workspace yet.

You can activate a Trial at [app.bitrise.io/build-cache](https://app.bitrise.io/build-cache),
or contact us at [support@bitrise.io](mailto:support@bitrise.io) to activate it.`

	BuildCacheCheckMsg        = "Checking whether Bitrise Build Cache is activated for this workspace ..."
	BuildCacheParsingFailed   = "Build Cache feature is not configured: %s"
	BuildCacheDisabledMsg     = "Build Cache feature is not enabled"
	BuildCacheNotActivatedMsg = "Bitrise Build Cache is not activated for this workspace but was enabled"
)

type BuildCache struct {
	Enabled         bool   `env:"build_cache_enabled,required"`
	Push            bool   `env:"build_cache_push,required"`
	ValidationLevel string `env:"build_cache_validation_level,opt[none,warning,error]"`
}

func BuildCacheFeature(
	inputParser stepconf.InputParser,
	envRepo env.Repository,
	logger log.Logger,
	annotator func(annotation service.Annotation) error,
) *BuildCache {
	logger.Debugf(BuildCacheCheckMsg)

	var bc BuildCache
	if err := inputParser.Parse(&bc); err != nil {
		logger.Debugf(BuildCacheParsingFailed, err)
		return nil
	}

	if !bc.Enabled {
		logger.Debugf(BuildCacheDisabledMsg)
		return nil
	}

	if envRepo.Get("BITRISEIO_BUILD_CACHE_ENABLED") != "true" {
		_ = annotator(service.Annotation{
			Context:  "",
			Markdown: unavailableMsg,
			Style:    "error",
		})
		logger.Errorf(BuildCacheNotActivatedMsg)
		return nil
	}

	return &bc
}

func (bc *BuildCache) CLIFlags() []string {
	if !bc.Enabled {
		return []string{}
	}

	return []string{
		"--cache",
		fmt.Sprintf("--cache-push=%t", bc.Push),
		fmt.Sprintf("--cache-validation=%s", bc.ValidationLevel),
	}
}
