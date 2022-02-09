package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type GithubIssue struct {
	GithubId int `json:"github_id"`
}

func TestGithubIssues(t *testing.T) {
	var issues []GithubIssue
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT github_id FROM github_issues")
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var issue GithubIssue
		if err := rows.Scan(&issue.GithubId); err != nil {
			panic(err)
		}
		issues = append(issues, issue)
	}
	assert.Equal(t, len(issues) == 0, true)
}
