package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opencontainers/go-digest"
)

func TestBlobVerify(t *testing.T) {
	tt := []struct {
		name      string
		transport http.RoundTripper
		handler   http.HandlerFunc
		valid     bool
	}{
		{"golden", &DefaultTransport, testVerifyBlobHandler, true},
		{"bad round trip", &badTransport{}, testVerifyBlobHandler, false},
		{"blob not found", &DefaultTransport, testVerifyBlobNotFoundHandler, false},
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
			err = api.VerifyBlob("foo/bar", digest.Digest("sha256:4dda4f89dd2bec1791d7b41c8d63e33e8f7f19092aa3057f2c2be127c6e4b2b9"))
			if err == nil && !tc.valid {
				t.Fatal("expected invalid test-case")
			}
		}
		t.Run(tc.name, tf)
	}
}

func testVerifyBlobHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/v2/foo/bar/blobs/sha256:4dda4f89dd2bec1791d7b41c8d63e33e8f7f19092aa3057f2c2be127c6e4b2b9" {
		w.WriteHeader(http.StatusOK)
	}
}

func testVerifyBlobNotFoundHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/v2/foo/bar/blobs/sha256:4dda4f89dd2bec1791d7b41c8d63e33e8f7f19092aa3057f2c2be127c6e4b2b9" {
		w.WriteHeader(http.StatusNotFound)
	}
}
