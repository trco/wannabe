package actions

import (
	"reflect"
	"testing"
	"wannabe/config"
)

type TestCaseProcessHeaders struct {
	Map      map[string][]string
	Config   config.Headers
	Expected []Header
}

func TestProcessHeaders(t *testing.T) {
	testCases := map[string]TestCaseProcessHeaders{
		"includeAllHeaders": {
			Map: testInitHeadersMap(),
			Config: config.Headers{
				Include:   []string{"Content-Type", "Authorization", "Accept", "X-test-header"},
				Wildcards: []config.WildcardKey{},
			},
			Expected: []Header{
				{Key: "Accept", Value: "test1,test2,test3"},
				{Key: "Authorization", Value: "test access token"},
				{Key: "Content-Type", Value: "application/json"},
				{Key: "X-test-header", Value: "test value"},
			},
		},
		"withPlaceholderTwoHeaders": {
			Map: testInitHeadersMap(),
			Config: config.Headers{
				Include:   []string{"Content-Type", "Authorization"},
				Wildcards: []config.WildcardKey{{Key: "Authorization", Placeholder: "{placeholder}"}},
			},
			Expected: []Header{
				{Key: "Authorization", Value: "{placeholder}"},
				{Key: "Content-Type", Value: "application/json"},
			},
		},
		"withoutPlaceholderTwoHeaders": {
			Map: testInitHeadersMap(),
			Config: config.Headers{
				Include:   []string{"Content-Type", "Authorization"},
				Wildcards: []config.WildcardKey{{Key: "Authorization"}},
			},
			Expected: []Header{
				{Key: "Authorization", Value: "{wannabe}"},
				{Key: "Content-Type", Value: "application/json"},
			},
		},
		"emptyHeadersMap": {
			Map: make(map[string][]string),
			Config: config.Headers{
				Include:   []string{"Content-Type", "Authorization"},
				Wildcards: []config.WildcardKey{{Key: "Authorization"}},
			},
			Expected: []Header{},
		},
	}

	for testKey, tc := range testCases {
		sortedHeaders := ProcessHeaders(tc.Map, tc.Config)

		if testKey == "emptyHeadersMap" && !(len(tc.Map) == 0 && len(sortedHeaders) == 0) {
			t.Errorf("failed test case: %v, expected headers: %v, actual headers: %v", testKey, tc.Expected, sortedHeaders)
		}

		if testKey != "emptyHeadersMap" && !reflect.DeepEqual(tc.Expected, sortedHeaders) {
			t.Errorf("failed test case: %v, expected headers: %v, actual headers: %v", testKey, tc.Expected, sortedHeaders)
		}
	}
}
