package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type Repo struct {
	ID string
}

func TestGitHubRepos(t *testing.T) {
	var repos []Repo
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT * FROM github_repos")
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var repo Repo
		// This fails because we need all the column names
		// It only passes when there are no repos in the DB.
		if err := rows.Scan(&repo.ID); err != nil {
			panic(err)
		}
		repos = append(repos, repo)
	}
	assert.Equal(t, len(repos) == 0, true)
}
