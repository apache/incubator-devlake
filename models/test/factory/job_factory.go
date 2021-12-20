package factory

import (
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/devops"
)

func CreateJob() (*devops.Job, error) {
	job := &devops.Job{
		DomainEntity: domainlayer.DomainEntity{
			Id: RandIntString(),
		},
		Name: "",
	}
	return job, nil
}
