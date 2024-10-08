package hash

import (
	"reflect"
	"testing"

	"github.com/trco/wannabe/internal/config"
)

func TestProcessPath(t *testing.T) {
	zero := 0

	tests := []struct {
		name    string
		path    string
		config  config.Path
		want    string
		wantErr string
	}{
		{
			name: "with placeholder",
			path: "/test1/test2/123456:test",
			config: config.Path{
				Wildcards: []config.WildcardIndex{{Index: &zero, Placeholder: "placeholder"}},
			},
			want:    "/placeholder/test2/123456:test",
			wantErr: "",
		},
		{
			name: "without placeholder",
			path: "/test1/test2/123456:test",
			config: config.Path{
				Wildcards: []config.WildcardIndex{{Index: &zero}},
			},
			want:    "/wannabe/test2/123456:test",
			wantErr: "",
		},
		{
			name: "with regex",
			path: "/test1/test2/123456:test",
			config: config.Path{
				Regexes: []config.Regex{{Pattern: "(\\d+):test", Placeholder: "{id}:test"}},
			},
			want:    "/test1/test2/{id}:test",
			wantErr: "",
		},
		{
			name: "empty string",
			path: "",
			config: config.Path{
				Wildcards: []config.WildcardIndex{{Index: &zero}},
			},
			want:    "",
			wantErr: "",
		},
		{
			name: "invalid regex",
			path: "/test1/test2/123456:test",
			config: config.Path{
				Regexes: []config.Regex{{Pattern: "(?P<foo"}},
			},
			want:    "",
			wantErr: "processPath: failed compiling regex: error parsing regexp: invalid named capture: `(?P<foo`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processPath(tt.path, tt.config)

			if (err != nil) && err.Error() != tt.wantErr {
				t.Errorf("processPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
