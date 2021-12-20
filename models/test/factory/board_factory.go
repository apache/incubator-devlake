package factory

import (
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
)

func CreateBoard() (*ticket.Board, error) {
	board := &ticket.Board{
		DomainEntity: domainlayer.DomainEntity{
			Id: RandIntString(),
		},
		Name: "",
		Url:  "",
	}
	return board, nil
}
