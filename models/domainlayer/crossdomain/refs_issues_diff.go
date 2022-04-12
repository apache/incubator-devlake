package crossdomain

import "github.com/merico-dev/lake/models/common"

type RefsIssuesDiffs struct {
	NewRefId        string `gorm:"type:varchar(255)"`
	OldRefId        string `gorm:"type:varchar(255)"`
	NewRefCommitSha string `gorm:"type:char(40)"`
	OldRefCommitSha string `gorm:"type:char(40)"`
	IssueNumber     string `gorm:"type:varchar(255)"`
	IssueId         string `gorm:";type:varchar(255)"`
	common.NoPKModel
}
