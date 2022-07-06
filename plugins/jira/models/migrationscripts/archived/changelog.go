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

package archived

import (
	"time"

	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
)

type JiraChangelog struct {
	archived.NoPKModel

	// collected fields
	ConnectionId      uint64 `gorm:"primaryKey"`
	ChangelogId       uint64 `gorm:"primarykey"`
	IssueId           uint64 `gorm:"index"`
	AuthorAccountId   string `gorm:"type:varchar(255)"`
	AuthorDisplayName string `gorm:"type:varchar(255)"`
	AuthorActive      bool
	Created           time.Time  `gorm:"index"`
	IssueUpdated      *time.Time `comment:"corresponding issue.updated time, changelog might need update IFF changelog.issue_updated < issue.updated"`
}

type JiraChangelogItem struct {
	archived.NoPKModel

	// collected fields
	ConnectionId uint64 `gorm:"primaryKey"`
	ChangelogId  uint64 `gorm:"primaryKey"`
	Field        string `gorm:"primaryKey"`
	FieldType    string
	FieldId      string
	FromValue    string
	FromString   string
	ToValue      string
	ToString     string
}

func (JiraChangelog) TableName() string {
	return "_tool_jira_changelogs"
}

func (JiraChangelogItem) TableName() string {
	return "_tool_jira_changelog_items"
}
