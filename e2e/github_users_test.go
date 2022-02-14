package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type GithubUser struct {
	Id int
}

func TestGithubUsers(t *testing.T) {
	var users []GithubUser
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT id FROM github_users"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var user GithubUser
		if err := rows.Scan(&user.Id); err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	assert.Equal(t, true, len(users) > 0)
}
