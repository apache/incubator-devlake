package models

type TapdIssue struct {
	SourceId Uint64s `gorm:"primaryKey"`
	ID       Uint64s `gorm:"primaryKey;type:BIGINT(100)" json:"id"`
}

func (TapdIssue) TableName() string {
	return "_tool_tapd_issues"
}
