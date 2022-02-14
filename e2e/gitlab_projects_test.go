package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type Project struct {
	GitlabId int `json:"gitlab_id"`
}

func TestGitLabProjects(t *testing.T) {
	var projects []Project
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT gitlab_id FROM gitlab_projects"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var project Project
		if err := rows.Scan(&project.GitlabId); err != nil {
			panic(err)
		}
		projects = append(projects, project)
	}
	assert.Equal(t, 1, len(projects))
}
