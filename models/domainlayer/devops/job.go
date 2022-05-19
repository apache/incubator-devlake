package devops

import (
	"github.com/apache/incubator-devlake/models/domainlayer"
)

type Job struct {
	Name string `gorm:"type:varchar(255)"`
	domainlayer.DomainEntity
}
