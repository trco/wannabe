package actions

import (
	"reflect"
	"testing"
	"wannabe/types"
)

func TestProcessHost(t *testing.T) {
	zero := 0

	tests := []struct {
		name    string
		host    string
		config  types.Host
		want    string
		wantErr string
	}{
		{
			name: "with http",
			host: "http://test1.test2.com",
			config: types.Host{
				Wildcards: []types.WildcardIndex{{Index: &zero, Placeholder: "{placeholder}"}},
			},
			want:    "{placeholder}.test2.com",
			wantErr: "",
		},
		{
			name: "with https",
			host: "https://test1.test2.com",
			config: types.Host{
				Wildcards: []types.WildcardIndex{{Index: &zero, Placeholder: "{placeholder}"}},
			},
			want:    "{placeholder}.test2.com",
			wantErr: "",
		},
		{
			name: "without placeholder",
			host: "https://test1.test2.com",
			config: types.Host{
				Wildcards: []types.WildcardIndex{{Index: &zero}},
			},
			want:    "{wannabe}.test2.com",
			wantErr: "",
		},
		{
			name: "with regex",
			host: "https://test1.test2.com",
			config: types.Host{
				Wildcards: []types.WildcardIndex{{Index: &zero, Placeholder: "{placeholder}"}},
				Regexes:   []types.Regex{{Pattern: "test2", Placeholder: "regexPlaceholder"}},
			},
			want:    "{placeholder}.regexPlaceholder.com",
			wantErr: "",
		},
		{
			name: "invalid regex",
			host: "https://test1.test2.com",
			config: types.Host{
				Regexes: []types.Regex{{Pattern: "(?P<foo", Placeholder: "regexPlaceholder"}},
			},
			want:    "",
			wantErr: "ProcessHost: failed compiling regex: error parsing regexp: invalid named capture: `(?P<foo`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProcessHost(tt.host, tt.config)

			if (err != nil) && err.Error() != tt.wantErr {
				t.Errorf("ProcessHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessHost() = %v, want %v", got, tt.want)
			}
		})
	}
}
