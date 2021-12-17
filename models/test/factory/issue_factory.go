package factory

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
)

func CreateIssue(boardId string) (*ticket.Issue, error) {
	issue := &ticket.Issue{
		DomainEntity: domainlayer.DomainEntity{
			Id: "1",
		},
		BoardId:                  boardId, // ref to board
		Url:                      "",
		Key:                      "",
		Title:                    "",
		Summary:                  "",
		EpicKey:                  "",
		Type:                     "",
		Status:                   "",
		StoryPoint:               1,
		OriginalEstimateMinutes:  1, // user input?
		AggregateEstimateMinutes: 1, // sum up of all subtasks?
		RemainingEstimateMinutes: 1, // could it be negative value?
		CreatorId:                "",
		AssigneeId:               "",
		ResolutionDate:           nil,
		Priority:                 "", // not sure how to deal with it yet, copy the name for now
		ParentId:                 "",
		SprintId:                 "",
		CreatedDate:              time.Now(),
		UpdatedDate:              time.Now(),
		SpentMinutes:             1,
		LeadTimeMinutes:          1,
	}
	return issue, nil
}
