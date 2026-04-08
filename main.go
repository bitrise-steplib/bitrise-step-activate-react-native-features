package main

import (
	"os"

	cli "github.com/bitrise-io/bitrise-build-cache-cli/cmd/common"
	_ "github.com/bitrise-io/bitrise-build-cache-cli/cmd/reactnative"
	"github.com/bitrise-io/bitrise-plugins-annotations/service"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/exitcode"

	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/bitrise-step-activate-react-native-features/step"
)

func main() {
	exitCode := run()
	os.Exit(int(exitCode))
}

func run() exitcode.ExitCode {
	logger := log.NewLogger()
	envRepo := env.NewRepository()
	inputParser := stepconf.NewInputParser(envRepo)

	stepInstance := step.New(
		logger,
		inputParser,
		service.Annotate,
		cli.RootCmd,
	)
	if err := stepInstance.Run(); err != nil {
		logger.Errorf(err.Error())
		return exitcode.Failure
	}

	return exitcode.Success
}
