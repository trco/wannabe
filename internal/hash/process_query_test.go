package hash

import (
	"reflect"
	"testing"

	"github.com/trco/wannabe/internal/config"
)

var queryMap = map[string][]string{
	"status": {"new"},
	"user":   {"paid"},
	"app":    {"1"},
}

func TestProcessQuery(t *testing.T) {
	tests := []struct {
		name     string
		queryMap map[string][]string
		config   config.Query
		want     string
		wantErr  string
	}{
		{
			name:     "with placeholder",
			queryMap: queryMap,
			config: config.Query{
				Wildcards: []config.WildcardKey{{Key: "user", Placeholder: "placeholder"}},
			},
			want:    "?app=1&status=new&user=placeholder",
			wantErr: "",
		},
		{
			name:     "without placeholder",
			queryMap: queryMap,
			config: config.Query{
				Wildcards: []config.WildcardKey{{Key: "user"}},
			},
			want:    "?app=1&status=new&user=wannabe",
			wantErr: "",
		},
		{
			name:     "with regex with placeholder",
			queryMap: queryMap,
			config: config.Query{
				Regexes: []config.Regex{{Pattern: "paid", Placeholder: "placeholder"}},
			},
			want:    "?app=1&status=new&user=placeholder",
			wantErr: "",
		},
		{
			name:     "with regex without placeholder",
			queryMap: queryMap,
			config: config.Query{
				Regexes: []config.Regex{{Pattern: "paid"}},
			},
			want:    "?app=1&status=new&user=wannabe",
			wantErr: "",
		},
		{
			name:     "empty string",
			queryMap: make(map[string][]string),
			config: config.Query{
				Wildcards: []config.WildcardKey{{Key: "user"}},
			},
			want:    "",
			wantErr: "",
		},
		{
			name:     "invalid regex",
			queryMap: queryMap,
			config: config.Query{
				Regexes: []config.Regex{{Pattern: "(?P<foo"}},
			},
			want:    "",
			wantErr: "processQuery: failed compiling regex: error parsing regexp: invalid named capture: `(?P<foo`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processQuery(tt.queryMap, tt.config)

			if (err != nil) && err.Error() != tt.wantErr {
				t.Errorf("processQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processQuery() = %v, want %v", got, tt.want)
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
			got := mapValuesToSingleString(tt.queryMap)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapValuesToSingleString() = %v, want %v", got, tt.want)
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
		wildcards []config.WildcardKey
		want      map[string]string
	}{
		{
			name:      "with placeholder",
			testMap:   testMap(),
			wildcards: []config.WildcardKey{{Key: "status", Placeholder: "placeholder"}},
			want: map[string]string{
				"status": "placeholder",
				"appId":  "test",
			},
		},
		{
			name:      "without placeholder",
			testMap:   testMap(),
			wildcards: []config.WildcardKey{{Key: "status"}},
			want: map[string]string{
				"status": "wannabe",
				"appId":  "test",
			},
		},
		{
			name:      "with and without placeholder",
			testMap:   testMap(),
			wildcards: []config.WildcardKey{{Key: "status", Placeholder: "placeholder"}, {Key: "appId"}},
			want: map[string]string{
				"status": "placeholder",
				"appId":  "wannabe",
			},
		},
		{
			name:      "non existing key",
			testMap:   testMap(),
			wildcards: []config.WildcardKey{{Key: "nonExistingKey"}},
			want:      testMap(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setWildcardsByKey(tt.testMap, tt.wildcards)

			if !reflect.DeepEqual(tt.testMap, tt.want) {
				t.Errorf("setWildcardsByKey() = %v, want %v", tt.testMap, tt.want)
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
			got := keyExists(tt.testMap, tt.key)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keyExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildQuery(t *testing.T) {
	t.Run("build query", func(t *testing.T) {
		query := testMap()

		want := "?appId=test&status=test"

		got := buildQuery(query)

		if got != want {
			t.Errorf("buildQuery() = %v, want %v", got, want)
		}
	})
}
