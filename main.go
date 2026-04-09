package main

import (
	"os"

	cli "github.com/bitrise-io/bitrise-build-cache-cli/cmd/common"
	_ "github.com/bitrise-io/bitrise-build-cache-cli/cmd/reactnative"
	"github.com/bitrise-io/bitrise-plugins-annotations/service"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/errorutil"
	"github.com/bitrise-io/go-utils/v2/exitcode"
	"github.com/bitrise-io/go-utils/v2/log"

	"github.com/bitrise-steplib/bitrise-step-activate-react-native-features/step"
)

type stepLogger struct {
	log.Logger
}

func (l stepLogger) FormattedErrorf(err error) {
	l.Errorf("%s", errorutil.FormattedError(err))
}

func main() {
	exitCode := run()
	os.Exit(int(exitCode))
}

func run() exitcode.ExitCode {
	logger, s := createStep()

	if err := s.ProcessConfig(); err != nil {
		logger.FormattedErrorf(err)
		return exitcode.Failure
	}

	if err := s.InstallDeps(); err != nil {
		logger.FormattedErrorf(err)
		return exitcode.Failure
	}

	if err := s.Run(); err != nil {
		logger.FormattedErrorf(err)
		return exitcode.Failure
	}

	if err := s.ExportOutputs(); err != nil {
		logger.FormattedErrorf(err)
		return exitcode.Failure
	}

	return exitcode.Success
}

func createStep() (stepLogger, step.Runner) {
	logger := stepLogger{log.NewLogger()}
	envRepo := env.NewRepository()
	inputParser := stepconf.NewInputParser(envRepo)

	return logger, step.New(
		logger,
		inputParser,
		service.Annotate,
		cli.RootCmd,
	)
}
