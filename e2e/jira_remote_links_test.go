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
	sqlCommand := "SELECT remotelink_id FROM jira_remotelinks;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
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
	assert.Equal(t, len(jiraRemoteLinks) > 0, true)
}
