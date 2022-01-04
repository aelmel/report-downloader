package report

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/aelmel/report-downloader/internal/api"
	mock_api "github.com/aelmel/report-downloader/mocks/api"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*client, *mock_api.MockClient, func()) {
	t.Helper()
	mockCtrl := gomock.NewController(t)
	apiClient := mock_api.NewMockClient(mockCtrl)

	bUrl, err := url.Parse("https://example.com")
	assert.Nil(t, err)

	log := logrus.New()
	log.Out = ioutil.Discard

	reportClient := &client{
		logger:  log,
		baseURL: bUrl,
		api:     apiClient,
	}

	close := func() { mockCtrl.Finish() }
	return reportClient, apiClient, close
}

func TestClient_GenerateReport_SuccessRequest(t *testing.T) {
	repClient, apiClient, close := setup(t)
	defer close()

	request, err := http.NewRequest("POST", "https://example.com/generation/report-requests", nil)
	assert.Nil(t, err)

	reportId := 99999
	resp := api.Response{
		Status:    "done",
		ReportURL: "",
		ReportID:  reportId,
		Error:     "",
	}
	apiClient.EXPECT().SendRequest(gomock.Eq(request)).Return(resp, nil).Times(1)

	id, err := repClient.GenerateReport(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, id, fmt.Sprintf("%d", reportId))
}

func TestClient_GenerateReport(t *testing.T) {
	repClient, apiClient, close := setup(t)
	defer close()

	request, err := http.NewRequest("POST", "https://example.com/generation/report-requests", nil)
	assert.Nil(t, err)

	apiClient.EXPECT().SendRequest(gomock.Eq(request)).Return(api.Response{}, errors.New("Something Went Wrong")).Times(1)

	id, err := repClient.GenerateReport(context.Background())
	assert.Empty(t, id)
	assert.Equal(t, "Something Went Wrong", err.Error())
}

func TestClient_GetReport(t *testing.T) {
	repClient, apiClient, close := setup(t)
	defer close()

	reportId := 12345

	request, err := http.NewRequest("GET", fmt.Sprintf("https://example.com/generation/reports/%d", reportId), nil)
	assert.Nil(t, err)
	resp := api.Response{
		Status:    "pending",
		ReportURL: "",
		ReportID:  reportId,
		Error:     "",
	}
	apiClient.EXPECT().SendRequest(gomock.Eq(request)).Return(resp, nil).Times(1)

	_, status, err := repClient.GetReport(context.Background(), fmt.Sprintf("%d", reportId))

	assert.Nil(t, err)
	assert.Equal(t, status, "pending")
}
