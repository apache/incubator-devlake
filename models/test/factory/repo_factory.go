package factory

import (
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
)

func CreateRepo() (*code.Repo, error) {
	repo := &code.Repo{
		DomainEntity: domainlayer.DomainEntity{
			OriginKey: "something",
		},
	}
	return repo, nil
}
