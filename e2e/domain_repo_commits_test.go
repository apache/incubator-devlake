package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainPullRequestCommit struct {
	CommitSha string `json:"commit_sha"`
}

func TestDomainPullRequestCommits(t *testing.T) {
	var domainPullRequestCommits []DomainPullRequestCommit
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT commit_sha FROM pull_request_commits;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var domainPullRequestCommit DomainPullRequestCommit
		if err := rows.Scan(&domainPullRequestCommit.CommitSha); err != nil {
			panic(err)
		}
		domainPullRequestCommits = append(domainPullRequestCommits, domainPullRequestCommit)
	}
	assert.Equal(t, len(domainPullRequestCommits), 6)
}
