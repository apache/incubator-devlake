package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainBoard struct {
	Id string
}

func TestDomainBoards(t *testing.T) {
	var domainBoards []DomainBoard
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT id FROM lake.builds where started_date < '2021-04-09 09:49:07.000';"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
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
	assert.Equal(t, 2904, len(domainBoards))
}
