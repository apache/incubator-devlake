package models

import (
	"github.com/merico-dev/lake/models/common"
)

// Please note that Issue Labels can also apply to Pull Requests.
// Pull Requests are considered Issues in GitHub.

type GithubIssueLabel struct {
	GithubId    int `gorm:"primaryKey"`
	Name        string
	Description string
	Color       string

	common.NoPKModel
}
