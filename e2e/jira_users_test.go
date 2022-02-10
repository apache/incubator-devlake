package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type JiraUser struct {
	Email string
}

func TestJiraUsers(t *testing.T) {
	var jiraUsers []JiraUser
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT email FROM jira_users;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var jiraUser JiraUser
		if err := rows.Scan(&jiraUser.Email); err != nil {
			panic(err)
		}
		jiraUsers = append(jiraUsers, jiraUser)
	}
	assert.Equal(t, true, len(jiraUsers) > 0)
}
