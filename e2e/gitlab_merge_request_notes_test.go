package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type MergeRequestNote struct {
	GitlabId int `json:"gitlab_id"`
}

func TestGitLabMergeRequestNotes(t *testing.T) {
	var mergeRequestNotes []MergeRequestNote
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT gitlab_id FROM gitlab_merge_request_notes"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var mergeRequestNote MergeRequestNote
		if err := rows.Scan(&mergeRequestNote.GitlabId); err != nil {
			panic(err)
		}
		mergeRequestNotes = append(mergeRequestNotes, mergeRequestNote)
	}
	assert.Equal(t, len(mergeRequestNotes) > 0, true)
}
