package entities

type GenerateCurlPayload struct {
	HttpMethod     string
	Path           string
	Query          map[string]string
	RequestHeaders map[string][]string
	RequestBody    []byte
}
