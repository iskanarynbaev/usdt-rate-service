package mocks

import (
	"net/http"
)

type HTTPClientMock struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *HTTPClientMock) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}
