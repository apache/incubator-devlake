package code

import (
	"database/sql"
	"time"

	"github.com/merico-dev/lake/plugins/domainlayer/models/base"
)

type Pr struct {
	base.DomainEntity
	RepoId      uint64 `gorm:"index"`
	State       string
	Title       string
	Url         string
	CreatedDate time.Time
	MergedDate  sql.NullTime
	ClosedAt    sql.NullTime
}
