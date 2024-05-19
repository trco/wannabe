package types

import "time"

type RecordPayload struct {
	Hash            string
	Curl            string
	HttpMethod      string
	Host            string
	Path            string
	Query           map[string][]string
	RequestHeaders  map[string][]string
	ResponseHeaders map[string][]string
	RequestBody     []byte
	StatusCode      int
	ResponseBody    []byte
}

type Record struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
	Metadata Metadata `json:"metadata"`
}

type Request struct {
	Hash       string              `json:"hash"`
	Curl       string              `json:"curl"`
	HttpMethod string              `json:"httpMethod" validate:"required,oneof=GET POST PUT DELETE PATCH HEAD CONNECT OPTIONS TRACE"`
	Host       string              `json:"host" validate:"required"`
	Path       string              `json:"path"`
	Query      map[string][]string `json:"query"`
	Headers    map[string][]string `json:"headers"`
	Body       interface{}         `json:"body" validate:"required_if=HttpMethod POST,required_if=HttpMethod PUT,required_if=HttpMethod PATCH"`
}

type Response struct {
	StatusCode int                 `json:"statusCode" validate:"required"`
	Headers    map[string][]string `json:"headers"`
	Body       interface{}         `json:"body" validate:"required"`
}

type Metadata struct {
	GeneratedAt   Timestamp `json:"generatedAt"`
	RegeneratedAt Timestamp `json:"regeneratedAt"`
}

type Timestamp struct {
	Unix int64     `json:"unix"`
	UTC  time.Time `json:"utc"`
}
