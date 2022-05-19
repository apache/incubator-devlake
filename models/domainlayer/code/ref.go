package code

import (
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer"
)

type Ref struct {
	domainlayer.DomainEntity
	RepoId      string `gorm:"type:varchar(255)"`
	Name        string `gorm:"type:varchar(255)"`
	CommitSha   string `gorm:"type:varchar(40)"`
	IsDefault   bool
	RefType     string `gorm:"type:varchar(255)"`
	CreatedDate *time.Time
}
