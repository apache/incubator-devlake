package models

import (
	"github.com/merico-dev/lake/models/common"
)

type GitlabProject struct {
	GitlabId          int `gorm:"primaryKey"`
	Name              string
	PathWithNamespace string
	WebUrl            string
	Visibility        string
	OpenIssuesCount   int
	StarCount         int

	common.NoPKModel
}
