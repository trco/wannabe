package actions

import (
	"reflect"
	"testing"

	"github.com/trco/wannabe/types"
)

func TestProcessQuery(t *testing.T) {
	tests := []struct {
		name     string
		queryMap map[string][]string
		config   types.Query
		want     string
		wantErr  string
	}{
		{
			name:     "with placeholder",
			queryMap: queryMap,
			config: types.Query{
				Wildcards: []types.WildcardKey{{Key: "user", Placeholder: "placeholder"}},
			},
			want:    "?app=1&status=new&user=placeholder",
			wantErr: "",
		},
		{
			name:     "without placeholder",
			queryMap: queryMap,
			config: types.Query{
				Wildcards: []types.WildcardKey{{Key: "user"}},
			},
			want:    "?app=1&status=new&user=wannabe",
			wantErr: "",
		},
		{
			name:     "with regex with placeholder",
			queryMap: queryMap,
			config: types.Query{
				Regexes: []types.Regex{{Pattern: "paid", Placeholder: "placeholder"}},
			},
			want:    "?app=1&status=new&user=placeholder",
			wantErr: "",
		},
		{
			name:     "with regex without placeholder",
			queryMap: queryMap,
			config: types.Query{
				Regexes: []types.Regex{{Pattern: "paid"}},
			},
			want:    "?app=1&status=new&user=wannabe",
			wantErr: "",
		},
		{
			name:     "empty string",
			queryMap: make(map[string][]string),
			config: types.Query{
				Wildcards: []types.WildcardKey{{Key: "user"}},
			},
			want:    "",
			wantErr: "",
		},
		{
			name:     "invalid regex",
			queryMap: queryMap,
			config: types.Query{
				Regexes: []types.Regex{{Pattern: "(?P<foo"}},
			},
			want:    "",
			wantErr: "ProcessQuery: failed compiling regex: error parsing regexp: invalid named capture: `(?P<foo`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProcessQuery(tt.queryMap, tt.config)

			if (err != nil) && err.Error() != tt.wantErr {
				t.Errorf("ProcessQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessQuery() = %v, want %v", got, tt.want)
			}
		})
	}

}

var queryMap = map[string][]string{
	"status": {"new"},
	"user":   {"paid"},
	"app":    {"1"},
}
