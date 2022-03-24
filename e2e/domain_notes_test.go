package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainNote struct {
	Id string
}

func TestDomainNotes(t *testing.T) {
	var domainNotes []DomainNote
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT id FROM lake.notes where created_date < '2019-10-24 02:23:43.000';"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var domainNote DomainNote
		if err := rows.Scan(&domainNote.Id); err != nil {
			panic(err)
		}
		domainNotes = append(domainNotes, domainNote)
	}
	assert.Equal(t, 4352, len(domainNotes))
}
