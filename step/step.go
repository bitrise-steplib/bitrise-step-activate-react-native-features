package step

import (
	"fmt"

	"github.com/bitrise-io/bitrise-plugins-annotations/service"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
)

const (
	FailedToParseInputsMsg          = "failed to parse inputs"
	NoFeaturesEnabledMsg            = "No features enabled"
	FailedToActivateMsg             = "failed to activate React Native features"
	ReactNativeFeaturesActivatedMsg = "React Native features activated successfully"
)

type Input struct {
	Verbose            bool `env:"verbose,required"`
	XcodeCacheEnabled  bool `env:"xcode_cache_enabled,required"`
	GradleCacheEnabled bool `env:"gradle_cache_enabled,required"`
}

type Step struct {
	logger      Logger
	inputParser InputParser
	annotator   func(annotation service.Annotation) error
	command     Command
}

func New(
	logger Logger,
	inputParser InputParser,
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

	stepconf.Print(input)
	step.logger.Println()

	if !input.XcodeCacheEnabled && !input.GradleCacheEnabled {
		step.logger.Infof(NoFeaturesEnabledMsg)

		return nil
	}

	args := []string{"activate", "react-native"}
	if !input.GradleCacheEnabled {
		args = append(args, "--gradle=false", "--cpp=false")
	}
	if !input.XcodeCacheEnabled {
		args = append(args, "--xcode=false")
	}
	if input.Verbose {
		args = append(args, "--debug")
	}

	step.command.SetArgs(args)
	if err := step.command.Execute(); err != nil {
		return fmt.Errorf(FailedToActivateMsg+": %w", err)
	}

	step.logger.Infof(ReactNativeFeaturesActivatedMsg)

	return nil
}
