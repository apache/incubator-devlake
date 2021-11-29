package code

import (
	"github.com/merico-dev/lake/models/domainlayer/base"
	"time"
)

type Note struct {
	base.DomainEntity
	PrId        uint64 `gorm:"index;comment:References the pull request for this note"`
	Type        string 
	Author      string
	Body        string
	Resolvable  bool `gorm:"comment:Is or is not a review comment"`
	System      bool `gorm:"comment:Is or is not auto-generated vs. human generated"`
	CreatedDate time.Time
}
