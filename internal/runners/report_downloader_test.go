package runners

import (
	"bytes"
	"errors"
	"github.com/aelmel/report-downloader/internal/report/store"
	mock_report "github.com/aelmel/report-downloader/mocks/report"
	mock_store "github.com/aelmel/report-downloader/mocks/report/store"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
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
	t.Run("success download", func(t *testing.T) {
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
	})
	t.Run("get report url error", func(t *testing.T) {
		reportClient.EXPECT().GetReport(gomock.Any(), "123456").Return("", "pending", errors.New("internal error"))
		_, err := d.getReportUrl("123456")
		assert.Equal(t, "internal error", err.Error())
	})
}

func Test_handle_report(t *testing.T) {
	log := logrus.New()
	log.Out = ioutil.Discard

	repo := store.NewReportStore(log)
	mockCtrl := gomock.NewController(t)
	reportClient := mock_report.NewMockClient(mockCtrl)
	defer mockCtrl.Finish()

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	rsp := &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString("1234,12345n\n1234,1245")),
	}

	reportClient.EXPECT().GetReport(gomock.Any(), "555666").Return("https://example.com/generated/1234567890.csv", "done", nil).
		Times(1)
	reportClient.EXPECT().DownloadReport("https://example.com/generated/1234567890.csv").Return(rsp, nil).
		Times(1).Do(func(_ interface{}) { wg.Done() })

	reportCh := store.NewReportChannel()

	d := downloader{
		reportClient: reportClient,
		logger:       log,
		repo:         repo,
		reportCh:     reportCh,
	}

	go d.handleReports()

	reportCh.AddReport("555666")
}

func Test_save_report(t *testing.T) {
	log := logrus.New()
	log.Out = ioutil.Discard

	mockCtrl := gomock.NewController(t)
	reportClient := mock_report.NewMockClient(mockCtrl)
	defer mockCtrl.Finish()
	d := downloader{
		reportClient: reportClient,
		logger:       log,
	}
	t.Run("success save report", func(t *testing.T) {
		rsp := &http.Response{
			Body: ioutil.NopCloser(bytes.NewBufferString(
				`1234,12345
			1234,1245`)),
		}

		reportClient.EXPECT().DownloadReport("https://example.com/generated/39d7f6d0c6a50.csv").Return(rsp, nil).Times(1)

		repFile, err := d.saveReport("https://example.com/generated/39d7f6d0c6a50.csv", "123456")
		assert.Nil(t, err)
		_, err = os.Stat(repFile)
		assert.Nil(t, err)
		err = os.Remove(repFile)
		assert.Nil(t, err)
	})
	t.Run("error downloading report", func(t *testing.T) {
		reportClient.EXPECT().DownloadReport("https://example.com/generated/39d7f6d0c6a50.csv").Return(nil, errors.New("internal error")).Times(1)
		_, err := d.saveReport("https://example.com/generated/39d7f6d0c6a50.csv", "123456")
		assert.NotNil(t, err)
	})
}
