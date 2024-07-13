package actions

import (
	"net/http"
	"strings"
)

func ProcessHttps(request *http.Request) *http.Request {
	if strings.Contains(request.RequestURI, "https://") {
		request.URL.Scheme = "https"
	}

	return request
}
