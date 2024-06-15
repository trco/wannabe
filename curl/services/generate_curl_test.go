package services

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"
	"wannabe/types"
)

func TestGenerateCurl(t *testing.T) {
	t.Run("generate curl", func(t *testing.T) {
		request := generateTestRequest()

		wannabe := types.Wannabe{
			RequestMatching: types.RequestMatching{
				Headers: types.Headers{
					Include: []string{"Content-Type", "Accept"},
				},
			},
		}

		want := "curl -X 'POST' -d '{\"test\":\"test\"}' -H 'Accept: test' -H 'Content-Type: application/json' 'test.com/test?test=test'"

		got, _ := GenerateCurl(request, wannabe)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("GenerateCurl() = %v, want %v", got, want)
		}
	})

}

func generateTestRequest() *http.Request {
	httpMethod := "POST"
	url := "http://test.com/test?test=test"
	body := "{\"test\":\"test\"}"
	bodyBuffer := bytes.NewBufferString(body)

	request, _ := http.NewRequest(httpMethod, url, bodyBuffer)

	request.Header.Set("Accept", "test")
	request.Header.Set("Content-Type", "application/json")

	return request
}
