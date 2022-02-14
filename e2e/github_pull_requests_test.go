package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type GithubPullRequest struct {
	GithubId int `json:"github_id"`
}

func TestGithubPullRequests(t *testing.T) {
	var pullRequests []GithubPullRequest
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT github_id FROM github_pull_requests where github_created_at < '2021-12-25 04:40:11.000'")
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var pullRequest GithubPullRequest
		if err := rows.Scan(&pullRequest.GithubId); err != nil {
			panic(err)
		}
		pullRequests = append(pullRequests, pullRequest)
	}
	assert.Equal(t, 512, len(pullRequests))
}
