package runners

import (
	"context"
	"errors"
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
	reportCh     store.Channel
}

var backoffSchedule = []time.Duration{
	10 * time.Second,
	10 * time.Second,
	10 * time.Second,
	10 * time.Second,
	10 * time.Second,
	10 * time.Second,
}

func NewDownloader(store store.ReportStore, reportCh store.Channel, client report.Client, logger *logrus.Logger) Runner {
	d := &downloader{
		reportClient: client,
		reportCh:     reportCh,
		repo:         store,
		logger:       logger,
	}

	go d.handleReports()
	return d
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
			file, err := d.saveReport(url, id)
			if err != nil {
				d.logger.Warnf("couldn't save report for reportId %s", id)
				return
			}
			d.logger.Infof("saved file %s report id %s", file, id)
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
			return url, nil
		}
		d.logger.Infof("get report id %s go status %s", id, status)
		time.Sleep(backoff)
	}

	return "", errors.New("retries attempts exceeded")
}

func (d *downloader) saveReport(url, reportID string) (string, error) {
	resp, err := d.reportClient.DownloadReport(url)
	if err != nil {
		d.logger.Warnf("error %s downloading report %s", url, err.Error())
		return "", err
	}

	now := time.Now().Format("2006_01_02_15_04_05")
	fileLocation := fmt.Sprintf("%s/%s_%s.csv", os.TempDir(), now, reportID)
	out, err := os.Create(fileLocation)

	if err != nil {
		d.logger.Warnf("error on file create %s", err.Error())
		return "", err
	}
	defer out.Close()
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		d.logger.Warnf("error %s copying resp to file ", err.Error())
	}

	return fileLocation, nil
}

func (d *downloader) handleReports() {
	for reportId := range d.reportCh.GetReportChannel() {
		d.logger.Debugf("got report id %s", reportId)
		go func(reportId string) {
			url, err := d.getReportUrl(reportId)
			if err != nil {
				d.logger.Warnf("error geting report url %s", reportId)
				d.repo.AddReport(reportId)
				return
			}
			file, err := d.saveReport(url, reportId)
			d.logger.Infof("save report to file %s id %s", file, reportId)
		}(reportId)
	}
}
