package models

import "github.com/merico-dev/lake/api/models"

type Board struct {
	models.Model
	JiraId string
	Name    string
}
