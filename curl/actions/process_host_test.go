package actions

import (
	"reflect"
	"testing"
	"wannabe/config"
)

type TestCaseProcessHost struct {
	Host     string
	Config   config.Host
	Expected string
}

func TestProcessHost(t *testing.T) {
	zero := 0

	testCases := map[string]TestCaseProcessHost{
		"withHttp": {
			Host: "http://analyticsdata.googleapis.com",
			Config: config.Host{
				Wildcards: []config.WildcardIndex{{Index: &zero, Placeholder: "{placeholder}"}},
			},
			Expected: "{placeholder}.googleapis.com",
		},
		"withHttps": {
			Host: "https://analyticsdata.googleapis.com",
			Config: config.Host{
				Wildcards: []config.WildcardIndex{{Index: &zero, Placeholder: "{placeholder}"}},
			},
			Expected: "{placeholder}.googleapis.com",
		},
		"withoutPlaceholder": {
			Host: "https://analyticsdata.googleapis.com",
			Config: config.Host{
				Wildcards: []config.WildcardIndex{{Index: &zero}},
			},
			Expected: "{wannabe}.googleapis.com",
		},
		"withRegex": {
			Host: "https://analyticsdata.googleapis.com",
			Config: config.Host{
				Wildcards: []config.WildcardIndex{{Index: &zero, Placeholder: "{placeholder}"}},
				Regexes:   []config.Regex{{Pattern: "googleapis", Placeholder: "regexPlaceholder"}},
			},
			Expected: "{placeholder}.regexPlaceholder.com",
		},
	}

	for testKey, tc := range testCases {
		processedHost, _ := ProcessHost(tc.Host, tc.Config)

		if !reflect.DeepEqual(tc.Expected, processedHost) {
			t.Errorf("Failed test case: %v, Expected: %v, Actual: %v", testKey, tc.Expected, processedHost)
		}
	}
}
