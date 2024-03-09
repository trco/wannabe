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
	}

	for testKey, tc := range testCases {
		processedPath, _ := ProcessPath(tc.Path, tc.Config)

		if !reflect.DeepEqual(tc.Expected, processedPath) {
			t.Errorf("Failed test case: %v, Expected: %v, Actual: %v", testKey, tc.Expected, processedPath)
		}
	}
}
