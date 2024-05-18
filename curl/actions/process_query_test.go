package actions

import (
	"reflect"
	"testing"
	"wannabe/config"
)

type TestCaseProcessQuery struct {
	QueryMap map[string][]string
	Config   config.Query
	Expected string
}

func TestProcessQuery(t *testing.T) {
	testCases := map[string]TestCaseProcessQuery{
		"withPlaceholder": {
			QueryMap: testMapQuery,
			Config: config.Query{
				Wildcards: []config.WildcardKey{{Key: "user", Placeholder: "{placeholder}"}},
			},
			Expected: "?app=1&status=new&user=%7Bplaceholder%7D",
		},
		"withoutPlaceholder": {
			QueryMap: testMapQuery,
			Config: config.Query{
				Wildcards: []config.WildcardKey{{Key: "user"}},
			},
			Expected: "?app=1&status=new&user=%7Bwannabe%7D",
		},
		"withRegexWithPlaceholder": {
			QueryMap: testMapQuery,
			Config: config.Query{
				Regexes: []config.Regex{{Pattern: "paid", Placeholder: "{placeholder}"}},
			},
			Expected: "?app=1&status=new&user=%7Bplaceholder%7D",
		},
		"withRegexWithoutPlaceholder": {
			QueryMap: testMapQuery,
			Config: config.Query{
				Regexes: []config.Regex{{Pattern: "paid"}},
			},
			Expected: "?app=1&status=new&user=%7Bwannabe%7D",
		},
		"emptyString": {
			QueryMap: make(map[string][]string),
			Config: config.Query{
				Wildcards: []config.WildcardKey{{Key: "user"}},
			},
			Expected: "",
		},
		"invalidRegex": {
			QueryMap: testMapQuery,
			Config: config.Query{
				Regexes: []config.Regex{{Pattern: "(?P<foo"}},
			},
			Expected: "",
		},
	}

	for testKey, tc := range testCases {
		processedQuery, err := ProcessQuery(tc.QueryMap, tc.Config)

		if testKey == "invalidRegex" && err != nil {
			expectedErr := "ProcessQuery: failed compiling regex: error parsing regexp: invalid named capture: `(?P<foo`"

			if err.Error() != expectedErr {
				t.Errorf("expected error: %s, actual error: %s", expectedErr, err.Error())
			}

			continue
		}

		if !reflect.DeepEqual(tc.Expected, processedQuery) {
			t.Errorf("failed test case: %v, expected query: %v, actual query: %v", testKey, tc.Expected, processedQuery)
		}
	}
}

// reusable variables

var testMapQuery = map[string][]string{
	"status": {"new"},
	"user":   {"paid"},
	"app":    {"1"},
}
