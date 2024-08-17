package hash

import (
	"testing"
)

func TestPrepareCurl(t *testing.T) {
	t.Run("prepare curl", func(t *testing.T) {
		request := generateTestRequest()

		want := "curl -X 'POST' -d '{\"test\":\"test\"}' -H 'Accept: test' -H 'Content-Type: application/json' 'https://test.com/test?test=test'"

		got, _ := prepareCurl(request)

		if got != want {
			t.Errorf("prepareCurl() = %v, want %v", got, want)
		}
	})
}
