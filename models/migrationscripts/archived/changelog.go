package archived

import (
	"time"
)

type Changelog struct {
	DomainEntity
	IssueId     string `gorm:"index"`
	AuthorId    string `gorm:"type:char(255)"`
	AuthorName  string `gorm:"type:char(255)"`
	FieldId     string `gorm:"type:char(255)"`
	FieldName   string `gorm:"type:char(255)"`
	From        string
	To          string
	CreatedDate time.Time
}
