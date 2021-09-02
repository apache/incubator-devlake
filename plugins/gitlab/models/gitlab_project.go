package models

import "github.com/merico-dev/lake/models"

type GitlabProject struct {
	GitlabId          int `gorm:"primary_key"`
	Name              string
	PathWithNamespace string
	WebUrl            string
	Visibility        string
	OpenIssuesCount   int
	StarCount         int

	models.NoPKModel
}
