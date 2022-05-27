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

	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	jiraArchived "github.com/apache/incubator-devlake/plugins/jira/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type JiraChangelogItem20220527 struct {
	archived.NoPKModel

	// collected fields
	SourceId    uint64 `gorm:"primaryKey"`
	ChangelogId uint64 `gorm:"primaryKey"`
	Field       string `gorm:"primaryKey"`
	FieldType   string
	FieldId     string
	FromValue   string
	FromString  string
	ToValue     string
	ToString    string
}

func (JiraChangelogItem20220527) TableName() string {
	return "_tool_jira_changelog_items"
}

type UpdateSchemas20220527 struct{}

func (*UpdateSchemas20220527) Up(ctx context.Context, db *gorm.DB) error {
	
	err := db.Migrator().RenameColumn(jiraArchived.JiraChangelogItem{}, "from", "from_value")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(jiraArchived.JiraChangelogItem{}, "to", "to_value")
	if err != nil {
		return err
	}
	return nil
}

func (*UpdateSchemas20220527) Version() uint64 {
	return 20220527154646
}

func (*UpdateSchemas20220527) Name() string {
	return "update `from` and `to` columns to `from_value` and `to_value` at _tool_jira_changelog_items"
}
