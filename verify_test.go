package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOCIVerify(t *testing.T) {
	tt := []struct {
		name      string
		transport http.RoundTripper
		handler   http.HandlerFunc
		valid     bool
	}{
		{"golden", &DefaultTransport, testVerifyHandler, true},
		{"bad round trip", &badTransport{}, testVerifyHandler, false},
		{"verify failed", &DefaultTransport, testInvalidVerifyHandler, false},
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
			err = api.Verify()
			if err == nil && !tc.valid {
				t.Fatal("expected invalid test-case")
			}
		}
		t.Run(tc.name, tf)
	}
}

func testVerifyHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/v2" {
		w.Header().Set(HeaderVersionCheck, "registry/2.0")
		w.WriteHeader(http.StatusOK)
	}
}

func testInvalidVerifyHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/v2" {
		w.Header().Set(HeaderVersionCheck, "invalid")
		w.WriteHeader(http.StatusOK)
	}
}
