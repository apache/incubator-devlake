package models

import "github.com/merico-dev/lake/models"

type JiraBoard struct {
	models.Model
	JiraId string
	Name   string
}
