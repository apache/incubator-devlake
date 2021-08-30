package models

import (
	"github.com/merico-dev/lake/models"
)

type JiraBoard struct {
	models.Model
	ProjectId uint
	Name      string
	Self      string
	Type      string
}
