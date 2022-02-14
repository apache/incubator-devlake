package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainChangelog struct {
	Id string
}

func TestDomainChangelogs(t *testing.T) {
	var domainChangelogs []DomainChangelog
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT id FROM lake.changelogs where created_date < '2020-06-20 06:18:24.880';"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
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
	assert.Equal(t, 1742, len(domainChangelogs))
}
