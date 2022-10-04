package commands

import (
	"io"
	"net/http"
	"time"
)

type TimedHttpClient interface {
	Get(url string) (*http.Response, time.Duration, error)
	Post(url string, contentType string, body io.Reader) (*http.Response, time.Duration, error)
}

type timedHTTPClient struct {
	*http.Client
}

func (client timedHTTPClient) Get(url string) (*http.Response, time.Duration, error) {
	startTime := time.Now()
	res, err := client.Client.Get(url)
	return res, time.Since(startTime), err
}

func (client timedHTTPClient) Post(url string, contentType string, body io.Reader) (*http.Response, time.Duration, error) {
	startTime := time.Now()
	res, err := client.Client.Post(url, contentType, body)
	return res, time.Since(startTime), err
}
