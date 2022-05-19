package ticket

import (
	"github.com/apache/incubator-devlake/models/domainlayer"
	"time"
)

type IssueComment struct {
	domainlayer.DomainEntity
	IssueId     string `gorm:"index"`
	Body        string
	UserId      string `gorm:"type:varchar(255)"`
	CreatedDate time.Time
}
