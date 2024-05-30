package tasks

import (
	"reflect"
	"testing"

	"github.com/apache/incubator-devlake/plugins/issue_trace/models"
)

func Test_buildStatusHistoryRecords(t *testing.T) {
	type args struct {
		logs []*StatusChangeLogResult
	}

	tests := []struct {
		name string
		args args
		want []*models.IssueStatusHistory
	}{
		{
			name: "empty",
			args: args{
				logs: make([]*StatusChangeLogResult, 0),
			},
			want: make([]*models.IssueStatusHistory, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildStatusHistoryRecords(tt.args.logs, "jira:JiraBoard:1:1"); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildStatusHistoryRecords() = %v, want %v", got, tt.want)
			}
		})
	}
}
