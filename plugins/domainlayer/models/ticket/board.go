package ticket

import (
	"github.com/merico-dev/lake/plugins/domainlayer/models/base"
)

type Board struct {
	base.DomainEntity
	Name string
	Url  string
}
