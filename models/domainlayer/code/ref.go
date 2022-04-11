package code

import (
	"github.com/merico-dev/lake/models/domainlayer"
	"time"
)

type Ref struct {
	domainlayer.DomainEntity
	RepoId      string `gorm:"type:varchar(255)"`
	Name        string `gorm:"type:varchar(255)"`
	CommitSha   string `gorm:"type:char(40)"`
	IsDefault   bool
	RefType     string `gorm:"type:varchar(255)"`
	CreatedDate *time.Time
}
