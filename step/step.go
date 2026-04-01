package step

import (
	"fmt"

	"github.com/bitrise-io/bitrise-plugins-annotations/service"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/bitrise-step-activate-react-native-features/step/features"
)

const (
	FailedToParseInputsMsg          = "failed to parse inputs"
	NoFeaturesEnabledMsg            = "No features enabled"
	FailedToActivateMsg             = "failed to activate React Native features"
	ReactNativeFeaturesActivatedMsg = "React Native features activated successfully"
)

type Input struct {
	Verbose bool `env:"verbose,required"`
}

type Step struct {
	logger      log.Logger
	inputParser stepconf.InputParser
	annotator   func(annotation service.Annotation) error
	command     Command
}

func New(
	logger log.Logger,
	inputParser stepconf.InputParser,
	annotator func(annotation service.Annotation) error,
	command Command,
) Step {
	return Step{
		logger:      logger,
		inputParser: inputParser,
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

	collectedCacheFeatures := step.collectCacheFeatures()

	var hasEnabledFeatures bool
	for _, feature := range collectedCacheFeatures {
		stepconf.Print(feature)
		hasEnabledFeatures = true
	}

	stepconf.Print(input)
	step.logger.Println()

	if !hasEnabledFeatures {
		step.logger.Infof(NoFeaturesEnabledMsg)
		return nil
	}

	if err := step.activate(input, collectedCacheFeatures); err != nil {
		return fmt.Errorf(FailedToActivateMsg+": %w", err)
	}

	step.logger.Infof(ReactNativeFeaturesActivatedMsg)

	return nil
}

func (step Step) collectCacheFeatures() []CacheFeature {
	collected := []CacheFeature{}

	if f := features.CPPCacheFeature(step.inputParser, step.logger); f != nil {
		collected = append(collected, f)
	}
	if f := features.XcodeCacheFeature(step.inputParser, step.logger); f != nil {
		collected = append(collected, f)
	}
	if f := features.GradleCacheFeature(step.inputParser, step.logger); f != nil {
		collected = append(collected, f)
	}

	return collected
}

func (step Step) activate(input Input, cacheFeatures []CacheFeature) error {
	for _, feature := range cacheFeatures {
		args := feature.ActivateArgs()
		if input.Verbose {
			args = append(args, "--debug")
		}
		step.command.SetArgs(args)
		if err := step.command.Execute(); err != nil {
			return err
		}
	}

	for _, feature := range cacheFeatures {
		serviceArgs := feature.ServiceArgs()
		if serviceArgs == nil {
			continue
		}
		args := serviceArgs
		if input.Verbose {
			args = append(args, "--debug")
		}
		go func() {
			step.command.SetArgs(args)
			if err := step.command.Execute(); err != nil {
				step.logger.Warnf("background service stopped with error: %s", err)
			}
		}()
	}

	return nil
}
