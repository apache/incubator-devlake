package factory

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/devops"
)

func CreateBuild(jobId string) (*devops.Build, error) {
	build := &devops.Build{
		DomainEntity: domainlayer.DomainEntity{
			Id: "1",
		},
		JobId:       jobId, // ref to job
		Name:        "",
		CommitSha:   "",
		DurationSec: 1,
		Status:      "",
		StartedDate: time.Now(),
	}
	return build, nil
}
