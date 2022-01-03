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
	c *cron.Cron
}

func NewScheduler(logger *logrus.Logger, runners ...runners.Runner) Scheduler {
	c := cron.New()

	cron.WithLogger(cron.VerbosePrintfLogger(logger))

	for _, runner := range runners {
		c.AddFunc("*/1 * * * *", runner.Execute)
	}

	c.Start()

	return &scheduler{c: c}
}

func (s *scheduler) Close() error {
	s.c.Stop()
	return nil
}
