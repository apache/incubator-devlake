package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type Pipeline struct {
	GitlabId string `json:"gitlab_id"`
}

func TestGitLabPipelines(t *testing.T) {
	var pipelines []Pipeline
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT gitlab_id FROM gitlab_pipelines where gitlab_created_at < '2019-06-25 02:41:42.000'"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var pipeline Pipeline
		if err := rows.Scan(&pipeline.GitlabId); err != nil {
			panic(err)
		}
		pipelines = append(pipelines, pipeline)
	}
	assert.Equal(t, 1333, len(pipelines))
}
