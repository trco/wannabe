package utils

import (
	"reflect"
	"testing"
	"wannabe/types"
)

type TestCaseByIndex struct {
	Slice     []string
	Wildcards []types.WildcardIndex
	Expected  []string
}

func testSlice() []string {
	return []string{"analyticsdata", "googleapis", "com"}
}

func TestSetWildcardsByIndex(t *testing.T) {
	zero := 0
	one := 1
	five := 5

	testCases := map[string]TestCaseByIndex{
		"withPlaceholder": {
			Slice:     testSlice(),
			Wildcards: []types.WildcardIndex{{Index: &zero, Placeholder: "{placeholder}"}},
			Expected:  []string{"{placeholder}", "googleapis", "com"},
		},
		"withoutPlaceholder": {
			Slice:     testSlice(),
			Wildcards: []types.WildcardIndex{{Index: &zero}},
			Expected:  []string{"{wannabe}", "googleapis", "com"},
		},
		"withAndWithoutPlaceholder": {
			Slice:     testSlice(),
			Wildcards: []types.WildcardIndex{{Index: &zero, Placeholder: "{placeholder}"}, {Index: &one}},
			Expected:  []string{"{placeholder}", "{wannabe}", "com"},
		},
		"indexOutOfBounds": {
			Slice:     testSlice(),
			Wildcards: []types.WildcardIndex{{Index: &five}},
			Expected:  testSlice(),
		},
	}

	for testKey, tc := range testCases {
		SetWildcardsByIndex(tc.Slice, tc.Wildcards)

		if !reflect.DeepEqual(tc.Expected, tc.Slice) {
			t.Errorf("failed test case: %v, expected slice: %v, actual slice: %v", testKey, tc.Expected, tc.Slice)
		}
	}
}

type TestCaseByKey struct {
	Map       map[string]string
	Wildcards []types.WildcardKey
	Expected  map[string]string
}

func testMap() map[string]string {
	return map[string]string{
		"status": "test",
		"appId":  "test",
	}
}

func TestSetWildcardsByKey(t *testing.T) {
	testCases := map[string]TestCaseByKey{
		"withPlaceholder": {
			Map:       testMap(),
			Wildcards: []types.WildcardKey{{Key: "status", Placeholder: "{placeholder}"}},
			Expected: map[string]string{
				"status": "{placeholder}",
				"appId":  "test",
			},
		},
		"withoutPlaceholder": {
			Map:       testMap(),
			Wildcards: []types.WildcardKey{{Key: "status"}},
			Expected: map[string]string{
				"status": "{wannabe}",
				"appId":  "test",
			},
		},
		"withAndWithoutPlaceholder": {
			Map:       testMap(),
			Wildcards: []types.WildcardKey{{Key: "status", Placeholder: "{placeholder}"}, {Key: "appId"}},
			Expected: map[string]string{
				"status": "{placeholder}",
				"appId":  "{wannabe}",
			},
		},
		"nonExistingKey": {
			Map:       testMap(),
			Wildcards: []types.WildcardKey{{Key: "nonExistingKey"}},
			Expected:  testMap(),
		},
	}

	for testKey, tc := range testCases {
		SetWildcardsByKey(tc.Map, tc.Wildcards)

		if !reflect.DeepEqual(tc.Expected, tc.Map) {
			t.Errorf("failed test case: %v, expected map: %v, actual map: %v", testKey, tc.Expected, tc.Map)
		}
	}
}

type TestCasePlaceholderByIndex struct {
	Slice     []string
	Wildcards types.WildcardIndex
	Expected  []string
}

func TestSetPlaceholderByIndex(t *testing.T) {
	zero := 0

	testCases := map[string]TestCasePlaceholderByIndex{
		"withPlaceholder": {
			Slice:     testSlice(),
			Wildcards: types.WildcardIndex{Index: &zero, Placeholder: "{placeholder}"},
			Expected:  []string{"{placeholder}", "googleapis", "com"},
		},
		"withoutPlaceholder": {
			Slice:     testSlice(),
			Wildcards: types.WildcardIndex{Index: &zero},
			Expected:  []string{"{wannabe}", "googleapis", "com"},
		},
	}

	for testKey, tc := range testCases {
		SetPlaceholderByIndex(tc.Slice, tc.Wildcards)

		if !reflect.DeepEqual(tc.Expected, tc.Slice) {
			t.Errorf("failed test case: %v, expected slice: %v, actual slice: %v", testKey, tc.Expected, tc.Slice)
		}
	}
}

type TestCasePlaceholderByKey struct {
	Map      map[string]string
	Wildcard types.WildcardKey
	Expected map[string]string
}

