package step_test

import (
	"testing"

	"github.com/bitrise-io/bitrise-plugins-annotations/service"
	utilsMocks "github.com/bitrise-io/go-utils/v2/mocks"
	"github.com/bitrise-steplib/bitrise-step-activate-react-native-features/step"
	"github.com/bitrise-steplib/bitrise-step-activate-react-native-features/step/features"
	"github.com/bitrise-steplib/bitrise-step-activate-react-native-features/step/mocks"
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

		mockParser := mocks.NewMockInputParser(t)
		mockParser.On("Parse", mock.Anything).Run(func(args mock.Arguments) {
			switch v := args.Get(0).(type) {
			case *step.Input:
				v.Verbose = true
			case *features.CPPCacheInput:
				v.CPPCacheEnabled = true
			case *features.XcodeCacheInput:
				v.XcodeCacheEnabled = true
			case *features.GradleCacheInput:
				v.GradleCacheEnabled = true
			}
		}).Return(nil)

		mockCmd := mocks.NewMockCommand(t)
		mockCmd.On("SetArgs", mock.Anything).Return()
		mockCmd.On("Execute").Return(nil)

		sut := step.New(
			logger,
			mockParser,
			func(annotation service.Annotation) error { return nil },
			mockCmd,
		)

		err := sut.Run()

		assert.Nil(t, err)
	})

	t.Run("Failed to parse input", func(t *testing.T) {
		logger := &utilsMocks.Logger{}

		mockParser := mocks.NewMockInputParser(t)
		mockParser.On("Parse", mock.AnythingOfType("*step.Input")).Return(assert.AnError)

		mockCmd := mocks.NewMockCommand(t)

		sut := step.New(
			logger,
			mockParser,
			func(annotation service.Annotation) error { return nil },
			mockCmd,
		)

		err := sut.Run()
		assert.ErrorContains(t, err, step.FailedToParseInputsMsg)
	})

	t.Run("No features enabled", func(t *testing.T) {
		logger := &utilsMocks.Logger{}
		logger.On("EnableDebugLog", false).Return().Once()
		logger.On("Println", mock.Anything).Return().Once()
		logger.On("Debugf", mock.Anything, mock.Anything).Return()
		logger.On("Infof", step.NoFeaturesEnabledMsg).Return().Once()

		mockParser := mocks.NewMockInputParser(t)
		mockParser.On("Parse", mock.Anything).Return(nil)
		// All feature inputs default to false (not enabled), so all features return nil

		mockCmd := mocks.NewMockCommand(t)

		sut := step.New(
			logger,
			mockParser,
			func(annotation service.Annotation) error { return nil },
			mockCmd,
		)

		err := sut.Run()
		assert.Nil(t, err)
	})

	t.Run("Failed to activate", func(t *testing.T) {
		logger := &utilsMocks.Logger{}
		logger.On("EnableDebugLog", false).Return().Once()
		logger.On("Println", mock.Anything).Return().Once()
		logger.On("Debugf", mock.Anything, mock.Anything).Return()

		mockParser := mocks.NewMockInputParser(t)
		mockParser.On("Parse", mock.Anything).Run(func(args mock.Arguments) {
			if v, ok := args.Get(0).(*features.CPPCacheInput); ok {
				v.CPPCacheEnabled = true
			}
		}).Return(nil)

		mockCmd := mocks.NewMockCommand(t)
		mockCmd.On("SetArgs", mock.Anything).Return()
		mockCmd.On("Execute").Return(assert.AnError)

		sut := step.New(
			logger,
			mockParser,
			func(annotation service.Annotation) error { return nil },
			mockCmd,
		)

		err := sut.Run()
		assert.EqualError(t, err, step.FailedToActivateMsg+": "+assert.AnError.Error())
	})
}
