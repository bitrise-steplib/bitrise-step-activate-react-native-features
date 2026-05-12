package step

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/bitrise-io/bitrise-plugins-annotations/service"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
)

const (
	FailedToParseInputsMsg          = "failed to parse inputs"
	FailedToInstallCLIMsg           = "failed to install CLI"
	NoFeaturesEnabledMsg            = "No features enabled"
	FailedToActivateMsg             = "failed to activate Build Cache for React Native"
	ReactNativeFeaturesActivatedMsg = "Build Cache for React Native activated successfully"
)

type Input struct {
	Verbose            bool `env:"verbose,required"`
	XcodeCacheEnabled  bool `env:"xcode_cache_enabled,required"`
	GradleCacheEnabled bool `env:"gradle_cache_enabled,required"`
}

type Step struct {
	logger        Logger
	inputParser   InputParser
	annotator     func(annotation service.Annotation) error
	input         Input
	cliBinaryPath string
}

func New(
	logger Logger,
	inputParser InputParser,
	annotator func(annotation service.Annotation) error,
) *Step {
	return &Step{
		logger:      logger,
		inputParser: inputParser,
		annotator:   annotator,
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
	path, err := installCLI(context.Background(), s.logger)
	if err != nil {
		return fmt.Errorf(FailedToInstallCLIMsg+": %w", err)
	}

	s.cliBinaryPath = path

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

	cmd := exec.Command(s.cliBinaryPath, args...) //nolint:gosec
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(FailedToActivateMsg+": %w", err)
	}

	s.logger.Infof(ReactNativeFeaturesActivatedMsg)

	return nil
}

// SetCLIBinaryPath overrides the CLI binary path. Used in tests.
func (s *Step) SetCLIBinaryPath(path string) {
	s.cliBinaryPath = path
}

func (s *Step) ExportOutputs() error {
	return nil
}
