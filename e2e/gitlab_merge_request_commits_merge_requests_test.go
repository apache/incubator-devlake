package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type MergeRequestCommitsJoin struct {
	GitlabId string `json:"gitlab_id"`
}

func TestGitLabMergeRequestCommitsJoins(t *testing.T) {
	var mergeRequestCommitJoins []MergeRequestCommitsJoin
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "Select gitlab_id from gitlab_merge_requests mr JOIN gitlab_merge_request_commit_merge_requests mrj ON mrj.merge_request_id = mr.gitlab_id where gitlab_created_at < '2019-04-25 04:40:11.000';"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var mergeRequestCommitJoin MergeRequestCommitsJoin
		if err := rows.Scan(&mergeRequestCommitJoin.GitlabId); err != nil {
			panic(err)
		}
		mergeRequestCommitJoins = append(mergeRequestCommitJoins, mergeRequestCommitJoin)
	}
	assert.Equal(t, 1957, len(mergeRequestCommitJoins))
}
