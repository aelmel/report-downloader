package report

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aelmel/report-downloader/internal/api"

	"github.com/sirupsen/logrus"
)

type Client interface {
	GenerateReport(ctx context.Context) (reportId string, err error)
	GetReport(ctx context.Context, reportId string) (string, string, error)
	DownloadReport(reportAddress string) (*http.Response, error)
	Close() error
}

type client struct {
	logger *logrus.Logger

	baseURL *url.URL
	api     api.Client
}

func NewReportClient(logger *logrus.Logger, baseUrl string) (Client, error) {
	bUrl, err := url.Parse(baseUrl)
	if err != nil {
		logger.Errorf("error parsing baseURL %s", err.Error())
		return nil, err
	}
	apiClient := api.NewClient(logger)
	return &client{
		logger:  logger,
		baseURL: bUrl,
		api:     apiClient,
	}, err
}

func (c *client) GenerateReport(ctx context.Context) (reportID string, err error) {
	path := c.buildUrl("/generation/report-requests")
	req, err := http.NewRequest("POST", path, nil)
	if err != nil {
		c.logger.Warnf("could not create request %v", err)
		return reportID, err
	}

	req = req.WithContext(ctx)
	resp, err := c.api.SendRequest(req)
	if err != nil {
		c.logger.Warnf("Error received from api %s", err.Error())
		return reportID, err
	}

	return fmt.Sprintf("%d", resp.ReportID), nil
}

func (c *client) GetReport(ctx context.Context, reportID string) (string, string, error) {
	path := c.buildUrl(fmt.Sprintf("generation/reports/%s", reportID))
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		c.logger.Warnf("could not create request %v", err)
		return "", "", err
	}

	req = req.WithContext(ctx)
	resp, err := c.api.SendRequest(req)
	if err != nil {
		c.logger.Warnf("error received from api %s", err.Error())
		return "", "", err
	}

	return resp.ReportURL, resp.Status, nil
}

func (c *client) DownloadReport(reportAddress string) (*http.Response, error) {
	req, err := http.NewRequest("GET", reportAddress, nil)
	if err != nil {
		c.logger.Warnf("error %s creating request %s", reportAddress, err.Error())
		return nil, err
	}

	return c.api.SendDownloadRequest(req)
}

func (c *client) buildUrl(resource string) string {
	base := *c.baseURL
	base.Path += resource
	return base.String()
}

func (c *client) Close() error {
	return c.api.Close()
}
