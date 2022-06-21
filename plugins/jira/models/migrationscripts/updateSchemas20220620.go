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
	"gorm.io/gorm"
	"time"
)

type UpdateSchemas20220620 struct {
}

type JiraIssue20220620 struct{}

func (JiraIssue20220620) TableName() string {
	return "_tool_jira_issues"
}

type JiraWorklog20220620 struct {
	IssueUpdated *time.Time
}

func (JiraWorklog20220620) TableName() string {
	return "_tool_jira_worklogs"
}

type JiraRemotelink20220620 struct {
	IssueUpdated *time.Time
}

func (JiraRemotelink20220620) TableName() string {
	return "_tool_jira_remotelinks"
}

func (*UpdateSchemas20220620) Up(ctx context.Context, db *gorm.DB) error {
	var err error
	err = db.Migrator().DropColumn(&JiraIssue20220620{}, "worklog_updated")
	if err != nil {
		return err
	}
	err = db.Migrator().DropColumn(&JiraIssue20220620{}, "remotelink_updated")
	if err != nil {
		return err
	}
	err = db.Migrator().AutoMigrate(&JiraWorklog20220620{}, &JiraRemotelink20220620{})
	if err != nil {
		return err
	}

	return nil
}

func (*UpdateSchemas20220620) Version() uint64 {
	return 20220620101111
}

func (*UpdateSchemas20220620) Name() string {
	return "add column issue_updated to _tool_jira_worklogs and _tool_jira_remotelinks"
}
