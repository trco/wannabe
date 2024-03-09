package actions

import (
	"io"
	"testing"
)

func TestPrepareRequest(t *testing.T) {
	httpMethod := "POST"
	url := "analyticsdata.googleapis.com/v1beta/properties/%7BpropertyId%7D:runReport?app=1&status=new&user=%7Bplaceholder%7D"
	body := "{\"dateRanges\":[{\"endDate\":\"2023-12-31\",\"startDate\":\"{placeholder}\"],\"dimensions\":\"{placeholder}\",\"limit\":10000,\"metrics\":[{\"name\":\"sessions\"}],\"returnPropertyQuota\":true}"

	request, _ := PrepareRequest(httpMethod, url, body)

	if request.Method != httpMethod {
		t.Errorf("Expected httpMethod: %s, Actual httpMethod: %s", httpMethod, request.Method)
	}

	if request.URL.String() != url {
		t.Errorf("Expected url: %s, Actual url: %s", url, request.URL)
	}

	requestBody, _ := io.ReadAll(request.Body)

	if string(requestBody) != body {
		t.Errorf("Expected body: %s, Actual body: %s", body, string(requestBody))
	}
}
