package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type JiraProject struct {
	Id int
}

func TestJiraProjects(t *testing.T) {
	var jiraProjects []JiraProject
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT id FROM jira_projects;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var jiraProject JiraProject
		if err := rows.Scan(&jiraProject.Id); err != nil {
			panic(err)
		}
		jiraProjects = append(jiraProjects, jiraProject)
	}
	assert.Equal(t,  true, len(jiraProjects) > 0)
}
