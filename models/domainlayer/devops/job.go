package devops

import (
	"github.com/merico-dev/lake/models/domainlayer/base"
)

type Job struct {
	base.DomainEntity
	Name string
}
