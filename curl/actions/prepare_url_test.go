package actions

import "testing"

func TestPrepareUrl(t *testing.T) {
	host := "test.com"
	path := "/test"
	query := "?test=test"

	url := PrepareUrl(host, path, query)

	expectedUrl := "test.com/test?test=test"

	if expectedUrl != url {
		t.Errorf("expected url: %s, actual url: %s", expectedUrl, url)
	}
}
