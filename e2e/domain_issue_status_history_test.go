package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainIssueSprintHistory struct {
	SprintId string `json:"sprint_id"`
}

func TestDomainIssueSprintHistorys(t *testing.T) {
	var domainBoardSprints []DomainIssueSprintHistory
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT sprint_id FROM lake.issue_sprints_history where start_date < '2020-06-17 07:29:24.052';;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var domainBoardSprint DomainIssueSprintHistory
		if err := rows.Scan(&domainBoardSprint.SprintId); err != nil {
			panic(err)
		}
		domainBoardSprints = append(domainBoardSprints, domainBoardSprint)
	}
	assert.Equal(t, 8, len(domainBoardSprints))
}
