package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is the global logger instance.
var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.InfoLevel)
}
