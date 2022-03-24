package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type JiraRemoteLink struct {
	RemotelinkId string `json:"remotelink_id"`
}

func TestJiraRemoteLinks(t *testing.T) {
	var jiraRemoteLinks []JiraRemoteLink
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT remotelink_id FROM lake.jira_remotelinks rl JOIN jira_issues ji ON ji.issue_id = rl.issue_id where resolution_date < '2020-06-19 06:31:18.495';"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var jiraRemoteLink JiraRemoteLink
		if err := rows.Scan(&jiraRemoteLink.RemotelinkId); err != nil {
			panic(err)
		}
		jiraRemoteLinks = append(jiraRemoteLinks, jiraRemoteLink)
	}
	assert.Equal(t, 43, len(jiraRemoteLinks))
}
