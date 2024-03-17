package dvlutil

import (
	"crypto/tls"
	"net/http"
)

// HTTPClient interface for sending HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// mockHTTPClient is a mock implementation of HTTPClient for testing
type mockHTTPClient struct {
	Resp *http.Response
	Err  error
}

// Do sends an HTTP request using the mock implementation
func (c mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.Resp, c.Err
}

// HttpConfig provides configuration options for Http Client
// for Rest APIs
type HttpConfig struct {
	Client HTTPClient
}

func DefaultHttpConfig() *HttpConfig {
	return &HttpConfig{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

// HttpOption is a functional option for configuring HttpConfig
type HttpOption func(*HttpConfig)

// WithHTTPClient sets the HTTP client for GetAuthToken
func WithHTTPClient(client HTTPClient) HttpOption {
	return func(config *HttpConfig) {
		config.Client = client
	}
}
