package client

import (
	"fmt"
	"net/http"
	"time"
)

// retryRoundTripper
// http.RoundTrip client middleware
// most of request retries will be handled in strategies
// serves as a general retry between client and webdriver
type retry struct {
	next       http.RoundTripper
	maxRetries int
	delay      time.Duration
}

// RoundTrip
// middleware for retries
func (rr retry) RoundTrip(r *http.Request) (*http.Response, error) {
	res, err := rr.next.RoundTrip(r)
	if err != nil {
		return res, err
	}

	return res, nil
}

// loggingRoundTripper
type loggin struct {
	next http.RoundTripper
}

// RountTrip
// middleware logger for Client
func (l loggin) RoundTrip(r *http.Request) (*http.Response, error) {
	res, err := l.next.RoundTrip(r)
	if err != nil {
		return nil, fmt.Errorf("error on %v request: %v", r, err)
	}

	return res, nil
}
