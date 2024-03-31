package actions

import (
	"reflect"
	"testing"
	"wannabe/config"
)

type TestCaseProcessPath struct {
	Path     string
	Config   config.Path
	Expected string
}

func TestProcessPath(t *testing.T) {
	zero := 0

	testCases := map[string]TestCaseProcessPath{
		"withPlaceholder": {
			Path: "/v1beta/properties/375748157:runReport",
			Config: config.Path{
				Wildcards: []config.WildcardIndex{{Index: &zero, Placeholder: "{placeholder}"}},
			},
			Expected: "/{placeholder}/properties/375748157:runReport",
		},
		"withoutPlaceholder": {
			Path: "/v1beta/properties/375748157:runReport",
			Config: config.Path{
				Wildcards: []config.WildcardIndex{{Index: &zero}},
			},
			Expected: "/{wannabe}/properties/375748157:runReport",
		},
		"withRegex": {
			Path: "/v1beta/properties/375748157:runReport",
			Config: config.Path{
				Regexes: []config.Regex{{Pattern: "(\\d+):runReport", Placeholder: "{propertyId}:runReport"}},
			},
			Expected: "/v1beta/properties/{propertyId}:runReport",
		},
		"emptyString": {
			Path: "",
			Config: config.Path{
				Wildcards: []config.WildcardIndex{{Index: &zero}},
			},
			Expected: "",
		},
		"invalidRegex": {
			Path: "/v1beta/properties/375748157:runReport",
			Config: config.Path{
				Regexes: []config.Regex{{Pattern: "(?P<foo"}},
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
