package infra

import (
	"net/http"
	"strings"
)

// NestedHandler allows to nest handler following the
// pattern /handlerA/{entityID-A}/handlerB/{entityID-B}.
type NestedHandler struct {
	handlerIDs  []string
	handlers    []http.Handler
	handlersLen int
}

// NewNestedHandler creates an empty nested handler.
func NewNestedHandler() *NestedHandler {
	return &NestedHandler{}
}

// AppendHandler adds one handler to the stack of handlers.
func (nh *NestedHandler) AppendHandler(id string, h http.Handler) {
	nh.handlerIDs = append(nh.handlerIDs, id)
	nh.handlers = append(nh.handlers, h)
	nh.handlersLen++
}

// AppendHandlerFunc adds one handler function to the stack of handlers.
func (nh *NestedHandler) AppendHandlerFunc(
	id string,
	hf func(http.ResponseWriter, *http.Request),
) {
	nh.AppendHandler(id, http.HandlerFunc(hf))
}

// ServeHTTP implements http.Handler.
func (nh *NestedHandler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	handler, ok := nh.handler(r.URL.Path)
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
func (nh *NestedHandler) handler(path string) (http.Handler, bool) {
	if strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}
	fields := strings.Split(path, "/")
	fieldsLen := len(fields)
	index := (fieldsLen - 1) / 2
	if index > nh.handlersLen-1 {
		return nil, false
	}
	if nh.handlerIDs[index] != fields[index*2] {
		return nil, false
	}
	return nh.handlers[index], true
}
