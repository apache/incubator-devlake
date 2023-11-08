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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*addConnIdToLabels)(nil)

type buglabel20231026 struct {
	ConnectionId uint64 `gorm:"primaryKey;autoIncrement:false"`
	BugId        uint64 `gorm:"primaryKey;autoIncrement:false"`
	LabelName    string `gorm:"primaryKey;type:varchar(255)"`
	archived.NoPKModel
}

type storylabel20231026 struct {
	ConnectionId uint64 `gorm:"primaryKey;autoIncrement:false"`
	StoryId      uint64 `gorm:"primaryKey;autoIncrement:false"`
	LabelName    string `gorm:"primaryKey;type:varchar(255)"`
	archived.NoPKModel
}

type tasklabel20231026 struct {
	ConnectionId uint64 `gorm:"primaryKey;autoIncrement:false"`
	TaskId       uint64 `gorm:"primaryKey;autoIncrement:false"`
	LabelName    string `gorm:"primaryKey;type:varchar(255)"`
	archived.NoPKModel
}

func (buglabel20231026) TableName() string {
	return "_tool_tapd_bug_labels"
}

func (storylabel20231026) TableName() string {
	return "_tool_tapd_story_labels"
}

func (tasklabel20231026) TableName() string {
	return "_tool_tapd_task_labels"
}

type addConnIdToLabels struct{}

func (script *addConnIdToLabels) Up(basicRes context.BasicRes) errors.Error {
	tables := []interface{}{&buglabel20231026{}, &storylabel20231026{}, &tasklabel20231026{}}
	db := basicRes.GetDal()
	err := db.DropTables(tables...)
	if err != nil {
		return err
	}
	return migrationhelper.AutoMigrateTables(basicRes, tables...)
}

func (*addConnIdToLabels) Version() uint64 {
	return 20231026000002
}

func (script *addConnIdToLabels) Name() string {
	return "add ConnectionId to Labels tables"
}
