package ticket

import (
	"github.com/merico-dev/lake/models/domainlayer"
)

type Board struct {
	domainlayer.DomainEntity
	Name string
	Url  string
}
