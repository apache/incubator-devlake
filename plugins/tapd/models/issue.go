package models

type TapdIssue struct {
	SourceId uint64 `gorm:"primaryKey"`
	ID       uint64 `gorm:"primaryKey;type:BIGINT(100)" json:"id,string"`
}

func (TapdIssue) TableName() string {
	return "_tool_tapd_issues"
}
