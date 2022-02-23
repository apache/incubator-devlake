package models

import (
	"github.com/merico-dev/lake/models/common"
)

type GitlabTag struct {
	Name               string `gorm:"primaryKey;type:char(60)"`
	Message            string
	Target             string
	Protected          bool
	ReleaseDescription string
	common.NoPKModel
}
