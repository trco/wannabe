package hash

import "testing"

func TestProcessHttpMethod(t *testing.T) {
	t.Run("process http method", func(t *testing.T) {
		httpMethod := "GET"

		want := "GET"

		got := processHttpMethod(httpMethod)

		if got != want {
			t.Errorf("processHttpMethod() = %v, want %v", got, want)
		}
	})
}
