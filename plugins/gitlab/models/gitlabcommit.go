package models

import "github.com/merico-dev/lake/models"

type GitlabCommit struct {
	models.Model
	Title   string
	Message string
}
