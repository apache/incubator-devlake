package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type JiraBoardSprint struct {
	BoardId int `json:"board_id"`
}

func TestJiraBoardSprints(t *testing.T) {
	var jiraBoardSprints []JiraBoardSprint
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT board_id FROM jira_board_sprints;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var jiraBoardSprint JiraBoardSprint
		if err := rows.Scan(&jiraBoardSprint.BoardId); err != nil {
			panic(err)
		}
		jiraBoardSprints = append(jiraBoardSprints, jiraBoardSprint)
	}
	assert.Equal(t, len(jiraBoardSprints), 66)
}
