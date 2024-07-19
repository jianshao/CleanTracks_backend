package logs

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func InitLog() {
	path := "~/logs/log"
	logf, err := rotatelogs.New(
		path+".%Y%m%d",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithMaxAge(30*time.Hour*24))
	if err != nil {
		logrus.Fatal(err.Error())
		return
	}

	logger.Out = logf
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)
}

func WriteLog(level logrus.Level, fields map[string]any, message string) {
	switch level {
	case logrus.DebugLevel:
		logger.WithFields(fields).Debug(message)
	case logrus.WarnLevel:
		logger.WithFields(fields).Warn(message)
	case logrus.InfoLevel:
		logger.WithFields(fields).Info(message)
	case logrus.ErrorLevel:
		logger.WithFields(fields).Error(message)
	case logrus.FatalLevel:
		logger.WithFields(fields).Fatal(message)
	}
}
