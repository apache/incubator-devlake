package unithelper

import (
	"github.com/apache/incubator-devlake/mocks"
	"github.com/stretchr/testify/mock"
)

func DummyLogger() *mocks.Logger {
	logger := new(mocks.Logger)
	logger.On("IsLevelEnabled", mock.Anything).Return(false).Maybe()
	logger.On("Printf", mock.Anything, mock.Anything).Maybe()
	logger.On("Log", mock.Anything, mock.Anything, mock.Anything).Maybe()
	logger.On("Debug", mock.Anything, mock.Anything).Maybe()
	logger.On("Info", mock.Anything, mock.Anything).Maybe()
	logger.On("Warn", mock.Anything, mock.Anything).Maybe()
	logger.On("Error", mock.Anything, mock.Anything).Maybe()
	return logger
}
