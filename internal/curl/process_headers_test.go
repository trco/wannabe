package curl

import (
	"reflect"
	"testing"

	"github.com/trco/wannabe/internal/config"
)

func TestProcessHeaders(t *testing.T) {
	tests := []struct {
		name       string
		headersMap map[string][]string
		config     config.Headers
		want       []Header
	}{
		{
			name:       "include all headers",
			headersMap: headersMap,
			config: config.Headers{
				Include:   []string{"Content-Type", "Authorization", "Accept", "X-test-header"},
				Wildcards: []config.WildcardKey{},
			},
			want: []Header{
				{Key: "Accept", Value: "test1,test2,test3"},
				{Key: "Authorization", Value: "test access token"},
				{Key: "Content-Type", Value: "application/json"},
				{Key: "X-test-header", Value: "test value"},
			},
		},
		{
			name:       "with placeholder and with two headers",
			headersMap: headersMap,
			config: config.Headers{
				Include:   []string{"Content-Type", "Authorization"},
				Wildcards: []config.WildcardKey{{Key: "Authorization", Placeholder: "placeholder"}},
			},
			want: []Header{
				{Key: "Authorization", Value: "placeholder"},
				{Key: "Content-Type", Value: "application/json"},
			},
		},
		{
			name:       "without placeholder and with two headers",
			headersMap: headersMap,
			config: config.Headers{
				Include:   []string{"Content-Type", "Authorization"},
				Wildcards: []config.WildcardKey{{Key: "Authorization"}},
			},
			want: []Header{
				{Key: "Authorization", Value: "wannabe"},
				{Key: "Content-Type", Value: "application/json"},
			},
		},

		{
			name:       "empty headers map",
			headersMap: make(map[string][]string),
			config: config.Headers{
				Include:   []string{"Content-Type", "Authorization"},
				Wildcards: []config.WildcardKey{{Key: "Authorization"}},
			},
			want: []Header{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProcessHeaders(tt.headersMap, tt.config)

			if tt.name == "empty headers map" {
				if len(tt.headersMap) != 0 || len(got) != 0 {
					t.Errorf("ProcessHeaders() = %v, want %v", got, tt.want)
				}
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

var headersMap = map[string][]string{
	"Content-Type":  {"application/json"},
	"Accept":        {"test1", "test2", "test3"},
	"Authorization": {"test access token"},
	"X-test-header": {"test value"},
}

func TestSetPlaceholderByKey(t *testing.T) {
	tests := []struct {
		name     string
		testMap  map[string]string
		wildcard config.WildcardKey
		want     map[string]string
	}{
		{
			name:     "with placeholder",
			testMap:  testMap(),
			wildcard: config.WildcardKey{Key: "status", Placeholder: "placeholder"},
			want: map[string]string{
				"status": "placeholder",
				"appId":  "test",
			},
		},
		{
			name:     "without placeholder",
			testMap:  testMap(),
			wildcard: config.WildcardKey{Key: "status"},
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
		want       []Header
	}{
		{
			name:       "include all headers",
			headersMap: testHeadersMap,
			want: []Header{
				{Key: "Accept", Value: "test1,test2,test3"},
				{Key: "Authorization", Value: "test access token"},
				{Key: "Content-Type", Value: "application/json"},
				{Key: "X-test-header", Value: "test value"},
			},
		},
		{
			name:       "empty headers map",
			headersMap: make(map[string]string),
			want:       []Header{},
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
		slice []Header
		want  []Header
	}{
		{
			name: "non empty header slice",
			slice: []Header{
				{Key: "Accept", Value: "test1,test2,test3"},
				{Key: "X-test-header", Value: "test value"},
				{Key: "Content-Type", Value: "application/json"},
				{Key: "Authorization", Value: "test access token"},
			},
			want: []Header{
				{Key: "Accept", Value: "test1,test2,test3"},
				{Key: "Authorization", Value: "test access token"},
				{Key: "Content-Type", Value: "application/json"},
				{Key: "X-test-header", Value: "test value"},
			},
		},
		{
			name:  "empty header slice",
			slice: []Header{},
			want:  []Header{},
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
