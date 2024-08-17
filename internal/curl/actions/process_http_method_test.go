package actions

import "testing"

func TestProcessHttpMethod(t *testing.T) {
	t.Run("process http method", func(t *testing.T) {
		httpMethod := "GET"

		want := "GET"

		got := ProcessHttpMethod(httpMethod)

		if got != want {
			t.Errorf("ProcessHttpMethod() = %v, want %v", got, want)
		}
	})
}
