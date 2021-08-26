package models

import "github.com/merico-dev/lake/api/models"

type Issue struct {
	models.Model
	JiraId string
	Key    string
}

func (i Issue) TableName() string {
	return `jira_plugin_issue`
}
