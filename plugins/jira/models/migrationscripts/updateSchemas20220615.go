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
	"time"

	"gorm.io/gorm"
)

type UpdateSchemas20220615 struct {
}

type JiraIssue20220615 struct{}

func (JiraIssue20220615) TableName() string {
	return "_tool_jira_issues"
}

type JiraChangelog20220615 struct {
	IssueUpdated *time.Time
}

func (JiraChangelog20220615) TableName() string {
	return "_tool_jira_changelogs"
}

func (*UpdateSchemas20220615) Up(ctx context.Context, db *gorm.DB) error {
	var err error
	err = db.Migrator().DropColumn(&JiraIssue20220615{}, "changelog_updated")
	if err != nil {
		return err
	}
	err = db.Migrator().AutoMigrate(&JiraChangelog20220615{})
	if err != nil {
		return err
	}

	return nil
}

func (*UpdateSchemas20220615) Version() uint64 {
	return 20220601154646
}

func (*UpdateSchemas20220615) Name() string {
	return "replace issues.changelog_updated with changelogs.issue_updated"
}
