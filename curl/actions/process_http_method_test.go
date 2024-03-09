package actions

import "testing"

func TestProcessHttpMethod(t *testing.T) {
	methodString := "GET"
	httpMethod := ProcessHttpMethod(methodString)

	expected := "GET"

	if httpMethod != expected {
		t.Errorf("Expected: %s, Actual: %s", expected, httpMethod)
	}
}
