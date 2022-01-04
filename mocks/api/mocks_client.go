// Code generated by MockGen. DO NOT EDIT.
// Source: internal/api/client.go

// Package mock_api is a generated GoMock package.
package mock_api

import (
	http "net/http"
	reflect "reflect"

	api "github.com/aelmel/report-downloader/internal/api"
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

// SendDownloadRequest mocks base method.
func (m *MockClient) SendDownloadRequest(req *http.Request) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendDownloadRequest", req)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendDownloadRequest indicates an expected call of SendDownloadRequest.
func (mr *MockClientMockRecorder) SendDownloadRequest(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendDownloadRequest", reflect.TypeOf((*MockClient)(nil).SendDownloadRequest), req)
}

// SendRequest mocks base method.
func (m *MockClient) SendRequest(req *http.Request) (api.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendRequest", req)
	ret0, _ := ret[0].(api.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendRequest indicates an expected call of SendRequest.
func (mr *MockClientMockRecorder) SendRequest(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendRequest", reflect.TypeOf((*MockClient)(nil).SendRequest), req)
}