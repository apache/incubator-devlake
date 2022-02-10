package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainIssue struct {
	Id string
}

func TestDomainIssues(t *testing.T) {
	var issues []DomainIssue
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT id FROM issues where created_at < '2021-12-25 04:40:11.000'")
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var issue DomainIssue
		if err := rows.Scan(&issue.Id); err != nil {
			panic(err)
		}
		issues = append(issues, issue)
	}
	assert.Equal(t, len(issues), 490)
}
