package main

import (
	"io"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
)

func graceful(logger *logrus.Logger, done chan int, signals []os.Signal, services ...io.Closer) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, signals...)
	sig := <-sigc

	logger.Infof("Shutting down signal %v", sig)
	for _, item := range services {
		if err := item.Close(); err != nil {
			logger.Errorf("Failed to close %v: %v,", item, err)
		}
	}
	done <- 0
}
