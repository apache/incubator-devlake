package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type GithubPullRequestCommit struct {
	CommitSha int `json:"commit_sha"`
}

func TestGithubPullRequestCommits(t *testing.T) {
	var PullRequestCommits []GithubPullRequestCommit
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT github_id FROM github_pull_request_commits")
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var PullRequestCommit GithubPullRequestCommit
		if err := rows.Scan(&PullRequestCommit.CommitSha); err != nil {
			panic(err)
		}
		PullRequestCommits = append(PullRequestCommits, PullRequestCommit)
	}
	assert.Equal(t, len(PullRequestCommits) == 0, true)
}
