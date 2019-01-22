package client

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
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
	BaseEndpoint *url.URL
	Transport    http.RoundTripper
	Auth         Authenticator
}

// Config defines the parameters for Client configuration.
type Config struct {
	BaseEndpoint string
	Transport    http.RoundTripper
	Auth         Authenticator
}

// Authenticator is the interface all auth methods must satisfy.
type Authenticator interface {
	Set(*http.Request)
}

// BasicAuth is an implementation of the Authenticator interface for
// username/password authentication.
type BasicAuth struct {
	Username string
	Password string
}

// Set adds an authentication header to a request.
func (a *BasicAuth) Set(req *http.Request) {
	req.SetBasicAuth(a.Username, a.Password)
}

// TokenAuth is an implementation of the Authenticator interface for token
// authentication.
type TokenAuth struct {
	Token string
}

// Set add an authentication header to a request.
func (a *TokenAuth) Set(req *http.Request) {
	bearer := strings.Join([]string{"bearer", a.Token}, " ")
	req.Header.Set("Authorization", bearer)
}

// New takes a Config and returns a fully initialized Client.
func New(cfg Config) (*Client, error) {
	host, err := url.Parse(cfg.BaseEndpoint)
	if err != nil {
		return nil, ErrURLInvalid
	}

	c := &Client{
		BaseEndpoint: host,
		Transport:    cfg.Transport,
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

// NewDistributionAPI returns a fully initialized API for interacting with
// a remote registry.
func NewDistributionAPI(c *Client) *DistributionAPI {
	return &DistributionAPI{client: c}
}

// DistributionAPI contains methods for interacting with a remote registry.
type DistributionAPI struct {
	client *Client
}
