package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type GithubIssueLabels struct {
	GithubId int `json:"github_id"`
}

func TestGithubIssueLabelss(t *testing.T) {
	var issues []GithubIssueLabels
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT github_id FROM lake.github_issue_labels")
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var issue GithubIssueLabels
		if err := rows.Scan(&issue.GithubId); err != nil {
			panic(err)
		}
		issues = append(issues, issue)
	}
	assert.Equal(t, true, len(issues) > 0)
}
