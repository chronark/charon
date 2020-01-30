package logging

import (
	"github.com/sirupsen/logrus"
)

func New(serviceName string) *logrus.Entry {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.JSONFormatter{})
	contextLogger := logger.WithFields(logrus.Fields{
		"service": serviceName,
	})
	return contextLogger
}
