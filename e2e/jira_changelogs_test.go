package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type JiraChangelog struct {
	ChangelogId int `json:"changelog_id"`
}

func TestJiraChangelogs(t *testing.T) {
	var jiraChangelogs []JiraChangelog
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT changelog_id FROM jira_changelogs where created < '2020-07-05 00:17:32.778';"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var jiraChangelog JiraChangelog
		if err := rows.Scan(&jiraChangelog.ChangelogId); err != nil {
			panic(err)
		}
		jiraChangelogs = append(jiraChangelogs, jiraChangelog)
	}
	assert.Equal(t, 3494, len(jiraChangelogs))
}
