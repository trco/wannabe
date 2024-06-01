package actions

import "testing"

func TestPrepareUrl(t *testing.T) {
	host := "test.com"
	query := "?test=test"

	path := "/test"

	url := PrepareUrl(host, path, query)

	expectedUrl := "test.com/test?test=test"

	if expectedUrl != url {
		t.Errorf("expected url: %s, actual url: %s, when path is '/test'", expectedUrl, url)
	}

	path = "/"

	url = PrepareUrl(host, path, query)

	expectedUrl = "test.com?test=test"

	if expectedUrl != url {
		t.Errorf("expected url: %s, actual url: %s, when path is '/'", expectedUrl, url)
	}
}
