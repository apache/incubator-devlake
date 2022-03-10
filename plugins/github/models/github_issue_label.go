package models

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/helper"
)

// Please note that Issue Labels can also apply to Pull Requests.
// Pull Requests are considered Issues in GitHub.

type GithubIssueLabel struct {
	IssueId   int    `gorm:"primaryKey;autoIncrement:false"`
	LabelName string `gorm:"primaryKey"`
	common.NoPKModel
	helper.RawDataOrigin
}
