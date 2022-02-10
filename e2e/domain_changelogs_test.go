package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainChangelog struct {
	Id int
}

func TestDomainChangelogs(t *testing.T) {
	var domainChangelogs []DomainChangelog
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT id FROM changelogs;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var domainChangelog DomainChangelog
		if err := rows.Scan(&domainChangelog.Id); err != nil {
			panic(err)
		}
		domainChangelogs = append(domainChangelogs, domainChangelog)
	}
	assert.Equal(t, len(domainChangelogs), 1)
}
