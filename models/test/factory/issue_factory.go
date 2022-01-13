package factory

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
)

func CreateIssue() (*ticket.Issue, error) {
	now := time.Now()
	issue := &ticket.Issue{
		DomainEntity: domainlayer.DomainEntity{
			Id: RandIntString(),
		},
		Url:                      "",
		Key:                      "",
		Title:                    "",
		Summary:                  "",
		EpicKey:                  "",
		Type:                     "",
		Status:                   "",
		StoryPoint:               1,
		OriginalEstimateMinutes:  1, // user input?
		CreatorId:                "",
		AssigneeId:               "",
		ResolutionDate:           nil,
		Priority:                 "", // not sure how to deal with it yet, copy the name for now
		ParentIssueId:            "",
		CreatedDate:              &now,
		UpdatedDate:              &now,
		LeadTimeMinutes:          1,
	}
	return issue, nil
}
