package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type JenkinsJob struct {
	Id int
}

func TestJenkinsJobs(t *testing.T) {
	var jenkinsJobs []JenkinsJob
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT id FROM jenkins_jobs;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var jenkinsJob JenkinsJob
		if err := rows.Scan(&jenkinsJob.Id); err != nil {
			panic(err)
		}
		jenkinsJobs = append(jenkinsJobs, jenkinsJob)
	}
	assert.Equal(t, true, len(jenkinsJobs) > 0)
}
