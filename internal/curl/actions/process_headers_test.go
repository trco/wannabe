package actions

import (
	"reflect"
	"testing"

	"github.com/trco/wannabe/types"
)

func TestProcessHeaders(t *testing.T) {
	tests := []struct {
		name       string
		headersMap map[string][]string
		config     types.Headers
		want       []types.Header
	}{
		{
			name:       "include all headers",
			headersMap: headersMap,
			config: types.Headers{
				Include:   []string{"Content-Type", "Authorization", "Accept", "X-test-header"},
				Wildcards: []types.WildcardKey{},
			},
			want: []types.Header{
				{Key: "Accept", Value: "test1,test2,test3"},
				{Key: "Authorization", Value: "test access token"},
				{Key: "Content-Type", Value: "application/json"},
				{Key: "X-test-header", Value: "test value"},
			},
		},
		{
			name:       "with placeholder and with two headers",
			headersMap: headersMap,
			config: types.Headers{
				Include:   []string{"Content-Type", "Authorization"},
				Wildcards: []types.WildcardKey{{Key: "Authorization", Placeholder: "placeholder"}},
			},
			want: []types.Header{
				{Key: "Authorization", Value: "placeholder"},
				{Key: "Content-Type", Value: "application/json"},
			},
		},
		{
			name:       "without placeholder and with two headers",
			headersMap: headersMap,
			config: types.Headers{
				Include:   []string{"Content-Type", "Authorization"},
				Wildcards: []types.WildcardKey{{Key: "Authorization"}},
			},
			want: []types.Header{
				{Key: "Authorization", Value: "wannabe"},
				{Key: "Content-Type", Value: "application/json"},
			},
		},

		{
			name:       "empty headers map",
			headersMap: make(map[string][]string),
			config: types.Headers{
				Include:   []string{"Content-Type", "Authorization"},
				Wildcards: []types.WildcardKey{{Key: "Authorization"}},
			},
			want: []types.Header{},
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
