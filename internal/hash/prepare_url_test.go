package hash

import "testing"

func TestPrepareUrl(t *testing.T) {
	scheme := "https"
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
			want: "https://test.com/test?test=test",
		},
		{
			name: "path is '/'",
			path: "/",
			want: "https://test.com?test=test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := prepareUrl(scheme, host, tt.path, query)

			if got != tt.want {
				t.Errorf("prepareUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
