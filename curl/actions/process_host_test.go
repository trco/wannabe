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
		"invalidRegex": {
			Host: "https://analyticsdata.googleapis.com",
			Config: config.Host{
				Regexes: []config.Regex{{Pattern: "(?P<foo", Placeholder: "regexPlaceholder"}},
			},
			Expected: "",
		},
	}

	for testKey, tc := range testCases {
		processedHost, err := ProcessHost(tc.Host, tc.Config)

		if testKey == "invalidRegex" && err != nil {
			expectedErr := "ProcessHost: failed compiling regex: error parsing regexp: invalid named capture: `(?P<foo`"

			if err.Error() != expectedErr {
				t.Errorf("expected error: %s, actual error: %s", expectedErr, err.Error())
			}

			continue
		}

		if !reflect.DeepEqual(tc.Expected, processedHost) {
			t.Errorf("failed test case: %v, expected host: %v, actual host: %v", testKey, tc.Expected, processedHost)
		}
	}
}
