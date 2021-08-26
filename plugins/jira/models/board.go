package models

import "github.com/merico-dev/lake/api/models"

type Board struct {
	models.Model
	JiraId string
	Name    string
}

func (m Board) TableName() string {
	return `jira_plugin_board`
}
