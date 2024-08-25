package hash

import (
	"fmt"
	"net/http"

	"moul.io/http2curl"
)

func prepareCurl(request *http.Request) (string, error) {
	curl, err := http2curl.GetCurlCommand(request)
	if err != nil {
		return "", fmt.Errorf("prepareCurl: failed preparing curl: %v", err)
	}

	return curl.String(), nil
}
