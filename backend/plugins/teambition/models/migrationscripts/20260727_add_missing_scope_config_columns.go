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
	"fmt"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

const teambitionScopeConfigTable20260727 = "_tool_teambition_scope_configs"

// teambitionScopeConfig20260727 mirrors models.TeambitionScopeConfig. The
// migration that created `_tool_teambition_scope_configs` did not include the
// columns of the embedded common.Model (`id`, `created_at`, `updated_at`),
// which the runtime model expects.
type teambitionScopeConfig20260727 struct {
	archived.Model
	Entities          []string          `gorm:"type:json;serializer:json" json:"entities"`
	ConnectionId      uint64            `json:"connectionId" gorm:"index"`
	Name              string            `json:"name" gorm:"type:varchar(255);uniqueIndex"`
	TypeMappings      map[string]string `json:"typeMappings" gorm:"serializer:json"`
	StatusMappings    map[string]string `json:"statusMappings" gorm:"serializer:json"`
	BugDueDateField   string            `json:"bugDueDateField" gorm:"column:bug_due_date_field"`
	TaskDueDateField  string            `json:"taskDueDateField" gorm:"column:task_due_date_field"`
	StoryDueDateField string            `json:"storyDueDateField" gorm:"column:story_due_date_field"`
}

func (teambitionScopeConfig20260727) TableName() string {
	return teambitionScopeConfigTable20260727
}

type addMissingScopeConfigColumns struct{}

// Up adds the columns of the embedded common.Model that the runtime model
// expects.
//
// `id` is an auto-increment primary key, which GORM's AutoMigrate cannot append
// to an existing table: it emits a plain `ADD COLUMN ... AUTO_INCREMENT`, which
// MySQL rejects with "Incorrect table definition; there can be only one auto
// column and it must be defined as a key". The column is therefore added with
// explicit DDL (the table has no primary key so far), letting the database
// backfill ids for existing rows and keep the sequence/counter in sync. The
// remaining columns (`created_at`, `updated_at`) and the indexes are then
// created by AutoMigrate as usual.
func (script *addMissingScopeConfigColumns) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if !db.HasColumn(teambitionScopeConfigTable20260727, "id") {
		ddl := fmt.Sprintf(
			"ALTER TABLE %s ADD COLUMN id BIGSERIAL PRIMARY KEY",
			teambitionScopeConfigTable20260727,
		)
		if db.Dialect() == "mysql" {
			ddl = fmt.Sprintf(
				"ALTER TABLE %s ADD COLUMN id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY",
				teambitionScopeConfigTable20260727,
			)
		}
		if err := db.Exec(ddl); err != nil {
			return err
		}
	}
	return migrationhelper.AutoMigrateTables(basicRes, &teambitionScopeConfig20260727{})
}

func (*addMissingScopeConfigColumns) Version() uint64 {
	return 20260727000001
}

func (*addMissingScopeConfigColumns) Name() string {
	return "add missing id/created_at/updated_at columns to _tool_teambition_scope_configs"
}
