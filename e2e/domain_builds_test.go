package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainBuild struct {
	Id string
}

func TestDomainBuilds(t *testing.T) {
	var domainBuilds []DomainBuild
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT id FROM lake.builds where started_date < '2021-05-14 15:06:50.000';;"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var domainBuild DomainBuild
		if err := rows.Scan(&domainBuild.Id); err != nil {
			panic(err)
		}
		domainBuilds = append(domainBuilds, domainBuild)
	}
	assert.Equal(t, 4266, len(domainBuilds))
}
