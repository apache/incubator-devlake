package archived

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type TapdIssueCommit struct {
	SourceId        uint64     `gorm:"primaryKey"`
	ID              uint64     `gorm:"primaryKey;type:BIGINT(100)" json:"id"`
	WorkspaceId     uint64     `json:"workspace_id"`
	UserID          string     `gorm:"type:varchar(255)"`
	UserName        string     `gorm:"type:varchar(255)"`
	HookUserName    string     `gorm:"type:varchar(255)"`
	CommitID        string     `gorm:"type:varchar(255)"`
	Message         string     `gorm:"type:varchar(255)"`
	Path            string     `gorm:"type:varchar(255)"`
	WebURL          string     `gorm:"type:varchar(255)"`
	HookProjectName string     `gorm:"type:varchar(255)"`
	CommitTime      *time.Time `json:"commit_time"`
	Created         *time.Time `json:"created"`
	Ref             string     `gorm:"type:varchar(255)"`
	RefStatus       string     `gorm:"type:varchar(255)"`
	GitEnv          string     `gorm:"type:varchar(255)"`
	FileCommit      string     `gorm:"type:varchar(255)"`
	IssueId         uint64
	IssueType       string `gorm:"type:varchar(255)"`
	common.NoPKModel
}

func (TapdIssueCommit) TableName() string {
	return "_tool_tapd_issue_commits"
}
