package ticket

import (
	"github.com/merico-dev/lake/models/domainlayer"
	"time"
)

type IssueComment struct {
	domainlayer.DomainEntity
	IssueId     string `gorm:"index"`
	Body        string
	UserId      string `gorm:"type:varchar(255)"`
	CreatedDate time.Time
}
