package infra

import (
	"net/http"
	"path"
	"strings"
)

// PathAt returns the nth part of the given request path and true if
// it exists. Otherwise an empty string and false. This way users of
// the nested handlers can retrieve an entity ID out of the path.
func PathAt(p string, n int) (string, bool) {
	if n < 0 {
		panic("illegal path index")
	}
	parts := strings.Split(p[1:], "/")
	if len(parts) < n+1 || parts[n] == "" {
		return "", false
	}
	return parts[n], true
}

// NestedHandler allows to nest handler following the
// pattern /handlerA/{entityID-A}/handlerB/{entityID-B}.
type NestedHandler struct {
	handlers []http.Handler
}

// NewNestedHandler creates an empty nested handler.
func NewNestedHandler() *NestedHandler {
	return &NestedHandler{}
}

// AppendHandler adds one handler to the stack of handlers.
func (nh *NestedHandler) AppendHandler(h http.Handler) {
	nh.handlers = append(nh.handlers, h)
}

// AppendHandlerFunc adds one handler function to the stack of handlers.
func (nh *NestedHandler) AppendHandlerFunc(
	hf func(http.ResponseWriter, *http.Request),
) {
	nh.handlers = append(nh.handlers, http.HandlerFunc(hf))
}

// ServeHTTP implements http.Handler.
func (nh *NestedHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	handler, ok := nh.handler(r)
	if !ok {
		http.Error(
			w,
			"cannot handle request",
			http.StatusRequestURITooLong,
		)
		return
	}
	handler.ServeHTTP(w, r)
}

// handler retrieves the correct handler from the stack.
func (nh *NestedHandler) handler(r *http.Request) (http.Handler, bool) {
	path := cleanPath(r.URL.Path)
	plen := len(strings.Split(path, "/"))
	prest := plen % 2
	index := plen/2 + prest
	if index > len(nh.handlers) {
		return nil, false
	}
	return nh.handlers[index-1], true
}

// cleanPath returns a canonical path. Borrowed from standard library
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		// Fast path for common case of p being the string we want:
		if len(p) == len(np)+1 && strings.HasPrefix(p, np) {
			np = p
		} else {
			np += "/"
		}
	}
	return np
}
