package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainWorklog struct {
	AuthorId string `json:"author_id"`
}

func TestDomainWorklogs(t *testing.T) {
	var domainWorklogs []DomainWorklog
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT author_id FROM lake.worklogs w JOIN lake.issues i ON w.issue_id = i.id where started_date < '2020-06-20 06:18:24.880';"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var domainWorklog DomainWorklog
		if err := rows.Scan(&domainWorklog.AuthorId); err != nil {
			panic(err)
		}
		domainWorklogs = append(domainWorklogs, domainWorklog)
	}
	assert.Equal(t, 987, len(domainWorklogs))
}
