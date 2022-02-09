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
	rows, err := db.Query("SELECT github_id FROM github_pull_requests")
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
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
	assert.Equal(t, len(pullRequests) == 0, true)
}
