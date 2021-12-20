package factory

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
)

func CreateChangelog(issueId string) (*ticket.Changelog, error) {
	changelog := &ticket.Changelog{
		DomainEntity: domainlayer.DomainEntity{
			Id: RandIntString(),
		},
		IssueId:     issueId, // ref to issue
		AuthorId:    "",
		AuthorName:  "",
		FieldId:     "",
		FieldName:   "",
		From:        "",
		To:          "",
		CreatedDate: time.Now(),
	}
	return changelog, nil
}
