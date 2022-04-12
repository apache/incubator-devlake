package devops

import (
	"github.com/merico-dev/lake/models/domainlayer"
)

type Job struct {
	Name string `gorm:"type:char(255)"`
	domainlayer.DomainEntity
}
