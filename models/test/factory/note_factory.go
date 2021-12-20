package factory

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/code"
)

func CreateNote(prId uint64) (*code.Note, error) {
	note := &code.Note{
		DomainEntity: domainlayer.DomainEntity{
			Id: RandIntString(),
		},
		PrId:        prId,
		Author:      "note.AuthorUsername",
		Body:        "note.Body",
		CreatedDate: time.Now(),
	}

	return note, nil
}
