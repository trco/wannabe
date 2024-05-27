package actions

import (
	"reflect"
	"testing"
	"wannabe/types"
)

type TestCaseProcessHeaders struct {
	Map      map[string][]string
	Config   types.Headers
	Expected []types.Header
}

func TestProcessHeaders(t *testing.T) {
	testCases := map[string]TestCaseProcessHeaders{
		"includeAllHeaders": {
			Map: testInitHeadersMap(),
			Config: types.Headers{
				Include:   []string{"Content-Type", "Authorization", "Accept", "X-test-header"},
				Wildcards: []types.WildcardKey{},
			},
			Expected: []types.Header{
				{Key: "Accept", Value: "test1,test2,test3"},
				{Key: "Authorization", Value: "test access token"},
				{Key: "Content-Type", Value: "application/json"},
				{Key: "X-test-header", Value: "test value"},
			},
		},
		"withPlaceholderTwoHeaders": {
			Map: testInitHeadersMap(),
			Config: types.Headers{
				Include:   []string{"Content-Type", "Authorization"},
				Wildcards: []types.WildcardKey{{Key: "Authorization", Placeholder: "{placeholder}"}},
			},
			Expected: []types.Header{
				{Key: "Authorization", Value: "{placeholder}"},
				{Key: "Content-Type", Value: "application/json"},
			},
		},
		"withoutPlaceholderTwoHeaders": {
			Map: testInitHeadersMap(),
			Config: types.Headers{
				Include:   []string{"Content-Type", "Authorization"},
				Wildcards: []types.WildcardKey{{Key: "Authorization"}},
			},
			Expected: []types.Header{
				{Key: "Authorization", Value: "{wannabe}"},
				{Key: "Content-Type", Value: "application/json"},
			},
		},
		"emptyHeadersMap": {
			Map: make(map[string][]string),
			Config: types.Headers{
				Include:   []string{"Content-Type", "Authorization"},
				Wildcards: []types.WildcardKey{{Key: "Authorization"}},
			},
			Expected: []types.Header{},
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
