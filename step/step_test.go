package step_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/bitrise-plugins-annotations/service"
	"github.com/bitrise-steplib/bitrise-step-activate-react-native-features/step"
	"github.com/bitrise-steplib/bitrise-step-activate-react-native-features/step/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// fakeBinary creates a shell script at a temp path that captures its arguments
// to a file and exits with the given code. Returns the binary path and args file path.
func fakeBinary(t *testing.T, exitCode int) (binaryPath, argsFile string) {
	t.Helper()

	dir := t.TempDir()
	binaryPath = filepath.Join(dir, "bitrise-build-cache")
	argsFile = filepath.Join(dir, "captured-args")

	script := `#!/bin/sh
echo "$@" > ` + argsFile + `
exit ` + fmt.Sprintf("%d", exitCode)

	require.NoError(t, os.WriteFile(binaryPath, []byte(script), 0o755))

	return binaryPath, argsFile
}

func newTestStep(t *testing.T, input step.Input, binaryPath string) *step.Step {
	t.Helper()

	mockLogger := mocks.NewMockLogger(t)
	mockLogger.On("EnableDebugLog", mock.Anything).Return().Maybe()
	mockLogger.On("Println").Return().Maybe()
	mockLogger.On("Infof", mock.Anything).Return().Maybe()
	mockLogger.On("Infof", mock.Anything, mock.Anything).Return().Maybe()

	mockParser := mocks.NewMockInputParser(t)
	mockParser.On("Parse", mock.Anything).Run(func(args mock.Arguments) {
		if v, ok := args.Get(0).(*step.Input); ok {
			*v = input
		}
	}).Return(nil)

	s := step.New(
		mockLogger,
		mockParser,
		func(_ service.Annotation) error { return nil },
	)

	require.NoError(t, s.ProcessConfig())
	s.SetCLIBinaryPath(binaryPath)

	return s
}

func Test_Step(t *testing.T) {
	t.Run("All features enabled — passes correct flags", func(t *testing.T) {
		binaryPath, argsFile := fakeBinary(t, 0)
		s := newTestStep(t, step.Input{
			Verbose:            true,
			XcodeCacheEnabled:  true,
			GradleCacheEnabled: true,
		}, binaryPath)

		err := s.Run()

		require.NoError(t, err)
		args, _ := os.ReadFile(argsFile)
		assert.Contains(t, string(args), "activate react-native")
		assert.Contains(t, string(args), "--debug")
	})

	t.Run("Only gradle enabled — xcode disabled", func(t *testing.T) {
		binaryPath, argsFile := fakeBinary(t, 0)
		s := newTestStep(t, step.Input{
			GradleCacheEnabled: true,
		}, binaryPath)

		err := s.Run()

		require.NoError(t, err)
		args, _ := os.ReadFile(argsFile)
		assert.Contains(t, string(args), "--xcode=false")
		assert.NotContains(t, string(args), "--gradle=false")
	})

	t.Run("Only xcode enabled — gradle and cpp disabled", func(t *testing.T) {
		binaryPath, argsFile := fakeBinary(t, 0)
		s := newTestStep(t, step.Input{
			XcodeCacheEnabled: true,
		}, binaryPath)

		err := s.Run()

		require.NoError(t, err)
		args, _ := os.ReadFile(argsFile)
		assert.Contains(t, string(args), "--gradle=false")
		assert.Contains(t, string(args), "--cpp=false")
		assert.NotContains(t, string(args), "--xcode=false")
	})

	t.Run("No features enabled — does not call CLI", func(t *testing.T) {
		binaryPath, argsFile := fakeBinary(t, 0)
		s := newTestStep(t, step.Input{}, binaryPath)

		err := s.Run()

		require.NoError(t, err)
		_, err = os.ReadFile(argsFile)
		assert.True(t, os.IsNotExist(err), "CLI should not have been called")
	})

	t.Run("CLI failure is propagated", func(t *testing.T) {
		binaryPath, _ := fakeBinary(t, 1)
		s := newTestStep(t, step.Input{
			GradleCacheEnabled: true,
		}, binaryPath)

		err := s.Run()

		assert.ErrorContains(t, err, step.FailedToActivateMsg)
	})

	t.Run("Failed to parse input", func(t *testing.T) {
		mockLogger := mocks.NewMockLogger(t)
		mockParser := mocks.NewMockInputParser(t)
		mockParser.On("Parse", mock.AnythingOfType("*step.Input")).Return(assert.AnError)

		s := step.New(
			mockLogger,
			mockParser,
			func(_ service.Annotation) error { return nil },
		)

		err := s.ProcessConfig()
		assert.ErrorContains(t, err, step.FailedToParseInputsMsg)
	})
}
