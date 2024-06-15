package actions

import (
	"testing"
)

func TestGenerateHash(t *testing.T) {
	curl := "curl -X 'POST' -d '{\"dateRanges\":[{\"endDate\":\"2023-12-31\",\"startDate\":\"{placeholder}\"],\"dimensions\":\"{placeholder}\",\"limit\":10000,\"metrics\":[{\"name\":\"sessions\"}],\"returnPropertyQuota\":true}' -H 'Accept: test1,test2,test3' -H 'Authorization: test access token' -H 'Content-Type: application/json' -H 'X-Test-Header: test value' 'analyticsdata.googleapis.com/v1beta/properties/%7BpropertyId%7D:runReport?app=1&status=new&user=%7Bplaceholder%7D'"

	want := "598532c061134aa3e816350fa27298931b72f80117541555edc1f0ad77d67566"

	t.Run("prepare curl", func(t *testing.T) {
		got, _ := GenerateHash(curl)

		if got != want {
			t.Errorf("GenerateHash() = %v, want %v", got, want)
		}
	})
}
