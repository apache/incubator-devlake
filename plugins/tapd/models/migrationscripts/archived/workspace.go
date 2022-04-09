package archived

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type TapdWorkspace struct {
	SourceId    uint64 `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID          uint64 `gorm:"primaryKey;type:BIGINT(100)" json:"id"`
	Name        string `gorm:"type:varchar(255)"`
	PrettyName  string `gorm:"type:varchar(255)"`
	Category    string `gorm:"type:varchar(255)"`
	Status      string `gorm:"type:varchar(255)"`
	Description string
	BeginDate   *time.Time `json:"begin_date"`
	EndDate     *time.Time `json:"end_date"`
	ExternalOn  string     `gorm:"type:varchar(255)"`
	Creator     string     `gorm:"type:varchar(255)"`
	Created     *time.Time `json:"created"`
	common.NoPKModel
}

type TapdWorkSpaceIssue struct {
	SourceId    uint64 `gorm:"primaryKey"`
	WorkspaceId uint64 `gorm:"primaryKey"`
	IssueId     uint64 `gorm:"primaryKey"`
	common.NoPKModel
}

func (TapdWorkspace) TableName() string {
	return "_tool_tapd_workspaces"
}

func (TapdWorkSpaceIssue) TableName() string {
	return "_tool_tapd_workspace_issues"
}
