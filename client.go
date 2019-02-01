package client

import (
	"net/http"
	"net/url"
	"time"
)

// DefaultTransport is the optional transport clients may use. Requests to the
// registry timeout after one second.
var DefaultTransport = http.Transport{ResponseHeaderTimeout: time.Second}

// Client is an implementation of http.RoundTripper.
type Client struct {
	Host      *url.URL
	Transport http.RoundTripper
	Auth      Authenticator
}

// Config defines the parameters for Client configuration.
type Config struct {
	Host      string
	Transport http.RoundTripper
	Auth      Authenticator
}

// New takes a Config and returns a fully initialized Client.
func New(cfg Config) (*Client, error) {
	host, _ := url.Parse(cfg.Host)

	var transport http.RoundTripper
	if cfg.Transport != nil {
		transport = cfg.Transport
	} else {
		transport = &DefaultTransport
	}

	c := &Client{
		Host:      host,
		Transport: transport,
		Auth:      cfg.Auth,
	}
	return c, nil
}

// RoundTrip is the Client implementation of http.RoundTripper. Used to hook
// into an http.Request before being set to the server.
func (c *Client) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if c.Auth != nil {
		c.Auth.Set(req)
	}
	return c.Transport.RoundTrip(req)
}
