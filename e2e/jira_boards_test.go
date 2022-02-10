package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type JiraBoard struct {
	BoardId int `json:"board_id"`
}

func TestJiraBoards(t *testing.T) {
	var jiraBoards []JiraBoard
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT board_id FROM jira_boards;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var jiraBoard JiraBoard
		if err := rows.Scan(&jiraBoard.BoardId); err != nil {
			panic(err)
		}
		jiraBoards = append(jiraBoards, jiraBoard)
	}
	assert.Equal(t, 1, len(jiraBoards))
}
