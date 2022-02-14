package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type User struct {
	Email string
}

func TestGitLabUsers(t *testing.T) {
	var users []User
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT email FROM gitlab_users"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Email); err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	assert.Equal(t, true, len(users) > 0)
}
