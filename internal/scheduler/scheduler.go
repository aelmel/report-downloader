package scheduler

import (
	"github.com/aelmel/report-downloader/internal/runners"
	"github.com/sirupsen/logrus"

	"github.com/robfig/cron/v3"
)

type Scheduler interface {
	Close() error
}

type scheduler struct {
	c      *cron.Cron
	logger *logrus.Logger
}

func NewScheduler(logger *logrus.Logger, frequency string, runners ...runners.Runner) Scheduler {
	c := cron.New()

	for _, runner := range runners {
		c.AddFunc(frequency, runner.Execute)
	}
	c.Start()

	return &scheduler{c: c, logger: logger}
}

func (s *scheduler) Close() error {
	s.c.Stop()
	return nil
}
