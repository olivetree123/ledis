package common

import "github.com/sirupsen/logrus"

var Logger = logrus.New()

type Command struct {
	Code int
	Args interface{}
}

func init() {
	Logger.SetReportCaller(true)
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}
