package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type GithubIssueEvent struct {
	GithubId int `json:"github_id"`
}

func TestGitHubIssueEvents(t *testing.T) {
	var issueEvents []GithubIssueEvent
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT github_id FROM github_issue_events where github_created_at < '2021-12-25 04:40:11.000'")
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var issueEvent GithubIssueEvent
		if err := rows.Scan(&issueEvent.GithubId); err != nil {
			panic(err)
		}
		issueEvents = append(issueEvents, issueEvent)
	}
	assert.Equal(t, 3899, len(issueEvents))
}
