package devops

import (
	"github.com/merico-dev/lake/models/domainlayer"
)

type Job struct {
	domainlayer.DomainEntity
	Name string
}
