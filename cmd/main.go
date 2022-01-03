package main

import (
	"os"
	"syscall"

	"github.com/aelmel/report-downloader/internal/report"
	"github.com/aelmel/report-downloader/internal/report/store"
	"github.com/aelmel/report-downloader/internal/runners"
	"github.com/aelmel/report-downloader/internal/scheduler"

	"github.com/namsral/flag"
	"github.com/sirupsen/logrus"
)

func main() {
	var (
		baseUrl         = flag.String("API_URL", "http://localhost:12345", "api host")
		parallelClients = flag.Int("API_CLIENTS", 2, "parallel clients for report generator")
		frequency       = flag.String("GENERATE_FREQUENCY", "*/1 * * * *", "frequency calling report generator")
	)

	flag.Parse()

	done := make(chan int, 1)

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	reportCli, err := report.NewReportClient(logger, *baseUrl)

	if err != nil {
		logger.Error("error init report client")
		done <- 1
	}

	repo := store.NewReportStore(logger)
	logger.Info("start scheduler")

	generator := runners.NewReportGenerator(repo, reportCli, *parallelClients, logger)
	downloader := runners.NewDownloader(repo, reportCli, logger)
	scheduler := scheduler.NewScheduler(logger, *frequency, generator, downloader)

	go graceful(logger, done, []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt}, reportCli, scheduler)
	<-done
}
