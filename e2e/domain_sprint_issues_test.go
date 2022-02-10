package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainSprintIssue struct {
	IssueId string `json:"issue_id"`
}

func TestDomainSprintIssues(t *testing.T) {
	var domainSprintIssues []DomainSprintIssue
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT issue_id FROM lake.sprint_issues si JOIN lake.issues i ON si.issue_id = i.id where created_date < '2020-06-20 06:18:24.880';"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
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
	assert.Equal(t, 999, len(domainSprintIssues))
}
