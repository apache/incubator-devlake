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
	sqlCommand := "SELECT issue_id FROM jira_issues;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
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
	assert.Equal(t, len(jiraIssues) > 0, true)
}
