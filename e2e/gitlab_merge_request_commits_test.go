package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type MergeRequestCommit struct {
	CommitId string `json:"commit_id"`
}

func TestGitLabMergeRequestCommits(t *testing.T) {
	var mergeRequestCommits []MergeRequestCommit
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT commit_id FROM _tool_gitlab_merge_request_commits where authored_date < '2019-06-25 02:41:42.000'"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var mergeRequestCommit MergeRequestCommit
		if err := rows.Scan(&mergeRequestCommit.CommitId); err != nil {
			panic(err)
		}
		mergeRequestCommits = append(mergeRequestCommits, mergeRequestCommit)
	}
	assert.Equal(t, 2496, len(mergeRequestCommits))
}
