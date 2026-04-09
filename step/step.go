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
	input       Input
}

func New(
	logger Logger,
	inputParser InputParser,
	annotator func(annotation service.Annotation) error,
	command Command,
) *Step {
	return &Step{
		logger:      logger,
		inputParser: inputParser,
		annotator:   annotator,
		command:     command,
	}
}

func (s *Step) ProcessConfig() error {
	if err := s.inputParser.Parse(&s.input); err != nil {
		return fmt.Errorf(FailedToParseInputsMsg+": %w", err)
	}
	s.logger.EnableDebugLog(s.input.Verbose)

	stepconf.Print(s.input)
	s.logger.Println()

	return nil
}

func (s *Step) InstallDeps() error {
	return nil
}

func (s *Step) Run() error {
	if !s.input.XcodeCacheEnabled && !s.input.GradleCacheEnabled {
		s.logger.Infof(NoFeaturesEnabledMsg)

		return nil
	}

	args := []string{"activate", "react-native"}
	if !s.input.GradleCacheEnabled {
		args = append(args, "--gradle=false", "--cpp=false")
	}
	if !s.input.XcodeCacheEnabled {
		args = append(args, "--xcode=false")
	}
	if s.input.Verbose {
		args = append(args, "--debug")
	}

	s.command.SetArgs(args)
	if err := s.command.Execute(); err != nil {
		return fmt.Errorf(FailedToActivateMsg+": %w", err)
	}

	s.logger.Infof(ReactNativeFeaturesActivatedMsg)

	return nil
}

func (s *Step) ExportOutputs() error {
	return nil
}
