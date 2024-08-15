package actions

import (
	"net/http"

	"github.com/trco/wannabe/types"
)

func SetHeaders(request *http.Request, headers []types.Header) {
	for _, v := range headers {
		request.Header.Set(v.Key, v.Value)
	}
}
