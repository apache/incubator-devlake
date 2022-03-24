package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainBoardSprint struct {
	BoardId string `json:"board_id"`
}

func TestDomainBoardSprints(t *testing.T) {
	var domainBoardSprints []DomainBoardSprint
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT bs.board_id FROM lake.board_sprints bs join sprints s ON bs.sprint_id = s.id where ended_date < '2022-01-14 14:10:00.000';"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var domainBoardSprint DomainBoardSprint
		if err := rows.Scan(&domainBoardSprint.BoardId); err != nil {
			panic(err)
		}
		domainBoardSprints = append(domainBoardSprints, domainBoardSprint)
	}
	assert.Equal(t, 61, len(domainBoardSprints))
}
