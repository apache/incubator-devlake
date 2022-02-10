package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type JiraBoardIssue struct {
	BoardIssueId int `json:"board_id"`
}

func TestJiraBoardIssues(t *testing.T) {
	var jiraBoardIssues []JiraBoardIssue
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT board_id FROM jira_board_issues;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var jiraBoardIssue JiraBoardIssue
		if err := rows.Scan(&jiraBoardIssue.BoardIssueId); err != nil {
			panic(err)
		}
		jiraBoardIssues = append(jiraBoardIssues, jiraBoardIssue)
	}
	assert.Equal(t, len(jiraBoardIssues), 5923)
}
