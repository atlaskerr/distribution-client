package client

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	testHost := "http://localhost:5050"
	tt := []struct {
		name, host string
		auth       Authenticator
		transport  http.RoundTripper
	}{
		{"basic auth", testHost, &BasicAuth{"user", "password"}, nil},
		{"token auth", testHost, &TokenAuth{"token"}, nil},
		{"custom transport", testHost, nil, &http.Transport{}},
	}

	for _, tc := range tt {
		tf := func(t *testing.T) {
			cfg := Config{
				Host:      tc.host,
				Auth:      tc.auth,
				Transport: tc.transport,
			}
			c, err := New(cfg)
			if err != nil {
				t.Fatal(err)
			}
			NewDistributionAPI(c)
		}
		t.Run(tc.name, tf)
	}
}

func TestAuthSchemes(t *testing.T) {
	tt := []struct {
		name string
		auth Authenticator
	}{
		{"basic auth", &BasicAuth{"user", "password"}},
		{"token auth", &TokenAuth{"token"}},
	}

	for _, tc := range tt {
		tf := func(t *testing.T) {
			req := new(http.Request)
			req.Header = make(http.Header)
			tc.auth.Set(req)

			val := req.Header.Get("Authorization")
			if val == "" {
				t.Fatal("authorization header not set")
			}
		}
		t.Run(tc.name, tf)
	}
}

type badTransport struct{}

func (t *badTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("bad round trip")
}

func TestRoundTrip(t *testing.T) {
	tt := []struct {
		name      string
		transport http.RoundTripper
		valid     bool
	}{
		{"valid round trip", &DefaultTransport, true},
		{"bad round trip", &badTransport{}, false},
	}

	for _, tc := range tt {
		tf := func(t *testing.T) {
			cfg := Config{
				Host:      "http://localhost",
				Transport: tc.transport,
			}
			c, err := New(cfg)
			if err != nil {
				t.Fatal(err)
			}

			req, _ := http.NewRequest("GET", "http://localhost", nil)
			_, err = c.Transport.RoundTrip(req)
			if err == nil && !tc.valid {
				t.Fatal("expected valid transport")
			}

		}
		t.Run(tc.name, tf)
	}
}
