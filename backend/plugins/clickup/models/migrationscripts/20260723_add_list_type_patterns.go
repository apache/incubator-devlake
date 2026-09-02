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
	"github.com/apache/devlake/core/context"
	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/models/migrationscripts/archived"
	"github.com/apache/devlake/helpers/migrationhelper"
)

// addListTypePatterns adds the list-name typing columns to
// _tool_clickup_scope_configs. ClickUp often carries no per-task type, so bugs
// live in a dedicated list (e.g. "QA Bugs"); these patterns let a folder scope
// classify a whole list's tasks as BUG/INCIDENT by list name. AutoMigrate is
// additive.

type clickUpScopeConfig20260723 struct {
	archived.ScopeConfig
	BugListPattern      string `gorm:"type:varchar(255)"`
	IncidentListPattern string `gorm:"type:varchar(255)"`
}

func (clickUpScopeConfig20260723) TableName() string { return "_tool_clickup_scope_configs" }

type addListTypePatterns struct{}

func (*addListTypePatterns) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&clickUpScopeConfig20260723{},
	)
}

func (*addListTypePatterns) Version() uint64 {
	return 20260723000001
}

func (*addListTypePatterns) Name() string {
	return "clickup: add bug/incident list-name type patterns to scope config"
}
