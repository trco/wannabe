package actions

import (
	"io"
	"testing"
)

func TestPrepareRequest(t *testing.T) {
	httpMethod := "POST"
	url := "http://test.com/test?test=test"
	body := "{\"test\":\"test\"}"

	request, _ := PrepareRequest(httpMethod, url, body)

	if request.Method != httpMethod {
		t.Errorf("expected httpMethod: %s, actual httpMethod: %s", httpMethod, request.Method)
	}

	if request.URL.String() != url {
		t.Errorf("expected url: %s, actual url: %s", url, request.URL)
	}

	requestBody, _ := io.ReadAll(request.Body)

	if string(requestBody) != body {
		t.Errorf("expected body: %s, actual body: %s", body, string(requestBody))
	}
}
