package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainBoardIssue struct {
	BoardId int `json:"board_id"`
}

func TestDomainBoardIssues(t *testing.T) {
	var domainBoardIssues []DomainBoardIssue
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT board_id FROM board_issues;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
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
	assert.Equal(t, len(domainBoardIssues), 5923)
}
