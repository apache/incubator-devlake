package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainPullRequest struct {
	Id string
}

func TestDomainPullRequests(t *testing.T) {
	var domainPullRequests []DomainPullRequest
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT id FROM pull_requests;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var domainPullRequest DomainPullRequest
		if err := rows.Scan(&domainPullRequest.Id); err != nil {
			panic(err)
		}
		domainPullRequests = append(domainPullRequests, domainPullRequest)
	}
	assert.Equal(t, len(domainPullRequests), 6)
}
