package logger

import (
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func NewLogger() *logrus.Logger {
	if Log != nil {
		return Log
	}

	Log = logrus.New()

	return Log
}
