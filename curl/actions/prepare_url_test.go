package actions

import "testing"

func TestPrepareUrl(t *testing.T) {
	host := "test.com"
	query := "?test=test"

	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "path is '/test'",
			path: "/test",
			want: "test.com/test?test=test",
		},
		{
			name: "path is '/'",
			path: "/",
			want: "test.com?test=test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PrepareUrl(host, tt.path, query)

			if got != tt.want {
				t.Errorf("PrepareUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
