package record

import (
	"reflect"
	"testing"
)

func TestFilterHeaders(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string][]string
		exclude []string
		want    map[string][]string
	}{
		{
			name:    "empty headers",
			headers: map[string][]string{},
			exclude: []string{"Content-Type"},
			want:    map[string][]string{},
		},
		{
			name:    "empty exclude",
			headers: map[string][]string{"Content-Type": {"application/json"}},
			exclude: []string{},
			want:    map[string][]string{"Content-Type": {"application/json"}},
		},
		{
			name:    "exclude some headers",
			headers: map[string][]string{"Content-Type": {"application/json"}, "Authorization": {"Bearer token"}},
			exclude: []string{"Authorization"},
			want:    map[string][]string{"Content-Type": {"application/json"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterHeaders(tt.headers, tt.exclude)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}
