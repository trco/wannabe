package actions

import (
	"reflect"
	"testing"
	"wannabe/types"
)

type TestCaseProcessPath struct {
	Path     string
	Config   types.Path
	Expected string
}

func TestProcessPath(t *testing.T) {
	zero := 0

	tests := []struct {
		name    string
		path    string
		config  types.Path
		want    string
		wantErr string
	}{
		{
			name: "with placeholder",
			path: "/test1/test2/123456:test",
			config: types.Path{
				Wildcards: []types.WildcardIndex{{Index: &zero, Placeholder: "{placeholder}"}},
			},
			want:    "/{placeholder}/test2/123456:test",
			wantErr: "",
		},
		{
			name: "without placeholder",
			path: "/test1/test2/123456:test",
			config: types.Path{
				Wildcards: []types.WildcardIndex{{Index: &zero}},
			},
			want:    "/{wannabe}/test2/123456:test",
			wantErr: "",
		},
		{
			name: "with regex",
			path: "/test1/test2/123456:test",
			config: types.Path{
				Regexes: []types.Regex{{Pattern: "(\\d+):test", Placeholder: "{id}:test"}},
			},
			want:    "/test1/test2/{id}:test",
			wantErr: "",
		},
		{
			name: "empty string",
			path: "",
			config: types.Path{
				Wildcards: []types.WildcardIndex{{Index: &zero}},
			},
			want:    "",
			wantErr: "",
		},
		{
			name: "invalid regex",
			path: "/test1/test2/123456:test",
			config: types.Path{
				Regexes: []types.Regex{{Pattern: "(?P<foo"}},
			},
			want:    "",
			wantErr: "ProcessPath: failed compiling regex: error parsing regexp: invalid named capture: `(?P<foo`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProcessPath(tt.path, tt.config)

			if (err != nil) && err.Error() != tt.wantErr {
				t.Errorf("ProcessPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
