package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type GithubRepo struct {
	GithubId string `json:"github_id"`
}

func TestGitHubRepos(t *testing.T) {
	var repos []GithubRepo
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT github_id FROM github_repos")
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var repo GithubRepo
		if err := rows.Scan(&repo.GithubId); err != nil {
			panic(err)
		}
		repos = append(repos, repo)
	}
	assert.Equal(t, 1, len(repos))
}
