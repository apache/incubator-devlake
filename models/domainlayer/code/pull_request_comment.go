package code

import (
	"github.com/merico-dev/lake/models/domainlayer"
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
