package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type JiraChangelogItem struct {
	ChangelogId int `json:"changelog_id"`
}

func TestJiraChangelogItems(t *testing.T) {
	var jiraChangelogItems []JiraChangelogItem
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT changelog_id FROM jira_changelog_items;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var jiraChangelogItem JiraChangelogItem
		if err := rows.Scan(&jiraChangelogItem.ChangelogId); err != nil {
			panic(err)
		}
		jiraChangelogItems = append(jiraChangelogItems, jiraChangelogItem)
	}
	assert.Equal(t, len(jiraChangelogItems), 218005)
}
