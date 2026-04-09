package step_test

import (
	"testing"

	"github.com/bitrise-io/bitrise-plugins-annotations/service"
	"github.com/bitrise-steplib/bitrise-step-activate-react-native-features/step"
	"github.com/bitrise-steplib/bitrise-step-activate-react-native-features/step/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Step(t *testing.T) {
	t.Run("All features enabled", func(t *testing.T) {
		mockLogger := mocks.NewMockLogger(t)
		mockLogger.On("EnableDebugLog", true).Return().Once()
		mockLogger.On("Println").Return().Once()
		mockLogger.On("Infof", step.ReactNativeFeaturesActivatedMsg).Return().Once()

		mockParser := mocks.NewMockInputParser(t)
		mockParser.On("Parse", mock.Anything).Run(func(args mock.Arguments) {
			if v, ok := args.Get(0).(*step.Input); ok {
				v.Verbose = true
				v.XcodeCacheEnabled = true
				v.GradleCacheEnabled = true
			}
		}).Return(nil)

		mockCmd := mocks.NewMockCommand(t)
		mockCmd.On("SetArgs", []string{"activate", "react-native", "--debug"}).Return()
		mockCmd.On("Execute").Return(nil)

		sut := step.New(
			mockLogger,
			mockParser,
			func(annotation service.Annotation) error { return nil },
			mockCmd,
		)

		err := sut.ProcessConfig()
		assert.NoError(t, err)

		err = sut.Run()
		assert.NoError(t, err)
		mockCmd.AssertCalled(t, "SetArgs", []string{"activate", "react-native", "--debug"})
	})

	t.Run("Only gradle enabled — cpp follows gradle, xcode disabled", func(t *testing.T) {
		mockLogger := mocks.NewMockLogger(t)
		mockLogger.On("EnableDebugLog", false).Return().Once()
		mockLogger.On("Println").Return().Once()
		mockLogger.On("Infof", step.ReactNativeFeaturesActivatedMsg).Return().Once()

		mockParser := mocks.NewMockInputParser(t)
		mockParser.On("Parse", mock.Anything).Run(func(args mock.Arguments) {
			if v, ok := args.Get(0).(*step.Input); ok {
				v.GradleCacheEnabled = true
			}
		}).Return(nil)

		mockCmd := mocks.NewMockCommand(t)
		mockCmd.On("SetArgs", []string{"activate", "react-native", "--xcode=false"}).Return()
		mockCmd.On("Execute").Return(nil)

		sut := step.New(
			mockLogger,
			mockParser,
			func(annotation service.Annotation) error { return nil },
			mockCmd,
		)

		err := sut.ProcessConfig()
		assert.NoError(t, err)

		err = sut.Run()
		assert.NoError(t, err)
		mockCmd.AssertCalled(t, "SetArgs", []string{"activate", "react-native", "--xcode=false"})
	})

	t.Run("Only xcode enabled — gradle and cpp disabled", func(t *testing.T) {
		mockLogger := mocks.NewMockLogger(t)
		mockLogger.On("EnableDebugLog", false).Return().Once()
		mockLogger.On("Println").Return().Once()
		mockLogger.On("Infof", step.ReactNativeFeaturesActivatedMsg).Return().Once()

		mockParser := mocks.NewMockInputParser(t)
		mockParser.On("Parse", mock.Anything).Run(func(args mock.Arguments) {
			if v, ok := args.Get(0).(*step.Input); ok {
				v.XcodeCacheEnabled = true
			}
		}).Return(nil)

		mockCmd := mocks.NewMockCommand(t)
		mockCmd.On("SetArgs", []string{"activate", "react-native", "--gradle=false", "--cpp=false"}).Return()
		mockCmd.On("Execute").Return(nil)

		sut := step.New(
			mockLogger,
			mockParser,
			func(annotation service.Annotation) error { return nil },
			mockCmd,
		)

		err := sut.ProcessConfig()
		assert.NoError(t, err)

		err = sut.Run()
		assert.NoError(t, err)
		mockCmd.AssertCalled(t, "SetArgs", []string{"activate", "react-native", "--gradle=false", "--cpp=false"})
	})

	t.Run("Failed to parse input", func(t *testing.T) {
		mockLogger := mocks.NewMockLogger(t)

		mockParser := mocks.NewMockInputParser(t)
		mockParser.On("Parse", mock.AnythingOfType("*step.Input")).Return(assert.AnError)

		mockCmd := mocks.NewMockCommand(t)

		sut := step.New(
			mockLogger,
			mockParser,
			func(annotation service.Annotation) error { return nil },
			mockCmd,
		)

		err := sut.ProcessConfig()
		assert.ErrorContains(t, err, step.FailedToParseInputsMsg)
	})

	t.Run("No features enabled", func(t *testing.T) {
		mockLogger := mocks.NewMockLogger(t)
		mockLogger.On("EnableDebugLog", false).Return().Once()
		mockLogger.On("Println").Return().Once()
		mockLogger.On("Infof", step.NoFeaturesEnabledMsg).Return().Once()

		mockParser := mocks.NewMockInputParser(t)
		mockParser.On("Parse", mock.Anything).Return(nil)

		mockCmd := mocks.NewMockCommand(t)

		sut := step.New(
			mockLogger,
			mockParser,
			func(annotation service.Annotation) error { return nil },
			mockCmd,
		)

		err := sut.ProcessConfig()
		assert.NoError(t, err)

		err = sut.Run()
		assert.NoError(t, err)
	})

	t.Run("Failed to activate", func(t *testing.T) {
		mockLogger := mocks.NewMockLogger(t)
		mockLogger.On("EnableDebugLog", false).Return().Once()
		mockLogger.On("Println").Return().Once()

		mockParser := mocks.NewMockInputParser(t)
		mockParser.On("Parse", mock.Anything).Run(func(args mock.Arguments) {
			if v, ok := args.Get(0).(*step.Input); ok {
				v.GradleCacheEnabled = true
			}
		}).Return(nil)

		mockCmd := mocks.NewMockCommand(t)
		mockCmd.On("SetArgs", mock.Anything).Return()
		mockCmd.On("Execute").Return(assert.AnError)

		sut := step.New(
			mockLogger,
			mockParser,
			func(annotation service.Annotation) error { return nil },
			mockCmd,
		)

		err := sut.ProcessConfig()
		assert.NoError(t, err)

		err = sut.Run()
		assert.EqualError(t, err, step.FailedToActivateMsg+": "+assert.AnError.Error())
	})
}
