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
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*addTableScopeIdToIncident)(nil)

type incident20240731 struct {
	Table   string `gorm:"index:idx_table_scope_id;type:varchar(255)"`
	ScopeId string `gorm:"index:idx_table_scope_id;type:varchar(255)"`
}

func (incident20240731) TableName() string {
	return "incidents"
}

type addTableScopeIdToIncident struct{}

func (*addTableScopeIdToIncident) Up(basicRes context.BasicRes) errors.Error {
	return basicRes.GetDal().AutoMigrate(&incident20240731{})
}

func (*addTableScopeIdToIncident) Version() uint64 {
	return 20240731180000
}

func (*addTableScopeIdToIncident) Name() string {
	return "add table and scope_id to incidents"
}
