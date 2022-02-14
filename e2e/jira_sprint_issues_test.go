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
	sqlCommand := "SELECT si.issue_id FROM lake.jira_sprint_issues si JOIN jira_issues ji ON ji.issue_id = si.issue_id where resolution_date < '2020-06-19 06:31:18.495';"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
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
	assert.Equal(t, 78, len(jiraSprintIssues))
}
