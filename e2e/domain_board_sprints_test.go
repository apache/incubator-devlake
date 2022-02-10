package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainBoardSprint struct {
	BoardId int `json:"board_id"`
}

func TestDomainBoardSprints(t *testing.T) {
	var domainBoardSprints []DomainBoardSprint
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT board_id FROM board_sprints;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
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
	assert.Equal(t, len(domainBoardSprints), 5923)
}
