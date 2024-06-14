package actions

import (
	"reflect"
	"testing"
)

func TestFilterHeaders(t *testing.T) {
	testCases := map[string]struct {
		headers map[string][]string
		exclude []string
		want    map[string][]string
	}{
		"empty headers": {
			headers: map[string][]string{},
			exclude: []string{"Content-Type"},
			want:    map[string][]string{},
		},
		"empty exclude": {
			headers: map[string][]string{"Content-Type": {"application/json"}},
			exclude: []string{},
			want:    map[string][]string{"Content-Type": {"application/json"}},
		},
		"exclude some headers": {
			headers: map[string][]string{"Content-Type": {"application/json"}, "Authorization": {"Bearer token"}},
			exclude: []string{"Authorization"},
			want:    map[string][]string{"Content-Type": {"application/json"}},
		},
	}

	for testKey, tc := range testCases {
		t.Run(testKey, func(t *testing.T) {
			got := FilterHeaders(tc.headers, tc.exclude)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("FilterHeaders() = %v, want %v", got, tc.want)
			}
		})
	}
}
