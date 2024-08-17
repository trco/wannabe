package curl

import (
	"reflect"
	"testing"

	"github.com/trco/wannabe/internal/config"
)

func TestProcessHost(t *testing.T) {
	zero := 0

	tests := []struct {
		name    string
		host    string
		config  config.Host
		want    string
		wantErr string
	}{
		{
			name: "with http",
			host: "http://test1.test2.com",
			config: config.Host{
				Wildcards: []config.WildcardIndex{{Index: &zero, Placeholder: "placeholder"}},
			},
			want:    "placeholder.test2.com",
			wantErr: "",
		},
		{
			name: "with https",
			host: "https://test1.test2.com",
			config: config.Host{
				Wildcards: []config.WildcardIndex{{Index: &zero, Placeholder: "placeholder"}},
			},
			want:    "placeholder.test2.com",
			wantErr: "",
		},
		{
			name: "without placeholder",
			host: "https://test1.test2.com",
			config: config.Host{
				Wildcards: []config.WildcardIndex{{Index: &zero}},
			},
			want:    "wannabe.test2.com",
			wantErr: "",
		},
		{
			name: "with regex",
			host: "https://test1.test2.com",
			config: config.Host{
				Wildcards: []config.WildcardIndex{{Index: &zero, Placeholder: "placeholder"}},
				Regexes:   []config.Regex{{Pattern: "test2", Placeholder: "regexPlaceholder"}},
			},
			want:    "placeholder.regexPlaceholder.com",
			wantErr: "",
		},
		{
			name: "invalid regex",
			host: "https://test1.test2.com",
			config: config.Host{
				Regexes: []config.Regex{{Pattern: "(?P<foo", Placeholder: "regexPlaceholder"}},
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

func testSlice() []string {
	return []string{"analyticsdata", "googleapis", "com"}
}

func TestSetWildcardsByIndex(t *testing.T) {
	zero := 0
	one := 1
	five := 5

	tests := []struct {
		name      string
		slice     []string
		wildcards []config.WildcardIndex
		want      []string
	}{
		{
			name:      "with placeholder",
			slice:     testSlice(),
			wildcards: []config.WildcardIndex{{Index: &zero, Placeholder: "placeholder"}},
			want:      []string{"placeholder", "googleapis", "com"},
		},
		{
			name:      "without placeholder",
			slice:     testSlice(),
			wildcards: []config.WildcardIndex{{Index: &zero}},
			want:      []string{"wannabe", "googleapis", "com"},
		},
		{
			name:      "with and without placeholder",
			slice:     testSlice(),
			wildcards: []config.WildcardIndex{{Index: &zero, Placeholder: "placeholder"}, {Index: &one}},
			want:      []string{"placeholder", "wannabe", "com"},
		},
		{
			name:      "index out of bounds",
			slice:     testSlice(),
			wildcards: []config.WildcardIndex{{Index: &five}},
			want:      testSlice(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetWildcardsByIndex(tt.slice, tt.wildcards)

			if !reflect.DeepEqual(tt.slice, tt.want) {
				t.Errorf("SetWildcardsByIndex() = %v, want %v", tt.slice, tt.want)
			}
		})
	}
}

func TestIsIndexOutOfBounds(t *testing.T) {
	tests := []struct {
		name  string
		slice []string
		index int
		want  bool
	}{
		{
			name:  "index not out of bounds",
			slice: testSlice(),
			index: 1,
			want:  false,
		},
		{
			name:  "index out of bounds",
			slice: testSlice(),
			index: 5,
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsIndexOutOfBounds(tt.slice, tt.index)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsIndexOutOfBounds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetPlaceholderByIndex(t *testing.T) {
	zero := 0

	tests := []struct {
		name      string
		slice     []string
		wildcards config.WildcardIndex
		want      []string
	}{
		{
			name:      "with placeholder",
			slice:     testSlice(),
			wildcards: config.WildcardIndex{Index: &zero, Placeholder: "placeholder"},
			want:      []string{"placeholder", "googleapis", "com"},
		},
		{
			name:      "without placeholder",
			slice:     testSlice(),
			wildcards: config.WildcardIndex{Index: &zero},
			want:      []string{"wannabe", "googleapis", "com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetPlaceholderByIndex(tt.slice, tt.wildcards)

			if !reflect.DeepEqual(tt.slice, tt.want) {
				t.Errorf("SetPlaceholderByIndex() = %v, want %v", tt.slice, tt.want)
			}
		})
	}
}

func TestReplaceRegexPatterns(t *testing.T) {
	tests := []struct {
		name       string
		testString string
		regexes    []config.Regex
		want       string
		isQuery    bool
	}{
		{
			name:       "match with placeholder",
			testString: "/v1beta/properties/375748157:runReport?user=paid&status=new&app=1",
			regexes:    []config.Regex{{Pattern: "(\\d+):runReport", Placeholder: "{propertyId}:runReport"}},
			want:       "/v1beta/properties/{propertyId}:runReport?user=paid&status=new&app=1",
			isQuery:    false,
		},
		{
			name:       "match without placeholder",
			testString: "/v1beta/properties/375748157:runReport?user=paid&status=new&app=1",
			regexes:    []config.Regex{{Pattern: "(\\d+):runReport"}},
			want:       "/v1beta/properties/wannabe?user=paid&status=new&app=1",
			isQuery:    false,
		},
		{
			name:       "match in query with placeholder",
			testString: "/v1beta/properties/375748157:runReport?user=paid&status=new&app=1",
			regexes:    []config.Regex{{Pattern: "paid", Placeholder: "placeholder"}},
			want:       "/v1beta/properties/375748157:runReport?user=placeholder&status=new&app=1",
			isQuery:    true,
		},
		{
			name:       "no match",
			testString: "/v1beta/properties/375748157:runReport?user=paid&status=new&app=1",
			regexes:    []config.Regex{{Pattern: "\"dimensions\":\\s*\\[(.*?)\\][,}]"}},
			want:       "/v1beta/properties/375748157:runReport?user=paid&status=new&app=1",
			isQuery:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := ReplaceRegexPatterns(tt.testString, tt.regexes, tt.isQuery)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReplaceRegexPatterns() = %v, want %v", got, tt.want)
			}
		})
	}
}
