package hash

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	t.Run("prepare curl", func(t *testing.T) {
		curl := "curl -X 'POST' -d '{\"dateRanges\":[{\"endDate\":\"2023-12-31\",\"startDate\":\"placeholder\"],\"dimensions\":\"placeholder\",\"limit\":10000,\"metrics\":[{\"name\":\"sessions\"}],\"returnPropertyQuota\":true}' -H 'Accept: test1,test2,test3' -H 'Authorization: test access token' -H 'Content-Type: application/json' -H 'X-Test-Header: test value' 'analyticsdata.googleapis.com/v1beta/properties/%7BpropertyId%7D:runReport?app=1&status=new&user=%7Bplaceholder%7D'"

		want := "8daeab258a6c5b5999800355fc64b51acefd7fe673aaafdadd0591ea7063c603"

		got, _ := Generate(curl)

		if got != want {
			t.Errorf("GenerateHash() = %v, want %v", got, want)
		}
	})
}
