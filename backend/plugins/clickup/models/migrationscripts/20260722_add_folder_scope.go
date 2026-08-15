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
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

// addFolderScope moves the plugin's data-source scope from List to Folder
// (Jira-parity: folder = board). It creates _tool_clickup_folders and adds the
// sprint/story-point/folder columns to the existing tool tables. AutoMigrate is
// additive: pre-existing columns from the list-scope schema are left in place.

// frozen snapshots for this migration (suffixed to avoid clashing with the
// 20260720 archived snapshots that share the same table names).

type clickUpFolder20260722 struct {
	archived.NoPKModel
	ConnectionId  uint64 `gorm:"primaryKey"`
	ScopeConfigId uint64
	FolderId      string `gorm:"primaryKey;type:varchar(255)"`
	Name          string `gorm:"type:varchar(255)"`
	SpaceId       string `gorm:"type:varchar(255)"`
	SpaceName     string `gorm:"type:varchar(255)"`
}

func (clickUpFolder20260722) TableName() string { return "_tool_clickup_folders" }

type clickUpList20260722 struct {
	archived.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey"`
	ListId       string `gorm:"primaryKey;type:varchar(255)"`
	FolderId     string `gorm:"index;type:varchar(255)"`
	SpaceId      string `gorm:"type:varchar(255)"`
	Name         string `gorm:"type:varchar(255)"`
	Archived     bool
	IsSprint     bool
	SprintName   string `gorm:"type:varchar(255)"`
	StartDate    *time.Time
	EndDate      *time.Time
}

func (clickUpList20260722) TableName() string { return "_tool_clickup_lists" }

type clickUpTask20260722 struct {
	ConnectionId uint64   `gorm:"primaryKey"`
	Id           string   `gorm:"primaryKey;type:varchar(255)"`
	FolderId     string   `gorm:"index;type:varchar(255)"`
	StoryPoint   *float64 `gorm:"column:story_point"`
	archived.NoPKModel
}

func (clickUpTask20260722) TableName() string { return "_tool_clickup_tasks" }

type clickUpScopeConfig20260722 struct {
	archived.ScopeConfig
	SprintNamePattern string `gorm:"type:varchar(255)"`
	StoryPointField   string `gorm:"type:varchar(255)"`
	DefaultIssueType  string `gorm:"type:varchar(100)"`
}

func (clickUpScopeConfig20260722) TableName() string { return "_tool_clickup_scope_configs" }

type addFolderScope struct{}

func (*addFolderScope) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&clickUpFolder20260722{},
		&clickUpList20260722{},
		&clickUpTask20260722{},
		&clickUpScopeConfig20260722{},
	)
}

func (*addFolderScope) Version() uint64 {
	return 20260722000001
}

func (*addFolderScope) Name() string {
	return "clickup: add folder scope + sprint/story-point columns"
}
