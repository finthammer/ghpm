package github

import (
	"net/http"
)

// Authenticator allows to implement the different authentication methods
// supported by GitHub.
type Authenticator interface {
	// AddAuthentication adds the authentication to the HTTP request.
	AddAuthentication(req *http.Request) (*http.Request, error)
}

// basicAuthenticator authorizes via HTTP basic with username and password.
type basicAuthenticator struct {
	username string
	password string
}

// NewBasicAuthenticator returns an authenticator using username and password.
func NewBasicAuthenticator(username, password string) Authenticator {
	return &basicAuthenticator{
		username: username,
		password: password,
	}
}

func (a *basicAuthenticator) AddAuthentication(req *http.Request) (*http.Request, error) {
	return req, nil
}
