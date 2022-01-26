package code

import (
	"github.com/merico-dev/lake/models/domainlayer"
)

type Ref struct {
	domainlayer.DomainEntity
	RepoId    string
	Ref       string
	CommitSha string `gorm:"type:char(40)"`
	IsDefault bool
	RefType   string `gorm:"type:varchar(20)"`
}
