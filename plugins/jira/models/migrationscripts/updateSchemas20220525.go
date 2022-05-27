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

type JiraProject20220525 struct {
	archived.NoPKModel

	// collected fields
	ConnectionId uint64 `gorm:"primarykey"`
	Id           string `gorm:"primaryKey;type:varchar(255)"`
	ProjectKey   string `gorm:"type:varchar(255)"`
	Name         string `gorm:"type:varchar(255)"`
}

func (JiraProject20220525) TableName() string {
	return "_tool_jira_projects"
}

type UpdateSchemas20220525 struct{}

func (*UpdateSchemas20220525) Up(ctx context.Context, db *gorm.DB) error {

	err := db.Migrator().RenameColumn(jiraArchived.JiraProject{}, "key", "project_key")
	if err != nil {
		return err
	}

	return nil
}

func (*UpdateSchemas20220525) Version() uint64 {
	return 20220525154646
}

func (*UpdateSchemas20220525) Name() string {
	return "update `key` columns to `project_key` at _tool_jira_projects"
}
