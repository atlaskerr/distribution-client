package client

import (
	"errors"
	"net/http"
	"net/url"
	"time"
)

var (
	// ErrURLInvalid is returned when the supplied endpoint URL is invalid.
	ErrURLInvalid = errors.New("client: Invalid endpoint URL")
)

// DefaultTransport is the optional transport clients may use. Requests to the
// registry timeout after one second.
var DefaultTransport = http.Transport{ResponseHeaderTimeout: time.Second}

// Client is an implementation of http.RoundTripper.
type Client struct {
	Endpoint  *url.URL
	Transport http.RoundTripper
}

// Config defines the parameters for Client configuration.
type Config struct {
	Endpoint  string
	Transport http.RoundTripper
}

// New takes a Config and returns a fully initialized Client.
func New(cfg Config) (*Client, error) {
	host, err := url.Parse(cfg.Endpoint)
	if err != nil {
		return nil, ErrURLInvalid
	}

	c := &Client{
		Endpoint:  host,
		Transport: cfg.Transport,
	}
	return c, nil
}

// Distribution defines the methods available for interacting with an
// OCI-compliant registry.
type Distribution interface {
	Verify() (bool, error)
}

