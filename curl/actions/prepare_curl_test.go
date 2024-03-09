package actions

import (
	"bytes"
	"net/http"
	"testing"
)

func TestPrepareCurl(t *testing.T) {
	httpMethod := "POST"
	url := "analyticsdata.googleapis.com/v1beta/properties/%7BpropertyId%7D:runReport?app=1&status=new&user=%7Bplaceholder%7D"
	body := "{\"dateRanges\":[{\"endDate\":\"2023-12-31\",\"startDate\":\"{placeholder}\"],\"dimensions\":\"{placeholder}\",\"limit\":10000,\"metrics\":[{\"name\":\"sessions\"}],\"returnPropertyQuota\":true}"
	bodyBuffer := bytes.NewBufferString(body)
	request, _ := http.NewRequest(httpMethod, url, bodyBuffer)
	headers := []Header{
		{Key: "Accept", Value: "test1,test2,test3"},
		{Key: "Authorization", Value: "test access token"},
		{Key: "Content-Type", Value: "application/json"},
		{Key: "X-test-header", Value: "test value"},
	}
	SetHeaders(request, headers)

	curl, _ := PrepareCurl(request)

	expectedCurl := "curl -X 'POST' -d '{\"dateRanges\":[{\"endDate\":\"2023-12-31\",\"startDate\":\"{placeholder}\"],\"dimensions\":\"{placeholder}\",\"limit\":10000,\"metrics\":[{\"name\":\"sessions\"}],\"returnPropertyQuota\":true}' -H 'Accept: test1,test2,test3' -H 'Authorization: test access token' -H 'Content-Type: application/json' -H 'X-Test-Header: test value' 'analyticsdata.googleapis.com/v1beta/properties/%7BpropertyId%7D:runReport?app=1&status=new&user=%7Bplaceholder%7D'"

	if curl != expectedCurl {
		t.Errorf("Expected curl: %s, Actual curl: %s", expectedCurl, curl)
	}
}
