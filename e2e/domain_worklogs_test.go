package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainWorklog struct {
	TimeSpentMinutes string `json:"time_spent_minutes"`
}

func TestDomainWorklogs(t *testing.T) {
	var domainWorklogs []DomainWorklog
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT time_spent_minutes FROM lake.issues;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var domainWorklog DomainWorklog
		if err := rows.Scan(&domainWorklog.TimeSpentMinutes); err != nil {
			panic(err)
		}
		domainWorklogs = append(domainWorklogs, domainWorklog)
	}
	assert.Equal(t, len(domainWorklogs), 41)
}
