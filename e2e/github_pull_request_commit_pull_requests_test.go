package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type GithubPRJoin struct {
	GithubId int `json:"github_id"`
}

func TestGithubPRJoins(t *testing.T) {
	var issues []GithubPRJoin
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT github_id FROM lake.github_pull_requests pr JOIN github_pull_request_commit_pull_requests prj ON prj.pull_request_id = pr.github_id where github_created_at < '2021-12-25 04:40:11.000';")
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var issue GithubPRJoin
		if err := rows.Scan(&issue.GithubId); err != nil {
			panic(err)
		}
		issues = append(issues, issue)
	}
	assert.Equal(t, 1705, len(issues))
}
