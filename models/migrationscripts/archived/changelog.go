package archived

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

type Changelog struct {
	domainlayer.DomainEntity
	IssueId     string `gorm:"index"`
	AuthorId    string `gorm:"type:char(255)"`
	AuthorName  string `gorm:"type:char(255)"`
	FieldId     string `gorm:"type:char(255)"`
	FieldName   string `gorm:"type:char(255)"`
	From        string
	To          string
	CreatedDate time.Time
}
