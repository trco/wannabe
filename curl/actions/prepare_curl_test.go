package actions

import (
	"testing"
)

func TestPrepareCurl(t *testing.T) {
	t.Run("prepare curl", func(t *testing.T) {
		request := generateTestRequest()

		want := "curl -X 'POST' -d '{\"test\":\"test\"}' -H 'Accept: test' -H 'Content-Type: application/json' 'http://test.com/test?test=test'"

		got, _ := PrepareCurl(request)

		if got != want {
			t.Errorf("PrepareCurl() = %v, want %v", got, want)
		}
	})
}
