package actions

import (
	"net/http"
	"strings"
)

// sets scheme to https if request uri contains it but the scheme is http due to proxying
func ProcessHttps(request *http.Request) *http.Request {
	if strings.Contains(request.RequestURI, "https://") {
		request.URL.Scheme = "https"
	}

	return request
}
