package logger

import "github.com/sirupsen/logrus"

func Error(message ...interface{}) {
	logrus.Error(message...)
}

func Errorf(format string, message ...interface{}) {
	logrus.Errorf(format, message...)
}
