package logging

import (
	"github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"log/syslog"
)

func New(appName string) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	hook, err := lSyslog.NewSyslogHook("udp", "syslog:514", syslog.LOG_INFO, "")
	if err != nil {
		logger.Fatal(err)
	}
	logger.Hooks.Add(hook)
	return logger
}
