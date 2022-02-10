package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainJob struct {
	Id string
}

func TestDomainJobs(t *testing.T) {
	var domainJobs []DomainJob
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT id FROM jobs;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var domainJob DomainJob
		if err := rows.Scan(&domainJob.Id); err != nil {
			panic(err)
		}
		domainJobs = append(domainJobs, domainJob)
	}
	assert.Equal(t, true, len(domainJobs) > 0)
}
