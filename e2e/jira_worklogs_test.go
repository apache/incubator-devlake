package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type JiraWorklog struct {
	Email string
}

func TestJiraWorklogs(t *testing.T) {
	var jiraWorklogs []JiraWorklog
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT email FROM jira_worklogs;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var jiraWorklog JiraWorklog
		if err := rows.Scan(&jiraWorklog.Email); err != nil {
			panic(err)
		}
		jiraWorklogs = append(jiraWorklogs, jiraWorklog)
	}
	assert.Equal(t, len(jiraWorklogs) > 0, true)
}
