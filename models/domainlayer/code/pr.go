package code

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

type Pr struct {
	domainlayer.DomainEntity
	RepoId      uint64 `gorm:"index"`
	State       string `gorm:"comment:open/closed or other"`
	Title       string
	Url         string
	CreatedDate time.Time
	MergedDate  *time.Time
	ClosedAt    *time.Time
}
