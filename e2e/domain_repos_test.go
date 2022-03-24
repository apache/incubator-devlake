package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainRepo struct {
	Id string
}

func TestDomainRepos(t *testing.T) {
	var repos []DomainRepo
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT id FROM repos")
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var repo DomainRepo
		if err := rows.Scan(&repo.Id); err != nil {
			panic(err)
		}
		repos = append(repos, repo)
	}
	assert.Equal(t, true, len(repos) > 0)
}
