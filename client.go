package client

import (
	"net/http"
	"time"
)

// DefaultTransport is the optional transport clients may use. Requests to the
// registry timeout after one second.
var DefaultTransport = http.Transport{ResponseHeaderTimeout: time.Second}

// Client is an implementation of http.RoundTripper.
type Client struct {
	Transport http.RoundTripper
}

// Config defines the parameters for Client configuration.
type Config struct {
	Transport http.RoundTripper
}

// New takes a Config and returns a fully initialized Client.
func New(cfg Config) (*Client, error) {
	var transport http.RoundTripper
	if cfg.Transport != nil {
		transport = cfg.Transport
	} else {
		transport = &DefaultTransport
	}

	c := &Client{
		Transport: transport,
	}
	return c, nil
}

// RoundTrip is the Client implementation of http.RoundTripper. Used to hook
// into an http.Request before being set to the server.
func (c *Client) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	return c.Transport.RoundTrip(req)
}
