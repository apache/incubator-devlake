package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type Mr struct {
	iid int
}

func TestGitLabMrs(t *testing.T) {
	var mrs []Mr
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT iid FROM gitlab_merge_requests"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var mr Mr
		if err := rows.Scan(&mr.iid); err != nil {
			panic(err)
		}
		mrs = append(mrs, mr)
	}
	assert.Equal(t, len(mrs), 0)
}
