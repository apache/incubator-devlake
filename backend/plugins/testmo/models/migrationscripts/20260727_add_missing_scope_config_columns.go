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
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

// testmoScopeConfig20260727 mirrors models.TestmoScopeConfig. The migration
// that created `_tool_testmo_scope_configs` omitted the `connection_id` and
// `name` columns of the embedded common.ScopeConfig, which the runtime model
// expects.
//
// The `uniqueIndex` on `name` is safe to add here: the new column is nullable,
// so pre-existing rows are backfilled with NULL, and both MySQL and PostgreSQL
// allow duplicate NULLs in a unique index (verified against both engines).
type testmoScopeConfig20260727 struct {
	archived.Model
	Entities              []string `gorm:"type:json;serializer:json" json:"entities"`
	ConnectionId          uint64   `json:"connectionId" gorm:"index"`
	Name                  string   `json:"name" gorm:"type:varchar(255);uniqueIndex"`
	AcceptanceTestPattern string   `json:"acceptanceTestPattern" gorm:"type:varchar(255)"`
	SmokeTestPattern      string   `json:"smokeTestPattern" gorm:"type:varchar(255)"`
	TeamPattern           string   `json:"teamPattern" gorm:"type:varchar(255)"`
}

func (testmoScopeConfig20260727) TableName() string {
	return "_tool_testmo_scope_configs"
}

type addMissingScopeConfigColumns struct{}

func (script *addMissingScopeConfigColumns) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes, &testmoScopeConfig20260727{})
}

func (*addMissingScopeConfigColumns) Version() uint64 {
	return 20260727000001
}

func (*addMissingScopeConfigColumns) Name() string {
	return "add missing connection_id/name columns to _tool_testmo_scope_configs"
}