func TestSetPlaceholderByKey(t *testing.T) {
	testCases := map[string]TestCasePlaceholderByKey{
		"withPlaceholder": {
			Map:      testMap(),
			Wildcard: types.WildcardKey{Key: "status", Placeholder: "{placeholder}"},
			Expected: map[string]string{
				"status": "{placeholder}",
				"appId":  "test",
			},
		},
		"withoutPlaceholder": {
			Map:      testMap(),
			Wildcard: types.WildcardKey{Key: "status"},
			Expected: map[string]string{
				"status": "{wannabe}",
				"appId":  "test",
			},
		},
	}

	for testKey, tc := range testCases {
		SetPlaceholderByKey(tc.Map, tc.Wildcard)

		if !reflect.DeepEqual(tc.Expected, tc.Map) {
			t.Errorf("failed test case: %v, expected map: %v, actual map: %v", testKey, tc.Expected, tc.Map)
		}
	}
}

type TestCaseReplaceRegexPatterns struct {
	String   string
	Regexes  []types.Regex
	Expected string
	IsQuery  bool
}

func TestReplaceRegexPatterns(t *testing.T) {
	testCases := map[string]TestCaseReplaceRegexPatterns{
		"matchWithPlaceholder": {
			String:   "/v1beta/properties/375748157:runReport?user=paid&status=new&app=1",
			Regexes:  []types.Regex{{Pattern: "(\\d+):runReport", Placeholder: "{propertyId}:runReport"}},
			Expected: "/v1beta/properties/{propertyId}:runReport?user=paid&status=new&app=1",
			IsQuery:  false,
		},
		"matchWithoutPlaceholder": {
			String:   "/v1beta/properties/375748157:runReport?user=paid&status=new&app=1",
			Regexes:  []types.Regex{{Pattern: "(\\d+):runReport"}},
			Expected: "/v1beta/properties/{wannabe}?user=paid&status=new&app=1",
			IsQuery:  false,
		},
		"matchInQueryWithPlaceholder": {
			String:   "/v1beta/properties/375748157:runReport?user=paid&status=new&app=1",
			Regexes:  []types.Regex{{Pattern: "paid", Placeholder: "{placeholder}"}},
			Expected: "/v1beta/properties/375748157:runReport?user=%7Bplaceholder%7D&status=new&app=1",
			IsQuery:  true,
		},
		"noMatch": {
			String:   "/v1beta/properties/375748157:runReport?user=paid&status=new&app=1",
			Regexes:  []types.Regex{{Pattern: "\"dimensions\":\\s*\\[(.*?)\\][,}]"}},
			Expected: "/v1beta/properties/375748157:runReport?user=paid&status=new&app=1",
			IsQuery:  false,
		},
	}

	for testKey, tc := range testCases {
		processedString, _ := ReplaceRegexPatterns(tc.String, tc.Regexes, tc.IsQuery)

		if processedString != tc.Expected {
			t.Errorf("failed test case: %v, expected string: %v, actual string: %v", testKey, tc.Expected, processedString)
		}
	}
}

type TestCaseMapValuesToSingleString struct {
	QueryMap map[string][]string
	Expected map[string]string
}

func TestMapValuesToSingleString(t *testing.T) {
	testCases := map[string]TestCaseMapValuesToSingleString{
		"singleValue": {
			QueryMap: map[string][]string{
				"key1": {"value1"},
			},
			Expected: map[string]string{
				"key1": "value1",
			},
		},
		"multipleValues": {
			QueryMap: map[string][]string{
				"key1": {"value1", "value2"},
				"key2": {"value3", "value4"},
			},
			Expected: map[string]string{
				"key1": "value1,value2",
				"key2": "value3,value4",
			},
		},
	}

	for testKey, tc := range testCases {
		queryMap := MapValuesToSingleString(tc.QueryMap)

		if !reflect.DeepEqual(tc.Expected, queryMap) {
			t.Errorf("failed test case: %v, expected map: %v, actual map: %v", testKey, tc.Expected, queryMap)
		}
	}
}

func TestBuildQuery(t *testing.T) {
	query := testMap()
	rebuiltQuery := BuildQuery(query)

	expected := "?appId=test&status=test"

	if rebuiltQuery != expected {
		t.Errorf("expected query: %s, actual query: %s", expected, rebuiltQuery)
	}
}

type TestCaseFilterHeadersToInclude struct {
	Map      map[string][]string
	Include  []string
	Expected map[string]string
}

func testInitHeadersMap() map[string][]string {
	return map[string][]string{
		"Content-Type":  {"application/json"},
		"Accept":        {"test1", "test2", "test3"},
		"Authorization": {"test access token"},
		"X-test-header": {"test value"},
	}
}

