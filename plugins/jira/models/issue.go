package models

import "github.com/merico-dev/lake/api/models"

type Issue struct {
	models.Model
	JiraId string
	Key    string
}
