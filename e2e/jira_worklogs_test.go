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

type JiraWorklog struct {
	TimeSpent string `json:"time_spent"`
}

func TestJiraWorklogs(t *testing.T) {
	var jiraWorklogs []JiraWorklog
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT time_spent FROM lake.jira_issues ji JOIN jira_worklogs jw ON ji.issue_id = jw.issue_id where resolution_date < '2020-06-19 06:31:18.495';"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var jiraWorklog JiraWorklog
		if err := rows.Scan(&jiraWorklog.TimeSpent); err != nil {
			panic(err)
		}
		jiraWorklogs = append(jiraWorklogs, jiraWorklog)
	}
	assert.Equal(t, 41, len(jiraWorklogs))
}
