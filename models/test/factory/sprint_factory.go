package factory

import (
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
)

func CreateSprint(boardId string) (*ticket.Sprint, error) {
	sprint := &ticket.Sprint{
		DomainEntity: domainlayer.DomainEntity{
			Id: RandIntString(),
		},
		Url:           "",
		Status:        "",
		Name:          "",
		StartedDate:   nil,
		EndedDate:     nil,
		CompletedDate: nil,
	}
	return sprint, nil
}
