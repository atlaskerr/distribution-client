package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifyManifest(t *testing.T) {
	tt := []struct {
		name      string
		transport http.RoundTripper
		handler   http.HandlerFunc
		valid     bool
	}{
		{"golden", &DefaultTransport, testVerifyManifestHandler, true},
		{"bad round trip", &badTransport{}, testVerifyManifestHandler, false},
		{"manifest not found", &DefaultTransport, testVerifyManifestNotFoundHandler, false},
	}

	for _, tc := range tt {
		tf := func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(tc.handler))
			defer srv.Close()

			cfg := Config{
				Host:      srv.URL,
				Transport: tc.transport,
				Auth:      &TokenAuth{"token"},
			}
			c, err := New(cfg)
			if err != nil {
				t.Fatal(err)
			}

			api := NewDistributionAPI(c)
			err = api.VerifyManifest("foo/bar", "latest")
			if err == nil && !tc.valid {
				t.Fatal("expected invalid test-case")
			}
		}
		t.Run(tc.name, tf)
	}
}

func testVerifyManifestHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/v2/foo/bar/manifests/latest" {
		w.WriteHeader(http.StatusOK)
	}
}

func testVerifyManifestNotFoundHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/v2/foo/bar/manifests/latest" {
		w.WriteHeader(http.StatusNotFound)
	}
}

func TestGetManifests(t *testing.T) {
	tt := []struct {
		name      string
		transport http.RoundTripper
		handler   http.HandlerFunc
		valid     bool
	}{}

	for _, tc := range tt {
		tf := func(t *testing.T) {
			h := http.HandlerFunc(tc.handler)
			srv := httptest.NewServer(h)
			defer srv.Close()

			cfg := Config{
				Host:      srv.URL,
				Transport: tc.transport,
				Auth:      &TokenAuth{"token"},
			}
			c, err := New(cfg)
			if err != nil {
				t.Fatal(err)
			}

			api := NewDistributionAPI(c)
			_, _, err = api.GetManifests("foo/bar", "latest")
			if err == nil && !tc.valid {
				t.Fatal("expected invalid test-case")
			}
		}
		t.Run(tc.name, tf)
	}
}
