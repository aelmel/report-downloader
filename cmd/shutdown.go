package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func graceful(logger *logrus.Logger, done chan int, signals []os.Signal) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, signals...)
	sig := <-sigc

	logger.Infof("Shutting down signal %v", sig)

	done <- 0
}
