package main

import "github.com/sirupsen/logrus"

func main() {
	done := make(chan int, 1)

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	logger.Info("Start scheduler")

	<-done
}
