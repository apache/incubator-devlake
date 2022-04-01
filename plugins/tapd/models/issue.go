package models

type TapdIssue struct {
	SourceId uint64 `gorm:"primaryKey"`
	ID       uint64 `gorm:"primaryKey;type:BIGINT(100)" json:"id"`
}
