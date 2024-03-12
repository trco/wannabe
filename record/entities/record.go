package entities

type Record struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

type Request struct {
	Hash       string              `json:"hash"`
	Curl       string              `json:"curl"`
	HttpMethod string              `json:"httpMethod"`
	Host       string              `json:"host"`
	Path       string              `json:"path"`
	Query      map[string]string   `json:"query"`
	Headers    map[string][]string `json:"headers"`
	Body       BodyMap             `json:"body"`
}

type Response struct {
	StatusCode int                 `json:"statusCode"`
	Headers    map[string][]string `json:"headers"`
	Body       BodyMap             `json:"body"`
}

type BodyMap map[string]interface{}
