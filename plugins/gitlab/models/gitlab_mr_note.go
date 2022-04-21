package models

import (
	"time"

	"github.com/merico-dev/lake/models/common"
)

type GitlabMergeRequestNote struct {
	GitlabId        int    `gorm:"primaryKey"`
	MergeRequestId  int    `gorm:"index"`
	MergeRequestIid int    `gorm:"comment:Used in API requests ex. /api/merge_requests/<THIS_IID>"`
	NoteableType    string `gorm:"type:varchar(100)"`
	AuthorUsername  string `gorm:"type:varchar(255)"`
	Body            string
	GitlabCreatedAt time.Time
	Confidential    bool
	Resolvable      bool `gorm:"comment:Is or is not review comment"`
	System          bool `gorm:"comment:Is or is not auto-generated vs. human generated"`

	common.NoPKModel
}

func (GitlabMergeRequestNote) TableName() string {
	return "_tool_gitlab_merge_request_notes"
}
