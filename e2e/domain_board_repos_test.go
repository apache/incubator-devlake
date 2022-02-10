package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainBoardRepo struct {
	BoardId int `json:"board_id"`
}

func TestDomainBoardRepos(t *testing.T) {
	var domainBoardRepos []DomainBoardRepo
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT board_id FROM board_repos;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("KEVIN >>> err", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var domainBoardRepo DomainBoardRepo
		if err := rows.Scan(&domainBoardRepo.BoardId); err != nil {
			panic(err)
		}
		domainBoardRepos = append(domainBoardRepos, domainBoardRepo)
	}
	assert.Equal(t, len(domainBoardRepos), 5923)
}
