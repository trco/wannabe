package utils

import (
	"reflect"
	"testing"
	"wannabe/types"
)

func testSlice() []string {
	return []string{"analyticsdata", "googleapis", "com"}
}

func TestSetWildcardsByIndex(t *testing.T) {
	zero := 0
	one := 1
	five := 5

	tests := []struct {
		name      string
		slice     []string
		wildcards []types.WildcardIndex
		want      []string
	}{
		{
			name:      "with placeholder",
			slice:     testSlice(),
			wildcards: []types.WildcardIndex{{Index: &zero, Placeholder: "placeholder"}},
			want:      []string{"placeholder", "googleapis", "com"},
		},
		{
			name:      "without placeholder",
			slice:     testSlice(),
			wildcards: []types.WildcardIndex{{Index: &zero}},
			want:      []string{"wannabe", "googleapis", "com"},
		},
		{
			name:      "with and without placeholder",
			slice:     testSlice(),
			wildcards: []types.WildcardIndex{{Index: &zero, Placeholder: "placeholder"}, {Index: &one}},
			want:      []string{"placeholder", "wannabe", "com"},
		},
		{
			name:      "index out of bounds",
			slice:     testSlice(),
			wildcards: []types.WildcardIndex{{Index: &five}},
			want:      testSlice(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetWildcardsByIndex(tt.slice, tt.wildcards)

			if !reflect.DeepEqual(tt.slice, tt.want) {
				t.Errorf("SetWildcardsByIndex() = %v, want %v", tt.slice, tt.want)
			}
		})
	}
}

func testMap() map[string]string {
	return map[string]string{
		"status": "test",
		"appId":  "test",
	}
}

func TestSetWildcardsByKey(t *testing.T) {
	tests := []struct {
		name      string
		testMap   map[string]string
		wildcards []types.WildcardKey
		want      map[string]string
	}{
		{
			name:      "with placeholder",
			testMap:   testMap(),
			wildcards: []types.WildcardKey{{Key: "status", Placeholder: "placeholder"}},
			want: map[string]string{
				"status": "placeholder",
				"appId":  "test",
			},
		},
		{
			name:      "without placeholder",
			testMap:   testMap(),
			wildcards: []types.WildcardKey{{Key: "status"}},
			want: map[string]string{
				"status": "wannabe",
				"appId":  "test",
			},
		},
		{
			name:      "with and without placeholder",
			testMap:   testMap(),
			wildcards: []types.WildcardKey{{Key: "status", Placeholder: "placeholder"}, {Key: "appId"}},
			want: map[string]string{
				"status": "placeholder",
				"appId":  "wannabe",
			},
		},
		{
			name:      "non existing key",
			testMap:   testMap(),
			wildcards: []types.WildcardKey{{Key: "nonExistingKey"}},
			want:      testMap(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetWildcardsByKey(tt.testMap, tt.wildcards)

			if !reflect.DeepEqual(tt.testMap, tt.want) {
				t.Errorf("SetWildcardsByKey() = %v, want %v", tt.testMap, tt.want)
			}
		})
	}
}

func TestSetPlaceholderByIndex(t *testing.T) {
	zero := 0

	tests := []struct {
		name      string
		slice     []string
		wildcards types.WildcardIndex
		want      []string
	}{
		{
			name:      "with placeholder",
			slice:     testSlice(),
			wildcards: types.WildcardIndex{Index: &zero, Placeholder: "placeholder"},
			want:      []string{"placeholder", "googleapis", "com"},
		},
		{
			name:      "without placeholder",
			slice:     testSlice(),
			wildcards: types.WildcardIndex{Index: &zero},
			want:      []string{"wannabe", "googleapis", "com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetPlaceholderByIndex(tt.slice, tt.wildcards)

			if !reflect.DeepEqual(tt.slice, tt.want) {
				t.Errorf("SetPlaceholderByIndex() = %v, want %v", tt.slice, tt.want)
			}
		})
	}
}

