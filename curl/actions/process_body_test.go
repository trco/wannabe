package actions

import (
	"reflect"
	"testing"

	"github.com/trco/wannabe/types"
)

type TestCaseProcessBody struct {
	Body     []byte
	Config   types.Body
	Expected string
}

func TestProcessBody(t *testing.T) {
	tests := []struct {
		name    string
		body    []byte
		config  types.Body
		want    string
		wantErr string
	}{
		{
			name: "with placeholder",
			body: testBody,
			config: types.Body{
				Regexes: []types.Regex{
					{Pattern: "\"dimensions\":\\s*\\[(.*?)\\][,}]", Placeholder: "\"dimensions\":\"placeholder\","},
					{Pattern: "\"startDate\":\\s*\"(.*?)\"[,}]", Placeholder: "\"startDate\":\"placeholder\""},
				},
			},
			want:    "{\"dateRanges\":[{\"endDate\":\"2023-12-31\",\"startDate\":\"placeholder\"],\"dimensions\":\"placeholder\",\"limit\":10000,\"metrics\":[{\"name\":\"sessions\"}],\"returnPropertyQuota\":true}",
			wantErr: "",
		},
		{
			name: "without placeholder",
			body: testBody,
			config: types.Body{
				Regexes: []types.Regex{
					{Pattern: "\"startDate\":\\s*\"(.*?)\"[,}\"]"},
				},
			},
			want:    "{\"dateRanges\":[{\"endDate\":\"2023-12-31\",wannabe],\"dimensions\":[{\"name\":\"date\"},{\"name\":\"source\"}],\"limit\":10000,\"metrics\":[{\"name\":\"sessions\"}],\"returnPropertyQuota\":true}",
			wantErr: "",
		},
		{
			name: "empty body",
			body: []byte{},
			config: types.Body{
				Regexes: []types.Regex{},
			},
			want:    "",
			wantErr: "",
		},
		{
			name: "invalid regex",
			body: testBody,
			config: types.Body{
				Regexes: []types.Regex{
					{Pattern: "(?P<foo"},
				},
			},
			want:    "",
			wantErr: "ProcessBody: failed compiling regex: error parsing regexp: invalid named capture: `(?P<foo`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProcessBody(tt.body, tt.config)

			if (err != nil) && err.Error() != tt.wantErr {
				t.Errorf("ProcessBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

var testBody = []byte{
	123, 10, 32, 32, 32, 32, 34, 100, 97, 116, 101, 82, 97, 110, 103, 101, 115, 34, 58, 32, 91, 10, 32, 32, 32,
	32, 32, 32, 32, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 34, 115, 116, 97, 114, 116, 68,
	97, 116, 101, 34, 58, 32, 34, 50, 48, 50, 51, 45, 49, 50, 45, 48, 49, 34, 44, 10, 32, 32, 32, 32, 32, 32, 32,
	32, 32, 32, 32, 32, 34, 101, 110, 100, 68, 97, 116, 101, 34, 58, 32, 34, 50, 48, 50, 51, 45, 49, 50, 45, 51,
	49, 34, 10, 32, 32, 32, 32, 32, 32, 32, 32, 125, 10, 32, 32, 32, 32, 93, 44, 10, 32, 32, 32, 32, 34, 109, 101,
	116, 114, 105, 99, 115, 34, 58, 32, 91, 10, 32, 32, 32, 32, 32, 32, 32, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32,
	32, 32, 32, 32, 32, 34, 110, 97, 109, 101, 34, 58, 32, 34, 115, 101, 115, 115, 105, 111, 110, 115, 34, 10, 32, 32,
	32, 32, 32, 32, 32, 32, 125, 10, 32, 32, 32, 32, 93, 44, 10, 32, 32, 32, 32, 34, 100, 105, 109, 101, 110, 115, 105,
	111, 110, 115, 34, 58, 32, 91, 123, 34, 110, 97, 109, 101, 34, 58, 32, 34, 100, 97, 116, 101, 34, 125, 44, 32, 123,
	34, 110, 97, 109, 101, 34, 58, 32, 34, 115, 111, 117, 114, 99, 101, 34, 125, 93, 44, 10, 32, 32, 32, 32, 34, 108, 105,
	109, 105, 116, 34, 58, 32, 49, 48, 48, 48, 48, 44, 10, 32, 32, 32, 32, 34, 114, 101, 116, 117, 114, 110, 80, 114, 111,
	112, 101, 114, 116, 121, 81, 117, 111, 116, 97, 34, 58, 32, 116, 114, 117, 101, 10, 125}
