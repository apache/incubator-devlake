package archived

import "github.com/merico-dev/lake/plugins/tapd/models"

type TapdIssue struct {
	SourceId models.Uint64s `gorm:"primaryKey"`
	ID       models.Uint64s `gorm:"primaryKey;type:BIGINT(100)" json:"id"`
}

func (TapdIssue) TableName() string {
	return "_tool_tapd_issues"
}
