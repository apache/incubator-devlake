package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainIssueAssigneeHistory struct {
	IssueId string `json:"issue_id"`
}

func TestDomainIssueAssigneeHistorys(t *testing.T) {
	var domainBoardSprints []DomainIssueAssigneeHistory
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT issue_id FROM lake.issue_assignee_history where start_date < '2022-06-12 00:44:34.880';"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var domainBoardSprint DomainIssueAssigneeHistory
		if err := rows.Scan(&domainBoardSprint.IssueId); err != nil {
			panic(err)
		}
		domainBoardSprints = append(domainBoardSprints, domainBoardSprint)
	}
	assert.Equal(t, 2591, len(domainBoardSprints))
}
