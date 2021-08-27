package main

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/jira/models"
)

func (plugin Jira) Init() {
	err := lakeModels.Db.AutoMigrate(&models.JiraIssue{}, &models.JiraBoard{})
	if err != nil {
		panic(err)
	}
}
