package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainSprintIssue struct {
	IssueId int `json:"issue_id"`
}

func TestDomainSprintIssues(t *testing.T) {
	var domainSprintIssues []DomainSprintIssue
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT issue_id FROM domain_sprint_issues;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var domainSprintIssue DomainSprintIssue
		if err := rows.Scan(&domainSprintIssue.IssueId); err != nil {
			panic(err)
		}
		domainSprintIssues = append(domainSprintIssues, domainSprintIssue)
	}
	assert.Equal(t, len(domainSprintIssues) > 0, true)
}
