package runners

import (
	"fmt"
	"github.com/aelmel/report-downloader/internal/report"
	mock_store "github.com/aelmel/report-downloader/mocks/report/store"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestGenerator_Execute(t *testing.T) {
	parallel := 4
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/generation/report-requests", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{
			"status":"done", 
			"reportId":99999
		}`))
	})

	log := logrus.New()
	log.Out = ioutil.Discard

	mockCtrl := gomock.NewController(t)
	reportCh := mock_store.NewMockChannel(mockCtrl)

	defer mockCtrl.Finish()
	reportClient, err := report.NewReportClient(log, server.URL)
	assert.Nil(t, err)

	var wg sync.WaitGroup
	wg.Add(parallel)
	defer wg.Wait()

	reportCh.EXPECT().AddReport(fmt.Sprintf("%d", 99999)).Times(parallel).Do(func(reportId interface{}) { wg.Done() })

	gen := generator{
		reportClient: reportClient,
		reportCh:     reportCh,
		parallel:     parallel,
		logger:       log,
	}

	gen.Execute()
}
