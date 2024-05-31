package actions

import (
	"testing"
)

func TestPrepareCurl(t *testing.T) {
	request := generateTestRequest()

	curl, _ := PrepareCurl(request)

	expectedCurl := "curl -X 'POST' -d '{\"test\":\"test\"}' -H 'Accept: test' -H 'Content-Type: application/json' 'http://test.com/test?test=test'"

	if curl != expectedCurl {
		t.Errorf("expected curl: %s, actual curl: %s", expectedCurl, curl)
	}
}
