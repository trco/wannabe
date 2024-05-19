package types

type CurlPayload struct {
	HttpMethod     string
	Host           string
	Path           string
	Query          map[string][]string
	RequestHeaders map[string][]string
	RequestBody    []byte
}
