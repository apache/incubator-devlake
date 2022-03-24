package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type Reviewer struct {
	GitlabId int `json:"gitlab_id"`
}

func TestGitLabReviewers(t *testing.T) {
	var reviewers []Reviewer
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT gitlab_id FROM gitlab_reviewers"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var reviewer Reviewer
		if err := rows.Scan(&reviewer.GitlabId); err != nil {
			panic(err)
		}
		reviewers = append(reviewers, reviewer)
	}
	assert.Equal(t, true, len(reviewers) > 0)
}
