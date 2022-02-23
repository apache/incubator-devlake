package models

import (
	"github.com/merico-dev/lake/models/common"
)

type GitlabTag struct {
	common.Model
	Name               string `gorm:"type:varchar(255)"`
	ProjectId          int `gorm:"index"`
	Message            string
	Target             string
	Protected          bool
	ReleaseDescription string
}
