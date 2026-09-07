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
	"encoding/json"

	"github.com/apache/devlake/core/context"
	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/models/migrationscripts/archived"
	"github.com/apache/devlake/helpers/migrationhelper"
)

// taigaScopeConfig20260727 mirrors models.TaigaScopeConfig. The initial
// migration created `_tool_taiga_scope_configs` without the `type_mappings`
// column, while the runtime model declares it — every read/write of the model
// would fail with "Unknown column 'type_mappings'".
//
// The `uniqueIndex` on `name` is safe to add here: the new column is nullable,
// so pre-existing rows are backfilled with NULL, and both MySQL and PostgreSQL
// allow duplicate NULLs in a unique index (verified against both engines).
type taigaScopeConfig20260727 struct {
	archived.Model
	Entities     []string                   `gorm:"type:json;serializer:json" json:"entities"`
	ConnectionId uint64                     `json:"connectionId" gorm:"index"`
	Name         string                     `json:"name" gorm:"type:varchar(255);uniqueIndex"`
	TypeMappings map[string]json.RawMessage `json:"typeMappings" gorm:"type:json;serializer:json"`
}

func (taigaScopeConfig20260727) TableName() string {
	return "_tool_taiga_scope_configs"
}

type addMissingScopeConfigColumns struct{}

func (script *addMissingScopeConfigColumns) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes, &taigaScopeConfig20260727{})
}

func (*addMissingScopeConfigColumns) Version() uint64 {
	return 20260727000001
}

func (*addMissingScopeConfigColumns) Name() string {
	return "add missing type_mappings column to _tool_taiga_scope_configs"
}
