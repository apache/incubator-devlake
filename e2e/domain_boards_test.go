package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainBoard struct {
	Id int
}

func TestDomainBoards(t *testing.T) {
	var domainBoards []DomainBoard
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT id FROM boards;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var domainBoard DomainBoard
		if err := rows.Scan(&domainBoard.Id); err != nil {
			panic(err)
		}
		domainBoards = append(domainBoards, domainBoard)
	}
	assert.Equal(t, len(domainBoards), 1)
}
