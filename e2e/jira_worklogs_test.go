package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type JiraWorklog struct {
	TimeSpent string `json:"time_spent"`
}

func TestJiraWorklogs(t *testing.T) {
	var jiraWorklogs []JiraWorklog
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT time_spent FROM lake.jira_issues ji JOIN jira_worklogs jw ON ji.issue_id = jw.issue_id where resolution_date < '2020-06-19 06:31:18.495';"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var jiraWorklog JiraWorklog
		if err := rows.Scan(&jiraWorklog.TimeSpent); err != nil {
			panic(err)
		}
		jiraWorklogs = append(jiraWorklogs, jiraWorklog)
	}
	assert.Equal(t, 41, len(jiraWorklogs))
}
