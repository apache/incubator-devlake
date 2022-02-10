package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type JiraSprintIssue struct {
	IssueId int `json:"issue_id"`
}

func TestJiraSprintIssues(t *testing.T) {
	var jiraSprintIssues []JiraSprintIssue
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT issue_id FROM jira_sprint_issues;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var jiraSprintIssue JiraSprintIssue
		if err := rows.Scan(&jiraSprintIssue.IssueId); err != nil {
			panic(err)
		}
		jiraSprintIssues = append(jiraSprintIssues, jiraSprintIssue)
	}
	assert.Equal(t, len(jiraSprintIssues) > 0, true)
}
