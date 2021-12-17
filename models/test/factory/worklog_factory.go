package factory

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
)

func CreateWorklog(boardId string, issueId string) (*ticket.Worklog, error) {
	worklog := &ticket.Worklog{
		DomainEntity: domainlayer.DomainEntity{
			Id: "1",
		},
		IssueId:          issueId, // ref to issue
		BoardId:          boardId, // ref to board
		AuthorId:         "",
		UpdateAuthorId:   "",
		TimeSpent:        "",
		TimeSpentSeconds: 1,
		Updated:          time.Now(),
		Started:          time.Now(),
	}
	return worklog, nil
}
