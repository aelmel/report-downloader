// Code generated by MockGen. DO NOT EDIT.
// Source: internal/report/client.go

// Package mock_report is a generated GoMock package.
package mock_report

import (
	context "context"
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockClient)(nil).Close))
}

// DownloadReport mocks base method.
func (m *MockClient) DownloadReport(reportAddress string) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadReport", reportAddress)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownloadReport indicates an expected call of DownloadReport.
func (mr *MockClientMockRecorder) DownloadReport(reportAddress interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadReport", reflect.TypeOf((*MockClient)(nil).DownloadReport), reportAddress)
}

// GenerateReport mocks base method.
func (m *MockClient) GenerateReport(ctx context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateReport", ctx)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateReport indicates an expected call of GenerateReport.
func (mr *MockClientMockRecorder) GenerateReport(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateReport", reflect.TypeOf((*MockClient)(nil).GenerateReport), ctx)
}

// GetReport mocks base method.
func (m *MockClient) GetReport(ctx context.Context, reportId string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReport", ctx, reportId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetReport indicates an expected call of GetReport.
func (mr *MockClientMockRecorder) GetReport(ctx, reportId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReport", reflect.TypeOf((*MockClient)(nil).GetReport), ctx, reportId)
}
