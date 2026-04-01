package step_test

import (
	"testing"

	"github.com/bitrise-io/bitrise-plugins-annotations/service"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	utilsMocks "github.com/bitrise-io/go-utils/v2/mocks"
	"github.com/bitrise-steplib/bitrise-step-activate-react-native-features/step"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Step(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		logger := &utilsMocks.Logger{}
		logger.On("EnableDebugLog", true).Return().Once()
		logger.On("Println", mock.Anything).Return().Once()
		logger.On("Debugf", mock.Anything, mock.Anything).Return()
		logger.On("Infof", step.ReactNativeFeaturesActivatedMsg).Return().Once()

		envRepo := NewMockEnvRepo()
		envRepo.Set("cpp_cache_enabled", "true")    //nolint: errcheck
		envRepo.Set("xcode_cache_enabled", "true")  //nolint: errcheck
		envRepo.Set("gradle_cache_enabled", "true") //nolint: errcheck
		envRepo.Set("verbose", "true")              //nolint: errcheck

		command := &MockCommand{}

		sut := step.New(
			logger,
			stepconf.NewInputParser(envRepo),
			func(annotation service.Annotation) error { return nil },
			command,
		)

		err := sut.Run()

		assert.Nil(t, err)
	})

	t.Run("Failed to parse input", func(t *testing.T) {
		logger := &utilsMocks.Logger{}
		command := &MockCommand{}
		envRepo := NewMockEnvRepo()

		sut := step.New(
			logger,
			stepconf.NewInputParser(envRepo),
			func(annotation service.Annotation) error { return nil },
			command,
		)

		err := sut.Run()
		assert.ErrorContains(t, err, step.FailedToParseInputsMsg)
		assert.Equal(t, 0, command.Executed)
	})

	t.Run("No features enabled", func(t *testing.T) {
		envRepo := NewMockEnvRepo()
		envRepo.Set("verbose", "false") //nolint: errcheck
		// cache feature inputs are not set, so all features return nil

		logger := &utilsMocks.Logger{}
		logger.On("EnableDebugLog", false).Return().Once()
		logger.On("Println", mock.Anything).Return().Once()
		logger.On("Debugf", mock.Anything, mock.Anything).Return()
		logger.On("Infof", step.NoFeaturesEnabledMsg).Return().Once()

		command := &MockCommand{}

		sut := step.New(
			logger,
			stepconf.NewInputParser(envRepo),
			func(annotation service.Annotation) error { return nil },
			command,
		)

		err := sut.Run()
		assert.Nil(t, err)
		assert.Equal(t, 0, command.Executed)
	})

	t.Run("Failed to activate", func(t *testing.T) {
		envRepo := NewMockEnvRepo()
		envRepo.Set("cpp_cache_enabled", "true") //nolint: errcheck
		envRepo.Set("verbose", "false")          //nolint: errcheck

		logger := &utilsMocks.Logger{}
		logger.On("EnableDebugLog", false).Return().Once()
		logger.On("Println", mock.Anything).Return().Once()
		logger.On("Debugf", mock.Anything, mock.Anything).Return()
		logger.On("Warnf", mock.Anything, mock.Anything).Return()

		command := &MockCommand{
			ExecutionError: assert.AnError,
		}

		sut := step.New(
			logger,
			stepconf.NewInputParser(envRepo),
			func(annotation service.Annotation) error { return nil },
			command,
		)

		err := sut.Run()
		assert.EqualError(t, err, step.FailedToActivateMsg+": "+assert.AnError.Error())
	})
}
