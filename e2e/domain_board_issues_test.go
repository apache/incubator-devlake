package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainBoardIssue struct {
	BoardId string `json:"board_id"`
}

func TestDomainBoardIssues(t *testing.T) {
	var domainBoardIssues []DomainBoardIssue
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT board_id FROM lake.issues i JOIN board_issues bi on i.id = bi.issue_id where resolution_date < '2021-10-25 17:00:58.000';"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var domainBoardIssue DomainBoardIssue
		if err := rows.Scan(&domainBoardIssue.BoardId); err != nil {
			panic(err)
		}
		domainBoardIssues = append(domainBoardIssues, domainBoardIssue)
	}
	assert.Equal(t, 2687, len(domainBoardIssues))
}
