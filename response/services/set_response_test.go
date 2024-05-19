package services

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"wannabe/types"
)

func TestSetResponse(t *testing.T) {
	request := &http.Request{
		Method: "POST",
		Host:   "test.com",
		URL: &url.URL{
			Path:     "/test",
			RawQuery: "test=test",
		},
		Header: map[string][]string{
			"Content-Type": {"application/json"},
			"Accept":       {"test"},
		},
		Body: io.NopCloser(bytes.NewBuffer([]byte(`{"test": "test"}`))),
	}

	record, _ := json.Marshal(types.Record{
		Request: types.Request{
			HttpMethod: "POST",
			Host:       "test.com",
			Path:       "/test",
			Query: map[string][]string{
				"test": {"test"},
			},
			Headers: map[string][]string{
				"Content-Type": {"application/json"},
				"Accept":       {"test"},
			},
			Body: map[string]interface{}{
				"test": "test",
			},
		},
		Response: types.Response{
			StatusCode: 200,
			Headers: map[string][]string{
				"Content-Type": {"application/json"},
				"Accept":       {"test"},
			},
			Body: map[string]interface{}{
				"test": "test",
			},
		},
	})

	// invalid record
	_, err := SetResponse([]byte{53}, request)
	if err == nil {
		t.Errorf("invalid record didn't return error.")
	}

	// valid record
	response, _ := SetResponse(record, request)

	expectedStatusCode := 200
	expectedResponseHeaders := map[string][]string{
		"Content-Type": {"application/json"},
		"Accept":       {"test"},
	}
	expectedResponseBody := map[string]interface{}{
		"test": "test",
	}

	if !reflect.DeepEqual(expectedStatusCode, response.StatusCode) {
		t.Errorf("expected status code: %v, actual status code: %v", expectedStatusCode, response.StatusCode)
	}

	if !reflect.DeepEqual(expectedResponseHeaders["Accept"], response.Header["Accept"]) {
		t.Errorf("expected response headers: %v, actual response headers: %v", expectedResponseHeaders, response.Header)
	}

	if !reflect.DeepEqual(expectedResponseHeaders["Content-Type"], response.Header["Content-Type"]) {
		t.Errorf("expected response headers: %v, actual response headers: %v", expectedResponseHeaders, response.Header)
	}

	responseBodyRaw, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	var responseBody map[string]interface{}
	json.Unmarshal(responseBodyRaw, &responseBody)

	if !reflect.DeepEqual(expectedResponseBody, responseBody) {
		t.Errorf("expected response body: %v, actual response body: %v", expectedResponseBody, responseBody)
	}
}
