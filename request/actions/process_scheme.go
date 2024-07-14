package actions

import (
	"net/http"
	"strings"
)

// sets scheme to https if request uri contains it but the scheme set on url is http
func ProcessScheme(request *http.Request) *http.Request {
	if strings.Contains(request.RequestURI, "https://") {
		request.URL.Scheme = "https"
	}

	return request
}
