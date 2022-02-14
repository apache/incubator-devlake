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
	sqlCommand := "SELECT board_id FROM lake.jira_issues ji JOIN jira_board_issues jbi ON ji.issue_id = jbi.issue_id where resolution_date < '2020-10-15 08:59:51.304';"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
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
	assert.Equal(t, 894, len(jiraBoardIssues))
}
