package models

import (
	"github.com/merico-dev/lake/plugins/helper"
	"time"
)

// Please note that Issue Labels can also apply to Pull Requests.
// Pull Requests are considered Issues in GitHub.

type GithubPullRequestLabel struct {
	PullId    int    `gorm:"primaryKey;autoIncrement:false"`
	LabelName string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"index"`

	helper.RawDataOrigin
}
