package factory

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
)

func CreateCommit(repoId uint64) (*code.Commit, error) {
	commit := &code.Commit{
		DomainEntity: domainlayer.DomainEntity{
			Id: "1",
		},
		RepoId:         repoId,
		Sha:            "dosifj9302hf80h23f",
		Additions:      1,
		Deletions:      1,
		DevEq:          1,
		Message:        "",
		AuthorName:     "",
		AuthorEmail:    "",
		AuthoredDate:   time.Now(),
		CommitterName:  "",
		CommitterEmail: "",
		CommittedDate:  time.Now(),
	}
	return commit, nil
}
