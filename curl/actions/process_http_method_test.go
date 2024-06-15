package actions

import "testing"

func TestProcessHttpMethod(t *testing.T) {
	httpMethod := "GET"

	want := "GET"

	t.Run("process http method", func(t *testing.T) {
		got := ProcessHttpMethod(httpMethod)

		if got != want {
			t.Errorf("ProcessHttpMethod() = %v, want %v", got, want)
		}
	})
}
