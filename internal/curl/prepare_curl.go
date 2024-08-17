package curl

import (
	"fmt"
	"net/http"

	"moul.io/http2curl"
)

func PrepareCurl(request *http.Request) (string, error) {
	curl, err := http2curl.GetCurlCommand(request)
	if err != nil {
		return "", fmt.Errorf("PrepareCurl: failed preparing curl: %v", err)
	}

	return curl.String(), nil
}
