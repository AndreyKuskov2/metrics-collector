package logger

import (
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func NewLogger(logPath string) *logrus.Logger {
	if Log != nil {
		return Log
	}

	Log = logrus.New()

	// file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err == nil {
	// 	Log.Out = file
	// } else {
	// 	Log.Info("Failed to log to file, using default stderr")
	// }

	return Log
}
