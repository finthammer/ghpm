package infra

import (
	"encoding/json"
	"net/http"
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

// ReplyJSON marshals the passed reply into JSON and sends it via
// the response writer.
func ReplyJSON(w http.ResponseWriter, reply interface{}) {
	b, err := json.Marshal(reply)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
