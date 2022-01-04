package runners

import (
	"bytes"
	"github.com/aelmel/report-downloader/internal/report/store"
	mock_report "github.com/aelmel/report-downloader/mocks/report"
	mock_store "github.com/aelmel/report-downloader/mocks/report/store"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func Test_getReportUrl(t *testing.T) {
	log := logrus.New()
	log.Out = ioutil.Discard

	mockCtrl := gomock.NewController(t)
	reportCh := mock_store.NewMockChannel(mockCtrl)
	reportClient := mock_report.NewMockClient(mockCtrl)
	defer mockCtrl.Finish()

	d := downloader{
		reportClient: reportClient,
		logger:       log,
		reportCh:     reportCh,
	}

	pending := reportClient.EXPECT().GetReport(gomock.Any(), "123456").Return("", "pending", nil).Times(3)
	done := reportClient.EXPECT().GetReport(gomock.Any(), "123456").Return("https://example.com/generated/39d7f6d0c6a50.csv", "done", nil).
		Times(1)

	gomock.InOrder(
		pending,
		done,
	)
	url, err := d.getReportUrl("123456")
	assert.Nil(t, err)
	assert.Equal(t, "https://example.com/generated/39d7f6d0c6a50.csv", url)
}

func Test_handle_report(t *testing.T) {
	log := logrus.New()
	log.Out = ioutil.Discard

	repo := store.NewReportStore(log)
	mockCtrl := gomock.NewController(t)
	reportClient := mock_report.NewMockClient(mockCtrl)
	defer mockCtrl.Finish()

	rsp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString("1234,12345n\n1234,1245")),
	}

	reportClient.EXPECT().GetReport(gomock.Any(), "123456").Return("https://example.com/generated/39d7f6d0c6a50.csv", "done", nil).
		Times(1)
	reportClient.EXPECT().DownloadReport("https://example.com/generated/39d7f6d0c6a50.csv").Return(rsp, nil).Times(1)

	reportCh := store.NewReportChannel()

	d := downloader{
		reportClient: reportClient,
		logger:       log,
		repo:         repo,
		reportCh:     reportCh,
	}

	go d.handleReports()

	reportCh.AddReport("123456")
}
