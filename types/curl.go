package types

type CurlPayload struct {
	Scheme         string
	HttpMethod     string
	Host           string
	Path           string
	Query          map[string][]string
	RequestHeaders map[string][]string
	RequestBody    []byte
}

type Header struct {
	Key   string
	Value string
}
