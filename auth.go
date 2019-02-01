package client

import (
	"net/http"
	"strings"
)

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

// TokenAuth is an implementation of the Authenticator interface for token
// authentication.
type TokenAuth struct {
	Token string
}

// Set adds an authentication header to a request.
func (a *BasicAuth) Set(req *http.Request) {
	req.SetBasicAuth(a.Username, a.Password)
}

// Set add an authentication header to a request.
func (a *TokenAuth) Set(req *http.Request) {
	bearer := strings.Join([]string{"bearer", a.Token}, " ")
	req.Header.Set("Authorization", bearer)
}