func TestSetPlaceholderByKey(t *testing.T) {
	tests := []struct {
		name     string
		testMap  map[string]string
		wildcard types.WildcardKey
		want     map[string]string
	}{
		{
			name:     "with placeholder",
			testMap:  testMap(),
			wildcard: types.WildcardKey{Key: "status", Placeholder: "placeholder"},
			want: map[string]string{
				"status": "placeholder",
				"appId":  "test",
			},
		},
		{
			name:     "without placeholder",
			testMap:  testMap(),
			wildcard: types.WildcardKey{Key: "status"},
			want: map[string]string{
				"status": "wannabe",
				"appId":  "test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetPlaceholderByKey(tt.testMap, tt.wildcard)

			if !reflect.DeepEqual(tt.testMap, tt.want) {
				t.Errorf("SetPlaceholderByKey() = %v, want %v", tt.testMap, tt.want)
			}
		})
	}
}

func TestReplaceRegexPatterns(t *testing.T) {
	tests := []struct {
		name       string
		testString string
		regexes    []types.Regex
		want       string
		isQuery    bool
	}{
		{
			name:       "match with placeholder",
			testString: "/v1beta/properties/375748157:runReport?user=paid&status=new&app=1",
			regexes:    []types.Regex{{Pattern: "(\\d+):runReport", Placeholder: "{propertyId}:runReport"}},
			want:       "/v1beta/properties/{propertyId}:runReport?user=paid&status=new&app=1",
			isQuery:    false,
		},
		{
			name:       "match without placeholder",
			testString: "/v1beta/properties/375748157:runReport?user=paid&status=new&app=1",
			regexes:    []types.Regex{{Pattern: "(\\d+):runReport"}},
			want:       "/v1beta/properties/wannabe?user=paid&status=new&app=1",
			isQuery:    false,
		},
		{
			name:       "match in query with placeholder",
			testString: "/v1beta/properties/375748157:runReport?user=paid&status=new&app=1",
			regexes:    []types.Regex{{Pattern: "paid", Placeholder: "placeholder"}},
			want:       "/v1beta/properties/375748157:runReport?user=%7Bplaceholder%7D&status=new&app=1",
			isQuery:    true,
		},
		{
			name:       "no match",
			testString: "/v1beta/properties/375748157:runReport?user=paid&status=new&app=1",
			regexes:    []types.Regex{{Pattern: "\"dimensions\":\\s*\\[(.*?)\\][,}]"}},
			want:       "/v1beta/properties/375748157:runReport?user=paid&status=new&app=1",
			isQuery:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := ReplaceRegexPatterns(tt.testString, tt.regexes, tt.isQuery)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReplaceRegexPatterns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapValuesToSingleString(t *testing.T) {
	tests := []struct {
		name     string
		queryMap map[string][]string
		want     map[string]string
	}{
		{
			name: "single value",
			queryMap: map[string][]string{
				"key1": {"value1"},
			},
			want: map[string]string{
				"key1": "value1",
			},
		},
		{

			name: "multiple values",
			queryMap: map[string][]string{
				"key1": {"value1", "value2"},
				"key2": {"value3", "value4"},
			},
			want: map[string]string{
				"key1": "value1,value2",
				"key2": "value3,value4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapValuesToSingleString(tt.queryMap)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapValuesToSingleString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildQuery(t *testing.T) {
	t.Run("build query", func(t *testing.T) {
		query := testMap()

		want := "?appId=test&status=test"

		got := BuildQuery(query)

		if got != want {
			t.Errorf("BuildQuery() = %v, want %v", got, want)
		}
	})
}

var initHeadersMap = map[string][]string{
	"Content-Type":  {"application/json"},
	"Accept":        {"test1", "test2", "test3"},
	"Authorization": {"test access token"},
	"X-test-header": {"test value"},
}

func TestFilterHeadersToInclude(t *testing.T) {
	tests := []struct {
		name       string
		headersMap map[string][]string
		include    []string
		want       map[string]string
	}{
		{
			name:       "include all headers",
			headersMap: initHeadersMap,
			include:    []string{"Accept", "Content-Type", "Authorization", "X-test-header"},
			want: map[string]string{
				"Accept":        "test1,test2,test3",
				"Authorization": "test access token",
				"Content-Type":  "application/json",
				"X-test-header": "test value",
			},
		},
		{
			name:       "include two headers",
			headersMap: initHeadersMap,
			include:    []string{"Content-Type", "X-test-header"},
			want: map[string]string{
				"Content-Type":  "application/json",
				"X-test-header": "test value",
			},
		},
		{
			name:       "non existing key",
			headersMap: initHeadersMap,
			include:    []string{"Non-Existing-Key", "Content-Type", "X-test-header"},
			want: map[string]string{
				"Content-Type":  "application/json",
				"X-test-header": "test value",
			},
		},
		{
			name:       "don't include headers with empty include",
			headersMap: initHeadersMap,
			include:    []string{},
			want:       map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterHeadersToInclude(tt.headersMap, tt.include)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterHeadersToInclude() = %v, want %v", got, tt.want)
			}
		})
	}
}

var testHeadersMap = map[string]string{
	"Content-Type":  "application/json",
	"Accept":        "test1,test2,test3",
	"Authorization": "test access token",
	"X-test-header": "test value",
}

func TestHeadersMapToSlice(t *testing.T) {
	tests := []struct {
		name       string
		headersMap map[string]string
		want       []types.Header
	}{
		{
			name:       "include all headers",
			headersMap: testHeadersMap,
			want: []types.Header{
				{Key: "Accept", Value: "test1,test2,test3"},
				{Key: "Authorization", Value: "test access token"},
				{Key: "Content-Type", Value: "application/json"},
				{Key: "X-test-header", Value: "test value"},
			},
		},
		{
			name:       "empty headers map",
			headersMap: make(map[string]string),
			want:       []types.Header{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := HeadersMapToSlice(tt.headersMap)
			got := SortHeaderSlice(headers)

			if tt.name == "empty headers map" {
				if len(tt.headersMap) != 0 || len(got) != 0 {
					t.Errorf("HeadersMapToSlice() + SortHeaderSlice() = %v, want %v", got, tt.want)
				}
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HeadersMapToSlice() + SortHeaderSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestSortHeaderSlice(t *testing.T) {
	tests := []struct {
		name  string
		slice []types.Header
		want  []types.Header
	}{
		{
			name: "non empty header slice",
			slice: []types.Header{
				{Key: "Accept", Value: "test1,test2,test3"},
				{Key: "X-test-header", Value: "test value"},
				{Key: "Content-Type", Value: "application/json"},
				{Key: "Authorization", Value: "test access token"},
			},
			want: []types.Header{
				{Key: "Accept", Value: "test1,test2,test3"},
				{Key: "Authorization", Value: "test access token"},
				{Key: "Content-Type", Value: "application/json"},
				{Key: "X-test-header", Value: "test value"},
			},
		},
		{
			name:  "empty header slice",
			slice: []types.Header{},
			want:  []types.Header{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortHeaderSlice(tt.slice)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortHeaderSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIndexOutOfBounds(t *testing.T) {
	tests := []struct {
		name  string
		slice []string
		index int
		want  bool
	}{
		{
			name:  "index not out of bounds",
			slice: testSlice(),
			index: 1,
			want:  false,
		},
		{
			name:  "index out of bounds",
			slice: testSlice(),
			index: 5,
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsIndexOutOfBounds(tt.slice, tt.index)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsIndexOutOfBounds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyExists(t *testing.T) {
	tests := []struct {
		name    string
		testMap map[string]string
		key     string
		want    bool
	}{
		{
			name:    "key exist",
			testMap: testMap(),
			key:     "status",
			want:    true,
		},
		{
			name:    "key doesn't exists",
			testMap: testMap(),
			key:     "test",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := KeyExists(tt.testMap, tt.key)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeyExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
