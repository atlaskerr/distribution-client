package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

//func TestOCIVerify(t *testing.T) {
//	tt := []struct {
//		name      string
//		transport http.RoundTripper
//		handler   http.Handler
//	}{
//		{""},
//	}
//}

func TestOCIVerify1(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(verifyHandler))
	defer srv.Close()

	cfg := Config{
		Host: srv.URL,
		Auth: &TokenAuth{
			Token: "token",
		},
	}
	c, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	api := NewDistributionAPI(c)
	err = api.Verify()
	if err != nil {
		t.Fatal(err)
	}
}

func verifyHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/v2" {
		w.Header().Set(HeaderVersionCheck, "registry/2.0")
		w.WriteHeader(http.StatusOK)
	}
}