func TestFilterHeadersToInclude(t *testing.T) {
	testCases := map[string]TestCaseFilterHeadersToInclude{
		"includeAllHeaders": {
			Map:     testInitHeadersMap(),
			Include: []string{"Accept", "Content-Type", "Authorization", "X-test-header"},
			Expected: map[string]string{
				"Accept":        "test1,test2,test3",
				"Authorization": "test access token",
				"Content-Type":  "application/json",
				"X-test-header": "test value",
			},
		},
		"includeTwoHeaders": {
			Map:     testInitHeadersMap(),
			Include: []string{"Content-Type", "X-test-header"},
			Expected: map[string]string{
				"Content-Type":  "application/json",
				"X-test-header": "test value",
			},
		},
		"nonExistingKey": {
			Map:     testInitHeadersMap(),
			Include: []string{"Non-Existing-Key", "Content-Type", "X-test-header"},
			Expected: map[string]string{
				"Content-Type":  "application/json",
				"X-test-header": "test value",
			},
		},
		"dontIncludeHeadersWithEmptyInclude": {
			Map:      testInitHeadersMap(),
			Include:  []string{},
			Expected: map[string]string{},
		},
	}

	for testName, tc := range testCases {
		headers := FilterHeadersToInclude(tc.Map, tc.Include)

		if !reflect.DeepEqual(tc.Expected, headers) {
			t.Errorf("failed test case: %v, expected headers: %v, actual headers: %v", testName, tc.Expected, headers)
		}
	}
}

type TestCaseHeadersMapToSlice struct {
	Map      map[string]string
	Expected []types.Header
}

func testHeadersMap() map[string]string {
	return map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "test1,test2,test3",
		"Authorization": "test access token",
		"X-test-header": "test value",
	}
}

func TestHeadersMapToSlice(t *testing.T) {
	testCases := map[string]TestCaseHeadersMapToSlice{
		"includeAllHeaders": {
			Map: testHeadersMap(),
			Expected: []types.Header{
				{Key: "Accept", Value: "test1,test2,test3"},
				{Key: "Authorization", Value: "test access token"},
				{Key: "Content-Type", Value: "application/json"},
				{Key: "X-test-header", Value: "test value"},
			},
		},
		"emptyHeadersMap": {
			Map:      make(map[string]string),
			Expected: []types.Header{},
		},
	}

	for testName, tc := range testCases {
		headers := HeadersMapToSlice(tc.Map)
		sortedSlice := SortHeaderSlice(headers)

		if testName == "emptyHeadersMap" && !(len(tc.Map) == 0 && len(sortedSlice) == 0) {
			t.Errorf("failed test case: %v, expected slice: %v, actual slice: %v", testName, tc.Expected, sortedSlice)
		}

		if testName != "emptyHeadersMap" && !reflect.DeepEqual(tc.Expected, sortedSlice) {
			t.Errorf("failed test case: %v, expected slice: %v, actual slice: %v", testName, tc.Expected, sortedSlice)
		}
	}
}

type TestCaseSortHeaderSlice struct {
	Slice    []types.Header
	Expected []types.Header
}

func TestSortHeaderSlice(t *testing.T) {
	testCases := map[string]TestCaseSortHeaderSlice{
		"nonEmptyHeaderSlice": {
			Slice: []types.Header{
				{Key: "Accept", Value: "test1,test2,test3"},
				{Key: "X-test-header", Value: "test value"},
				{Key: "Content-Type", Value: "application/json"},
				{Key: "Authorization", Value: "test access token"},
			},
			Expected: []types.Header{
				{Key: "Accept", Value: "test1,test2,test3"},
				{Key: "Authorization", Value: "test access token"},
				{Key: "Content-Type", Value: "application/json"},
				{Key: "X-test-header", Value: "test value"},
			},
		},
		"emptyHeaderSlice": {
			Slice:    []types.Header{},
			Expected: []types.Header{},
		},
	}

	for testName, tc := range testCases {
		sortedSlice := SortHeaderSlice(tc.Slice)

		if !reflect.DeepEqual(tc.Expected, sortedSlice) {
			t.Errorf("failed test case: %v, expected slice: %v, actual slice: %v", testName, tc.Expected, sortedSlice)
		}
	}
}

func TestIsIndexOutOfBounds(t *testing.T) {
	indexOutOfBounds := IsIndexOutOfBounds(testSlice(), 1)
	if indexOutOfBounds {
		t.Errorf("index out of bounds although it's not")
	}

	indexOutOfBounds = IsIndexOutOfBounds(testSlice(), 5)
	if !indexOutOfBounds {
		t.Errorf("index within the bounds although it's not")
	}
}

func TestKeyExists(t *testing.T) {
	exists := KeyExists(testMap(), "status")
	if !exists {
		t.Errorf("key doesn't exist, but it should exists")
	}

	exists = KeyExists(testMap(), "test")
	if exists {
		t.Errorf("key exists, but it shouldn't")
	}
}