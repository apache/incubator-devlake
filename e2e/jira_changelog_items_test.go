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

type JiraChangelogItem struct {
	ChangelogId int `json:"changelog_id"`
}

func TestJiraChangelogItems(t *testing.T) {
	var jiraChangelogItems []JiraChangelogItem
	db, err := InitializeDb()
	assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}
	sqlCommand := "SELECT jc.changelog_id FROM lake.jira_changelogs jc JOIN jira_changelog_items jci ON jci.changelog_id = jc.changelog_id where created < '2020-07-05 00:17:32.778';"
	rows, err := db.Query(sqlCommand)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	assert.Nil(t, err)
	defer rows.Close()
	for rows.Next() {
		var jiraChangelogItem JiraChangelogItem
		if err := rows.Scan(&jiraChangelogItem.ChangelogId); err != nil {
			panic(err)
		}
		jiraChangelogItems = append(jiraChangelogItems, jiraChangelogItem)
	}
	assert.Equal(t, 4293, len(jiraChangelogItems))
}
