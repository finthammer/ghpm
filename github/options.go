package github

import (
	"net/http"
)

// Options contains all options of the different GitHub API types.
type Options struct {
	Authentication Authentication
}

// Option is a function type modifying the passed options.
type Option func(os *Options)

// Authentication allows to implement the different authentication methods
// supported by GitHub.
type Authentication interface {
	// Apply adds the authentication to the HTTP request.
	Apply(req *http.Request) error
}

// basicAuth realizes the HTTP basic authentication.
type basicAuth struct {
	username string
	password string
}

// Apply implements Authentication.
func (a *basicAuth) Apply(req *http.Request) error {
	req.SetBasicAuth(a.username, a.password)
	return nil
}

// BasicAuth creates the basic authentication as Option.
func BasicAuth(username, password string) Option {
	return func(os *Options) {
		os.Authentication = &basicAuth{
			username: username,
			password: password,
		}
	}
}

// oauth2Token realizes the OAuth 2 authentication.
type oauth2Token string

// Apply implements Authentication.
func (t oauth2Token) Apply(req *http.Request) error {
	req.Header.Set("Authorization", "token "+string(t))
	return nil
}

// OAuth2Token creates the OAuth 2 authentication as Option.
func OAuth2Token(token string) Option {
	return func(os *Options) {
		os.Authentication = oauth2Token(token)
	}
}
