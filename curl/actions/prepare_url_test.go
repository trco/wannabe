package actions

import "testing"

func TestPrepareUrl(t *testing.T) {
	host := "analyticsdata.googleapis.com"
	path := "/v1beta/properties/{propertyId}:runReport"
	query := "?app=1&status=new&user=%7Bplaceholder%7D"

	url := PrepareUrl(host, path, query)

	expectedUrl := "analyticsdata.googleapis.com/v1beta/properties/{propertyId}:runReport?app=1&status=new&user=%7Bplaceholder%7D"

	if expectedUrl != url {
		t.Errorf("expected url: %s, actual url: %s", expectedUrl, url)
	}
}
