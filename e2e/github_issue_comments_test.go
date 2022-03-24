package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type GithubIssueComment struct {
	GithubId int `json:"github_id"`
}

func TestGitHubIssueComments(t *testing.T) {
	var issueComments []GithubIssueComment
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT github_id FROM github_issue_comments where github_created_at < '2021-12-25 04:40:11.000'")
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var issueComment GithubIssueComment
		if err := rows.Scan(&issueComment.GithubId); err != nil {
			panic(err)
		}
		issueComments = append(issueComments, issueComment)
	}
	assert.Equal(t, 781, len(issueComments))
}
