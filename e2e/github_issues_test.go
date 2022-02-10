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
	rows, err := db.Query("SELECT github_id FROM github_issues where github_created_at < '2021-12-25 04:40:11.000'")
	if err != nil {
		fmt.Println("ERROR: ", err)
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
	assert.Equal(t, 490, len(issues))
}
