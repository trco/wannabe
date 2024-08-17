package curl

import (
	"net/http"
)

func SetHeaders(request *http.Request, headers []Header) {
	for _, v := range headers {
		request.Header.Set(v.Key, v.Value)
	}
}
