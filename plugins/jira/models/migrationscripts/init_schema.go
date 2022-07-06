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

package migrationscripts

import (
	"context"

	"github.com/apache/incubator-devlake/plugins/jira/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type InitSchemas struct{}

func (*InitSchemas) Up(ctx context.Context, db *gorm.DB) error {
	return db.Migrator().AutoMigrate(
		&archived.JiraAccount{},
		&archived.JiraBoardIssue{},
		&archived.JiraBoard{},
		&archived.JiraChangelogItem{},
		&archived.JiraChangelog{},
		&archived.JiraConnection{},
		&archived.JiraIssueCommit{},
		&archived.JiraIssueLabel{},
		&archived.JiraIssue{},
		&archived.JiraProject{},
		&archived.JiraRemotelink{},
		&archived.JiraSprint{},
		&archived.JiraBoardSprint{},
		&archived.JiraSprintIssue{},
		&archived.JiraStatus{},
		&archived.JiraWorklog{},
	)
}

func (*InitSchemas) Version() uint64 {
	return 20220407201138
}

func (*InitSchemas) Name() string {
	return "Jira init schemas"
}
