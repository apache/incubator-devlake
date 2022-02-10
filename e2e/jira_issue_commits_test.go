package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type JiraIssueCommit struct {
	IssueId int `json:"issue_id"`
}

func TestJiraIssueCommits(t *testing.T) {
	var jiraIssueCommits []JiraIssueCommit
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT issue_id FROM jira_board_sprints;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var jiraIssueCommit JiraIssueCommit
		if err := rows.Scan(&jiraIssueCommit.IssueId); err != nil {
			panic(err)
		}
		jiraIssueCommits = append(jiraIssueCommits, jiraIssueCommit)
	}
	assert.Equal(t, len(jiraIssueCommits) > 0, true)
}
