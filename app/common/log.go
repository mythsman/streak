package common

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func InitLogger() {
	logrus.SetOutput(os.Stdout)
	level, err := logrus.ParseLevel(viper.GetString("logger.level"))
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableQuote:    true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logrus.Infoln("Logger init success")
}
