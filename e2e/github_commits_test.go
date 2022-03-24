package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type GithubCommit struct {
	Sha string `json:"sha"`
}

func TestGitHubCommits(t *testing.T) {
	var commits []GithubCommit
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT sha FROM github_commits where authored_date < '2021-11-24 19:22:29.000'")
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var commit GithubCommit
		if err := rows.Scan(&commit.Sha); err != nil {
			panic(err)
		}
		commits = append(commits, commit)
	}
	assert.Equal(t, 874, len(commits))
}
