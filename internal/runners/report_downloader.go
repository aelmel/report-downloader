package runners

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aelmel/report-downloader/internal/report"
	"github.com/aelmel/report-downloader/internal/report/store"

	"github.com/sirupsen/logrus"
)

type downloader struct {
	reportClient report.Client
	repo         store.ReportStore
	logger       *logrus.Logger
}

var backoffSchedule = []time.Duration{
	10 * time.Second,
	10 * time.Second,
	10 * time.Second,
	10 * time.Second,
	10 * time.Second,
	10 * time.Second,
}

func NewDownloader(store store.ReportStore, client report.Client, logger *logrus.Logger) Runner {
	return &downloader{
		reportClient: client,
		repo:         store,
		logger:       logger,
	}
}

func (d *downloader) Execute() {
	reports := d.repo.GetReportsId()
	d.logger.Infof("check existing reports found %d", len(reports))
	for _, id := range reports {
		go func(id string) {
			defer d.repo.RemoveReport(id)
			url, err := d.getReportUrl(id)
			if err != nil {
				d.logger.Warnf("could not get report address %s", err.Error())
				return
			}
			_ = d.saveReport(url, id)
		}(id)
	}
}

func (d *downloader) getReportUrl(id string) (url string, err error) {
	for _, backoff := range backoffSchedule {
		url, status, err := d.reportClient.GetReport(context.Background(), id)
		if err != nil {
			d.logger.Warnf("error %s recieved for id %s", id, err.Error())
			return url, err
		}
		if status == "done" {
			break
		}
		time.Sleep(backoff)
	}

	return url, nil
}

func (d *downloader) saveReport(url, reportId string) error {
	resp, err := d.reportClient.DownloadReport(url)
	if err != nil {
		d.logger.Warnf("error %s downloading report %s", url, err.Error())
		return err
	}

	now := time.Now().Format("2006_01_02_15_04_05")
	out, err := os.Create(fmt.Sprintf("%s/%s_%s", os.TempDir(), now, reportId))

	if err != nil {
		d.logger.Warnf("error on file create %s", err.Error())
		return err
	}
	defer out.Close()
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		d.logger.Warnf("error %s copying resp to file ", err.Error())
	}

	return nil
}
