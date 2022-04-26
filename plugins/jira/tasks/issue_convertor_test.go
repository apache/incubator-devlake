package tasks

import "testing"

func Test_convertURL(t *testing.T) {
	type args struct {
		api      string
		issueKey string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"",
			args{"https://merico.atlassian.net/rest/agile/1.0/issue/10458", "EE-8194"},
			"https://merico.atlassian.net/browse/EE-8194",
		},
		{
			"",
			args{"http://8.142.68.162:8080/rest/agile/1.0/issue/10003", "TEST-4"},
			"http://8.142.68.162:8080/browse/TEST-4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertURL(tt.args.api, tt.args.issueKey); got != tt.want {
				t.Errorf("convertURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
