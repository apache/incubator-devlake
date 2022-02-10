package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type JiraIssue struct {
	IssueId int `json:"issue_id"`
}

func TestJiraIssues(t *testing.T) {
	var jiraIssues []JiraIssue
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT issue_id FROM lake.jira_issues ji where resolution_date < '2020-06-23 10:21:23.562';"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var jiraIssue JiraIssue
		if err := rows.Scan(&jiraIssue.IssueId); err != nil {
			panic(err)
		}
		jiraIssues = append(jiraIssues, jiraIssue)
	}
	assert.Equal(t, 130, len(jiraIssues))
}
