package runners

import (
	"context"

	"github.com/aelmel/report-downloader/internal/report"
	"github.com/aelmel/report-downloader/internal/report/store"

	"github.com/sirupsen/logrus"
)

type generator struct {
	reportClient report.Client
	repo         store.ReportStore
	parallel     int
	logger       *logrus.Logger
}

func NewReportGenerator(store store.ReportStore, client report.Client, parallel int, logger *logrus.Logger) Runner {
	return &generator{
		reportClient: client,
		repo:         store,
		parallel:     parallel,
		logger:       logger,
	}
}

func (g *generator) Execute() {
	g.logger.Info("generate report")
	for i := 0; i < g.parallel; i++ {
		go func() {
			reportId, err := g.reportClient.GenerateReport(context.Background())
			if err != nil {
				g.logger.Warnf("error generating report")
				return
			}
			g.repo.AddReport(reportId)
		}()
	}
}
