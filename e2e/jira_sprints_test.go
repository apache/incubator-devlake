package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type JiraSprint struct {
	SprintId int `json:"sprint_id"`
}

func TestJiraSprints(t *testing.T) {
	var jiraSprints []JiraSprint
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT sprint_id FROM jira_sprints;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var jiraSprint JiraSprint
		if err := rows.Scan(&jiraSprint.SprintId); err != nil {
			panic(err)
		}
		jiraSprints = append(jiraSprints, jiraSprint)
	}
	assert.Equal(t, len(jiraSprints) > 0, true)
}
