package factory

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/devops"
)

func CreateBuild(jobId string) (*devops.Build, error) {
	build := &devops.Build{
		DomainEntity: domainlayer.DomainEntity{
			Id: RandIntString(),
		},
		JobName:     jobId, // ref to job
		Name:        "",
		CommitSha:   "",
		DurationSec: uint64(RandInt()),
		Status:      "",
		StartedDate: time.Now(),
	}
	return build, nil
}
