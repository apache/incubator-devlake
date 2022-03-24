package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type GithubPullRequestComment struct {
	GithubId int `json:"github_id"`
}

func TestGithubPullRequestComments(t *testing.T) {
	var PullRequestComments []GithubPullRequestComment
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT github_id FROM github_pull_request_comments where github_created_at < '2021-12-25 04:40:11.000'")
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var PullRequestComment GithubPullRequestComment
		if err := rows.Scan(&PullRequestComment.GithubId); err != nil {
			panic(err)
		}
		PullRequestComments = append(PullRequestComments, PullRequestComment)
	}
	assert.Equal(t, 401, len(PullRequestComments))
}
