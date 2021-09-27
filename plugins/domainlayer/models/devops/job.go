package devops

import (
	"github.com/merico-dev/lake/plugins/domainlayer/models/base"
)

type Job struct {
	base.DomainEntity
	Name string
}
