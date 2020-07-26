package mocks

import (
	"net/http"
)

var (
	// GetDoFunc returns the mock client's do function
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

// MockClient implements a mock of the HTTPClient interface
type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do allows MockClient to satisfy the HTTPClient interface
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}
