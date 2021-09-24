package code

import (
	"time"

	"github.com/merico-dev/lake/plugins/domainlayer/models/base"
)

type Note struct {
	base.DomainEntity
	PrId        uint64 `gorm:"index"`
	Type        string
	Author      string
	Body        string
	Resolvable  bool // Resolvable means a comment is a code review comment
	System      bool // System means the comment is generated automatically
	CreatedDate time.Time
}
