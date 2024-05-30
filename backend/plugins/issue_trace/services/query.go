package services

import (
	jiraModel "github.com/apache/incubator-devlake/plugins/jira/models"
	tapdModel "github.com/apache/incubator-devlake/plugins/tapd/models"
)

func GetTicketBoardModel(plugin string) interface{} {
	switch plugin {
	case "jira":
		return &jiraModel.JiraBoard{}
	case "tapd":
		return &tapdModel.TapdWorkspace{}
	default:
		return nil
	}
}
