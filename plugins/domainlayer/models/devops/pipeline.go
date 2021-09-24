package devops

import (
	"database/sql"
	"time"

	"github.com/merico-dev/lake/plugins/domainlayer/models/base"
)

type Pipeline struct {
	base.DomainEntity
	RepoId       uint64
	CommitId     uint64
	Status       string
	Duration     int
	StartedDate  time.Time
	FinishedDate sql.NullTime
}
