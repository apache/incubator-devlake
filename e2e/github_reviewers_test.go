package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type GithubReviewer struct {
	GithubId int `json:"github_id"`
}

func TestGithubReviewers(t *testing.T) {
	var Reviewers []GithubReviewer
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT github_id FROM github_reviewers")
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var Reviewer GithubReviewer
		if err := rows.Scan(&Reviewer.GithubId); err != nil {
			panic(err)
		}
		Reviewers = append(Reviewers, Reviewer)
	}
	assert.Equal(t, len(Reviewers) == 0, true)
}
