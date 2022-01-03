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
	Close() error
}

type client struct {
	logger *logrus.Logger

	baseUrl *url.URL
	api     api.Client
}

func NewReportClient(logger *logrus.Logger, baseUrl string) (Client, error) {
	bUrl, err := url.Parse(baseUrl)
	if err != nil {
		logger.Errorf("error parsing baseUrl %s", err.Error())
		return nil, err
	}
	apiClient := api.NewClient(logger)
	return &client{
		logger:  logger,
		baseUrl: bUrl,
		api:     apiClient,
	}, err
}

func (c *client) GenerateReport(ctx context.Context) (reportId string, err error) {
	path := c.buildUrl("/generation/report-requests")
	req, err := http.NewRequest("POST", path, nil)
	if err != nil {
		c.logger.Warnf("could not create request %v", err)
		return reportId, err
	}

	req = req.WithContext(ctx)
	resp, err := c.api.SendRequest(req)
	if err != nil {
		c.logger.Warnf("Error received from api %s", err.Error())
		return reportId, err
	}

	return resp.ReportId, nil
}

func (c *client) GetReport(ctx context.Context, reportId string) (string, string, error) {
	path := c.buildUrl(fmt.Sprintf("generation/reports/%s", reportId))
	req, err := http.NewRequest("POST", path, nil)
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

	return resp.ReportUrl, resp.Status, nil
}

func (c *client) buildUrl(resource string) string {
	base := *c.baseUrl
	base.Path += resource
	return base.String()
}

func (c *client) Close() error {
	return c.api.Close()
}
