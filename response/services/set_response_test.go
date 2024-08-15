package services

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/trco/wannabe/types"
)

func TestSetResponse(t *testing.T) {
	tests := []struct {
		name                string
		record              []byte
		request             *http.Request
		wantStatusCode      int
		wantResponseHeaders map[string][]string
		wantResponseBody    map[string]interface{}
		wantErr             string
	}{
		{
			name:                "invalid record",
			record:              []byte{53},
			request:             request,
			wantStatusCode:      0,
			wantResponseHeaders: nil,
			wantResponseBody:    nil,
			wantErr:             "SetResponse: failed unmarshaling record: json: cannot unmarshal number into Go value of type types.Record",
		},
		{
			name:           "valid record",
			record:         record,
			request:        request,
			wantStatusCode: 200,
			wantResponseHeaders: map[string][]string{
				"Content-Type": {"application/json"},
				"Accept":       {"test"},
			},
			wantResponseBody: map[string]interface{}{
				"test": "test",
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetResponse(tt.record, tt.request)

			if err != nil {
				if err.Error() != tt.wantErr {
					t.Errorf("SetResponse() error = %v, wantErr %v", err.Error(), tt.wantErr)
				}
				return
			}

			if !reflect.DeepEqual(got.StatusCode, tt.wantStatusCode) {
				t.Errorf("got status code %v, want status code %v", got.StatusCode, tt.wantStatusCode)
			}

			if !(len(got.Header) == len(tt.wantResponseHeaders)) {
				t.Errorf("got %v response headers , want %v response headers", len(got.Header), len(tt.wantResponseHeaders))
			}

			responseBodyEncoded, _ := io.ReadAll(got.Body)
			defer got.Body.Close()

			var gotResponseBody map[string]interface{}
			json.Unmarshal(responseBodyEncoded, &gotResponseBody)

			if !reflect.DeepEqual(gotResponseBody, tt.wantResponseBody) {
				t.Errorf("got response body: %v, want response body: %v", gotResponseBody, tt.wantResponseBody)
			}
		})
	}
}

var request = &http.Request{
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

var record, _ = json.Marshal(types.Record{
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
