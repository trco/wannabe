package actions

import (
	"reflect"
	"testing"
	"wannabe/config"
)

type TestCaseProcessQuery struct {
	Query    map[string]string
	Config   config.Query
	Expected string
}

func testMapQuery() map[string]string {
	return map[string]string{
		"status": "new",
		"user":   "paid",
		"app":    "1",
	}
}

func TestProcessQuery(t *testing.T) {
	testCases := map[string]TestCaseProcessQuery{
		"withPlaceholder": {
			Query: testMapQuery(),
			Config: config.Query{
				Wildcards: []config.WildcardKey{{Key: "user", Placeholder: "{placeholder}"}},
			},
			Expected: "?app=1&status=new&user=%7Bplaceholder%7D",
		},
		"withoutPlaceholder": {
			Query: testMapQuery(),
			Config: config.Query{
				Wildcards: []config.WildcardKey{{Key: "user"}},
			},
			Expected: "?app=1&status=new&user=%7Bwannabe%7D",
		},
		"withRegex": {
			Query: testMapQuery(),
			Config: config.Query{
				Regexes: []config.Regex{{Pattern: "app=1", Placeholder: "app=123"}},
			},
			Expected: "?app=123&status=new&user=paid",
		},
		"emptyString": {
			Query: make(map[string]string),
			Config: config.Query{
				Wildcards: []config.WildcardKey{{Key: "user"}},
			},
			Expected: "",
		},
	}

	for testKey, tc := range testCases {
		processedQuery, _ := ProcessQuery(tc.Query, tc.Config)

		if !reflect.DeepEqual(tc.Expected, processedQuery) {
			t.Errorf("Failed test case: %v, Expected: %v, Actual: %v", testKey, tc.Expected, processedQuery)
		}
	}
}
