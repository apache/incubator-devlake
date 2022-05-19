package code

import (
	"github.com/apache/incubator-devlake/models/domainlayer"
	"time"
)

type PullRequestComment struct {
	domainlayer.DomainEntity
	PullRequestId string `gorm:"index"`
	Body          string
	UserId        string `gorm:"type:varchar(255)"`
	CreatedDate   time.Time
	CommitSha     string `gorm:"type:varchar(255)"`
	Position      int
}
