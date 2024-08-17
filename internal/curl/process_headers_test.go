package curl

import (
	"reflect"
	"testing"

	"github.com/trco/wannabe/internal/config"
)

var headersMap = map[string][]string{
	"Content-Type":  {"application/json"},
	"Accept":        {"test1", "test2", "test3"},
	"Authorization": {"test access token"},
	"X-test-header": {"test value"},
}

func TestProcessHeaders(t *testing.T) {
	tests := []struct {
		name       string
		headersMap map[string][]string
		config     config.Headers
		want       []header
	}{
		{
			name:       "include all headers",
			headersMap: headersMap,
			config: config.Headers{
				Include:   []string{"Content-Type", "Authorization", "Accept", "X-test-header"},
				Wildcards: []config.WildcardKey{},
			},
			want: []header{
				{key: "Accept", value: "test1,test2,test3"},
				{key: "Authorization", value: "test access token"},
				{key: "Content-Type", value: "application/json"},
				{key: "X-test-header", value: "test value"},
			},
		},
		{
			name:       "with placeholder and with two headers",
			headersMap: headersMap,
			config: config.Headers{
				Include:   []string{"Content-Type", "Authorization"},
				Wildcards: []config.WildcardKey{{Key: "Authorization", Placeholder: "placeholder"}},
			},
			want: []header{
				{key: "Authorization", value: "placeholder"},
				{key: "Content-Type", value: "application/json"},
			},
		},
		{
			name:       "without placeholder and with two headers",
			headersMap: headersMap,
			config: config.Headers{
				Include:   []string{"Content-Type", "Authorization"},
				Wildcards: []config.WildcardKey{{Key: "Authorization"}},
			},
			want: []header{
				{key: "Authorization", value: "wannabe"},
				{key: "Content-Type", value: "application/json"},
			},
		},

		{
			name:       "empty headers map",
			headersMap: make(map[string][]string),
			config: config.Headers{
				Include:   []string{"Content-Type", "Authorization"},
				Wildcards: []config.WildcardKey{{Key: "Authorization"}},
			},
			want: []header{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processHeaders(tt.headersMap, tt.config)

			if tt.name == "empty headers map" {
				if len(tt.headersMap) != 0 || len(got) != 0 {
					t.Errorf("processHeaders() = %v, want %v", got, tt.want)
				}
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
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
			setPlaceholderByKey(tt.testMap, tt.wildcard)

			if !reflect.DeepEqual(tt.testMap, tt.want) {
				t.Errorf("setPlaceholderByKey() = %v, want %v", tt.testMap, tt.want)
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
			got := filterHeadersToInclude(tt.headersMap, tt.include)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterHeadersToInclude() = %v, want %v", got, tt.want)
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
		want       []header
	}{
		{
			name:       "include all headers",
			headersMap: testHeadersMap,
			want: []header{
				{key: "Accept", value: "test1,test2,test3"},
				{key: "Authorization", value: "test access token"},
				{key: "Content-Type", value: "application/json"},
				{key: "X-test-header", value: "test value"},
			},
		},
		{
			name:       "empty headers map",
			headersMap: make(map[string]string),
			want:       []header{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := headersMapToSlice(tt.headersMap)
			got := sortHeaderSlice(headers)

			if tt.name == "empty headers map" {
				if len(tt.headersMap) != 0 || len(got) != 0 {
					t.Errorf("headersMapToSlice() + sortHeaderSlice() = %v, want %v", got, tt.want)
				}
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("headersMapToSlice() + sortHeaderSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestSortHeaderSlice(t *testing.T) {
	tests := []struct {
		name  string
		slice []header
		want  []header
	}{
		{
			name: "non empty header slice",
			slice: []header{
				{key: "Accept", value: "test1,test2,test3"},
				{key: "X-test-header", value: "test value"},
				{key: "Content-Type", value: "application/json"},
				{key: "Authorization", value: "test access token"},
			},
			want: []header{
				{key: "Accept", value: "test1,test2,test3"},
				{key: "Authorization", value: "test access token"},
				{key: "Content-Type", value: "application/json"},
				{key: "X-test-header", value: "test value"},
			},
		},
		{
			name:  "empty header slice",
			slice: []header{},
			want:  []header{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortHeaderSlice(tt.slice)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortHeaderSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
