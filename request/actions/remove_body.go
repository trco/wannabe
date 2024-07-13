package actions

import "net/http"

// prevents sending body and related headers in GET requests
func RemoveBody(request *http.Request) *http.Request {
	isGet := request.Method == "GET"
	if isGet {
		request.Body = http.NoBody
		request.ContentLength = 0
		request.Header.Del("Content-Length")
		request.Header.Del("Content-Type")
	}

	return request
}
