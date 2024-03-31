package actions

import "testing"

func TestProcessHttpMethod(t *testing.T) {
	methodString := "GET"
	httpMethod := ProcessHttpMethod(methodString)

	expected := "GET"

	if httpMethod != expected {
		t.Errorf("expected http method: %s, actual http method: %s", expected, httpMethod)
	}
}
