package archived

import (
	"time"
)

type Changelog struct {
	DomainEntity
	IssueId     string `gorm:"index;type:varchar(255)"`
	AuthorId    string `gorm:"type:varchar(255)"`
	AuthorName  string `gorm:"type:varchar(255)"`
	FieldId     string `gorm:"type:varchar(255)"`
	FieldName   string `gorm:"type:varchar(255)"`
	From        string
	To          string
	CreatedDate time.Time
}
