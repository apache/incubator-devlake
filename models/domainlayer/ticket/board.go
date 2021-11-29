package ticket

import (
	"github.com/merico-dev/lake/models/domainlayer/base"
)

type Board struct {
	base.DomainEntity
	Name string
	Url  string
}
