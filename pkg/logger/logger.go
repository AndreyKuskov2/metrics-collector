package logger

import (
	"github.com/sirupsen/logrus"
)

// Логгер
var Log *logrus.Logger

// Функция создания нового логгера
func NewLogger() *logrus.Logger {
	if Log != nil {
		return Log
	}

	Log = logrus.New()

	return Log
}
