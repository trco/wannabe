package curl

import (
	"reflect"
	"testing"

	"github.com/trco/wannabe/internal/config"
)

func TestGenerateCurl(t *testing.T) {
	t.Run("generate curl", func(t *testing.T) {
		request := generateTestRequest()

		wannabe := config.Wannabe{
			RequestMatching: config.RequestMatching{
				Headers: config.Headers{
					Include: []string{"Content-Type", "Accept"},
				},
			},
		}

		want := "curl -X 'POST' -d '{\"test\":\"test\"}' -H 'Accept: test' -H 'Content-Type: application/json' 'https://test.com/test?test=test'"

		got, _ := generateCurl(request, wannabe)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("GenerateCurl() = %v, want %v", got, want)
		}
	})

}
