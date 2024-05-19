package services

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"wannabe/types"
)

func TestGenerateRecordPayload(t *testing.T) {
	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Host: "test.com",
			Path: "/test",
		},
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewBufferString("request body")),
	}

	res := &http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewBufferString("response body")),
	}

	session := types.WannabeSession{
		Req: req,
		Res: res,
	}

	hash := "test hash"
	curl := "test curl"

	expectedRecordPayload := types.RecordPayload{
		Hash:            hash,
		Curl:            curl,
		HttpMethod:      req.Method,
		Host:            req.URL.Host,
		Path:            req.URL.Path,
		Query:           req.URL.Query(),
		RequestHeaders:  req.Header,
		RequestBody:     []byte{114, 101, 113, 117, 101, 115, 116, 32, 98, 111, 100, 121},
		StatusCode:      res.StatusCode,
		ResponseHeaders: res.Header,
		ResponseBody:    []byte{114, 101, 115, 112, 111, 110, 115, 101, 32, 98, 111, 100, 121},
	}

	payload, _ := GenerateRecordPayload(session, hash, curl)

	if !reflect.DeepEqual(expectedRecordPayload, payload) {
		t.Errorf("expected record payload: %v, actual record payload: %v", expectedRecordPayload, payload)
	}
}
