package helper

import (
	"github.com/merico-dev/lake/plugins/core"
	"github.com/sirupsen/logrus"
)

// bridge to current implementation at this point
// TODO: implement another TaskLogger for distributed runner/worker

func NewDefaultTaskLogger(log *logrus.Logger, prefix string, loggerPool map[string]*logrus.Logger) core.Logger {
	return NewDefaultLogger(log, prefix, loggerPool)
}
