package factory

import (
	"fmt"
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
)

func CreatePr(repoId uint64) (*code.PullRequest, error) {
	pr := &code.PullRequest{
		DomainEntity: domainlayer.DomainEntity{
			Id: RandIntString(),
		},
		RepoId:      fmt.Sprintf("%v", repoId),
		Status:      "",
		Title:       "",
		Url:         "",
		CreatedDate: time.Now(),
		MergedDate:  nil,
		ClosedAt:    nil,
	}
	return pr, nil
}
