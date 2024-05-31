package actions

import (
	"reflect"
	"testing"
	"wannabe/types"
)

type TestCaseProcessPath struct {
	Path     string
	Config   types.Path
	Expected string
}

func TestProcessPath(t *testing.T) {
	zero := 0

	testCases := map[string]TestCaseProcessPath{
		"withPlaceholder": {
			Path: "/test1/test2/123456:test",
			Config: types.Path{
				Wildcards: []types.WildcardIndex{{Index: &zero, Placeholder: "{placeholder}"}},
			},
			Expected: "/{placeholder}/test2/123456:test",
		},
		"withoutPlaceholder": {
			Path: "/test1/test2/123456:test",
			Config: types.Path{
				Wildcards: []types.WildcardIndex{{Index: &zero}},
			},
			Expected: "/{wannabe}/test2/123456:test",
		},
		"withRegex": {
			Path: "/test1/test2/123456:test",
			Config: types.Path{
				Regexes: []types.Regex{{Pattern: "(\\d+):test", Placeholder: "{id}:test"}},
			},
			Expected: "/test1/test2/{id}:test",
		},
		"emptyString": {
			Path: "",
			Config: types.Path{
				Wildcards: []types.WildcardIndex{{Index: &zero}},
			},
			Expected: "",
		},
		"invalidRegex": {
			Path: "/test1/test2/123456:test",
			Config: types.Path{
				Regexes: []types.Regex{{Pattern: "(?P<foo"}},
			},
			Expected: "",
		},
	}

	for testKey, tc := range testCases {
		processedPath, err := ProcessPath(tc.Path, tc.Config)

		if testKey == "invalidRegex" && err != nil {
			expectedErr := "ProcessPath: failed compiling regex: error parsing regexp: invalid named capture: `(?P<foo`"

			if err.Error() != expectedErr {
				t.Errorf("expected error: %s, actual error: %s", expectedErr, err.Error())
			}

			continue
		}

		if !reflect.DeepEqual(tc.Expected, processedPath) {
			t.Errorf("Failed test case: %v, Expected: %v, Actual: %v", testKey, tc.Expected, processedPath)
		}
	}
}
