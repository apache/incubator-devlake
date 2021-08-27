package models

import "github.com/merico-dev/lake/models"

type JiraIssue struct {
	models.Model
	JiraId string
	Key    string
}
