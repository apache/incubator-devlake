package archived

import (
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type TapdWorkspace struct {
	ConnectionId uint64          `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ID           uint64          `gorm:"primaryKey;type:BIGINT" json:"id,string"`
	Name         string          `gorm:"type:varchar(255)" json:"name"`
	PrettyName   string          `gorm:"type:varchar(255)" json:"pretty_name"`
	Category     string          `gorm:"type:varchar(255)" json:"category"`
	Status       string          `gorm:"type:varchar(255)" json:"status"`
	Description  string          `json:"description"`
	BeginDate    *helper.CSTTime `json:"begin_date"`
	EndDate      *helper.CSTTime `json:"end_date"`
	ExternalOn   string          `gorm:"type:varchar(255)" json:"external_on"`
	Creator      string          `gorm:"type:varchar(255)" json:"creator"`
	Created      *helper.CSTTime `json:"created"`
	common.NoPKModel
}

type TapdWorkSpaceIssue struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	WorkspaceID  uint64 `gorm:"primaryKey"`
	IssueId      uint64 `gorm:"primaryKey"`
	common.NoPKModel
}

func (TapdWorkspace) TableName() string {
	return "_tool_tapd_workspaces"
}

func (TapdWorkSpaceIssue) TableName() string {
	return "_tool_tapd_workspace_issues"
}
