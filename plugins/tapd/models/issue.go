package models

type TapdIssue struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	ID           uint64 `gorm:"primaryKey;type:BIGINT" json:"id,string"`
}

func (TapdIssue) TableName() string {
	return "_tool_tapd_issues"
}
