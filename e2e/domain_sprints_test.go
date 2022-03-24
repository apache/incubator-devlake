package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainSprint struct {
	Id string
}

func TestDomainSprints(t *testing.T) {
	var domainSprints []DomainSprint
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT id FROM lake.sprints where started_date < '2020-07-27 01:26:13.465';"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var domainSprint DomainSprint
		if err := rows.Scan(&domainSprint.Id); err != nil {
			panic(err)
		}
		domainSprints = append(domainSprints, domainSprint)
	}
	assert.Equal(t, 4, len(domainSprints))
}
