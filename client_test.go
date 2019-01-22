package client

import (
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
