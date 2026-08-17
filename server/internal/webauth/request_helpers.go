package webauth

import (
	"net/http"
	"strings"
)

func isPartialRequest(r *http.Request) bool {
	return isDatastarRequest(r)
}

func isDatastarRequest(r *http.Request) bool {
	return strings.EqualFold(r.Header.Get("Datastar-Request"), "true")
}
