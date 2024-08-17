package hash

import (
	"net/http"
)

func setHeaders(request *http.Request, headers []header) {
	for _, v := range headers {
		request.Header.Set(v.key, v.value)
	}
}
