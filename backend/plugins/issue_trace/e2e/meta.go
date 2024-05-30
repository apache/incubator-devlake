package e2e

import "github.com/apache/incubator-devlake/plugins/issue_trace/tasks"

var TaskData = &tasks.TaskData{
	Options: tasks.Options{
		Plugin:       "jira",
		ConnectionId: 2,
		BoardId:      8,
	},
	BoardId: "jira:JiraBoard:2:8",
}
