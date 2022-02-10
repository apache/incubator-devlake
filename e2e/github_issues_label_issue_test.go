package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type GithubIssueLabelIssue struct {
	GithubId int `json:"github_id"`
}

func TestGithubIssueLabelIssues(t *testing.T) {
	var issues []GithubIssueLabelIssue
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT github_id FROM lake.github_issues gi JOIN lake.github_issue_label_issues gili ON gili.issue_id = gi.github_id where github_created_at < '2021-11-24 19:22:29.000';")
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var issue GithubIssueLabelIssue
		if err := rows.Scan(&issue.GithubId); err != nil {
			panic(err)
		}
		issues = append(issues, issue)
	}
	assert.Equal(t, 561, len(issues))
}
