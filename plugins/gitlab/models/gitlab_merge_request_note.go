package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type GitlabMergeRequestNote struct {
	GitlabId        int `gorm:"primaryKey"`
	MergeRequestId  int `gorm:"index"`
	MergeRequestIid int `gorm:"comment:Used in API requests ex. /api/merge_requests/<THIS_IID>"`
	NoteableType    string
	AuthorUsername  string
	Body            string
	GitlabCreatedAt time.Time
	Confidential    bool
	Resolvable      bool `gorm:"comment:Is or is not review comment"`
	System          bool `gorm:"comment:Is or is not auto-generated vs. human generated"`

	common.NoPKModel
}
