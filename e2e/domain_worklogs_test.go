/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test should only run once main_test is complete and ready

type DomainWorklog struct {
	AuthorId string `json:"author_id"`
}

func TestDomainWorklogs(t *testing.T) {
	var domainWorklogs []DomainWorklog
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT author_id FROM lake.worklogs w JOIN lake.issues i ON w.issue_id = i.id where started_date < '2020-06-20 06:18:24.880';"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var domainWorklog DomainWorklog
		if err := rows.Scan(&domainWorklog.AuthorId); err != nil {
			panic(err)
		}
		domainWorklogs = append(domainWorklogs, domainWorklog)
	}
	assert.Equal(t, 987, len(domainWorklogs))
}
