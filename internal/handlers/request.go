package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/clbanning/mxj"
	"github.com/trco/wannabe/internal/config"
	"github.com/trco/wannabe/internal/hash"
	"github.com/trco/wannabe/internal/record"
	"github.com/trco/wannabe/internal/storage"

	"github.com/AdguardTeam/gomitmproxy"
	"github.com/AdguardTeam/gomitmproxy/proxyutil"
)

type RequestHandler func(*gomitmproxy.Session) (*http.Request, *http.Response)

func Request(cfg config.Config, storageProvider storage.Provider) RequestHandler {
	return func(session *gomitmproxy.Session) (*http.Request, *http.Response) {
		return processRequest(cfg, storageProvider, session)
	}
}

func processRequest(cfg config.Config, storageProvider storage.Provider, session *gomitmproxy.Session) (*http.Request, *http.Response) {
	request := session.Request()

	isConnect := request.Method == "CONNECT"
	if isConnect {
		return nil, nil
	}

	removeBody(request)
	setScheme(request)

	host := request.URL.Host
	wannabe := cfg.Wannabes[host]

	curl, err := hash.GenerateCurl(request, wannabe)
	if err != nil {
		return internalErrorOnRequest(session, request, err)
	}
	session.SetProp("curl", curl)

	hash, err := hash.Generate(curl)
	if err != nil {
		return internalErrorOnRequest(session, request, err)
	}
	session.SetProp("hash", hash)

	isNotProxyMode := cfg.Mode != config.ProxyMode
	if isNotProxyMode {
		records, err := storageProvider.ReadRecords(host, []string{hash})
		if err != nil {
			return internalErrorOnRequest(session, request, err)
		}

		isSingleRecord := len(records) == 1
		if isSingleRecord {
			return setResponseFromRecord(session, request, records[0])
		}

		isServerMode := cfg.Mode == config.ServerMode
		if isServerMode {
			return internalErrorOnRequest(session, request, fmt.Errorf("no record found for the request"))
		}
	}

	hasBody := request.Body != nil
	if hasBody {
		requestBody, err := copyBody(request)
		if err != nil {
			return internalErrorOnRequest(session, request, err)
		}
		session.SetProp("requestBody", requestBody)
	}

	return nil, nil
}

// prevents sending body and related headers in GET requests
func removeBody(request *http.Request) {
	isGet := request.Method == "GET"
	if isGet {
		request.Body = http.NoBody
		request.ContentLength = 0
		request.Header.Del("Content-Length")
		request.Header.Del("Content-Type")
	}
}

// sets scheme to "https" if request uri contains it but the scheme set on url is "http"
func setScheme(request *http.Request) {
	if strings.Contains(request.RequestURI, "https://") {
		request.URL.Scheme = "https"
	}
}

func setResponseFromRecord(session *gomitmproxy.Session, request *http.Request, record []byte) (*http.Request, *http.Response) {
	response, err := setResponse(record, request)
	if err != nil {
		return internalErrorOnRequest(session, request, err)
	}

	session.SetProp("responseSet", true)

	fmt.Println("Response successfully read from configured StorageProvider.")

	return nil, response
}

func setResponse(encodedRecord []byte, request *http.Request) (*http.Response, error) {
	var record record.Record

	err := json.Unmarshal(encodedRecord, &record)
	if err != nil {
		return nil, fmt.Errorf("SetResponse: failed unmarshaling record: %v", err)
	}

	decodedBody := record.Response.Body
	contentTypeHeader := record.Response.Headers["Content-Type"]
	contentEncodingHeader := record.Response.Headers["Content-Encoding"]
	body, err := encodeBody(decodedBody, contentTypeHeader, contentEncodingHeader)
	if err != nil {
		return nil, fmt.Errorf("SetResponse: failed encoding response body: %v", err)
	}

	statusCode := record.Response.StatusCode

	response := proxyutil.NewResponse(statusCode, bytes.NewReader(body), request)

	headers := record.Response.Headers
	setHeaders(response, headers)

	return response, nil
}

func encodeBody(decodedBody interface{}, contentTypeHeader []string, contentEncodingHeader []string) ([]byte, error) {
	var body []byte

	contentType := record.GetContentType(contentTypeHeader)
	switch {
	case contentType == "application/json":
		var err error
		body, err = json.Marshal(decodedBody)
		if err != nil {
			return nil, fmt.Errorf("SetResponse: failed marshaling body: %v", err)
		}
	case contentType == "application/xml", contentType == "text/xml":
		var err error
		mapValue := mxj.Map(decodedBody.(map[string]interface{}))
		body, err = mapValue.Xml()
		if err != nil {
			return body, fmt.Errorf("Generate: failed unmarshaling XML body: %v", err)
		}
	case contentType == "text/plain", contentType == "text/html":
		body = []byte(decodedBody.(string))
	default:
		return nil, fmt.Errorf("SetResponse: unsupported content type: %s", contentTypeHeader)
	}

	// gzip the body that was unzipped before storing it to the record
	contentEncoding := record.GetContentEncoding(contentEncodingHeader)
	if contentEncoding == "gzip" {
		compressedBody, err := record.Gzip(body)
		if err != nil {
			return nil, fmt.Errorf("SetResponse: failed compressing response body: %s", err)
		}

		return compressedBody, nil
	}

	return body, nil
}

func setHeaders(response *http.Response, headers map[string][]string) {
	for key, value := range headers {
		for _, v := range value {
			response.Header.Set(key, v)
		}
	}
}

func copyBody(request *http.Request) (io.ReadCloser, error) {
	var requestBody io.ReadCloser

	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, fmt.Errorf("CopyBody: failed reading request body: %v", err)
	}
	defer request.Body.Close()

	requestBody = io.NopCloser(bytes.NewReader(body))
	// set body back to the request
	request.Body = io.NopCloser(bytes.NewBuffer(body))

	return requestBody, nil
}

type Session interface {
	SetProp(key string, value interface{})
	GetProp(key string) (interface{}, bool)
	Request() *http.Request
	Response() *http.Response
}

func internalErrorOnRequest(session Session, request *http.Request, err error) (*http.Request, *http.Response) {
	session.SetProp("blocked", true)

	body := PrepareResponseBody(err)
	response := proxyutil.NewResponse(http.StatusInternalServerError, body, request)
	response.Header.Set("Content-Type", "application/json")

	return nil, response
}

func PrepareResponseBody(err error) *bytes.Reader {
	body, err := json.Marshal(InternalError{
		Error: err.Error(),
	})
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	bodyReader := bytes.NewReader(body)

	return bodyReader
}

type InternalError struct {
	Error string `json:"error"`
}
