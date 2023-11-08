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
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/zentao/models/migrationscripts/archived"
)

type addScopeConfigTables struct{}

type bug20230601 struct {
	Url       string `json:"url"`
	StdStatus string `json:"stdStatus" gorm:"type:varchar(20)"`
	StdType   string `json:"stdType" gorm:"type:varchar(20)"`
}

func (bug20230601) TableName() string {
	return "_tool_zentao_bugs"
}

type story20230601 struct {
	Url       string `json:"url"`
	StdStatus string `json:"stdStatus" gorm:"type:varchar(20)"`
	StdType   string `json:"stdType" gorm:"type:varchar(20)"`
}

func (story20230601) TableName() string {
	return "_tool_zentao_stories"
}

type task20230601 struct {
	Url       string `json:"url"`
	StdStatus string `json:"stdStatus" gorm:"type:varchar(20)"`
	StdType   string `json:"stdType" gorm:"type:varchar(20)"`
}

func (task20230601) TableName() string {
	return "_tool_zentao_tasks"
}

type project20230602 struct {
	ScopeConfigId uint64 `json:"scopeConfigId,omitempty" mapstructure:"scopeConfigId"`
}

func (project20230602) TableName() string {
	return "_tool_zentao_projects"
}

type product20230602 struct {
	ScopeConfigId uint64 `json:"scopeConfigId,omitempty" mapstructure:"scopeConfigId"`
}

func (product20230602) TableName() string {
	return "_tool_zentao_products"
}

func (*addScopeConfigTables) Up(basicRes context.BasicRes) errors.Error {

	return migrationhelper.AutoMigrateTables(
		basicRes,
		&archived.ZentaoScopeConfig{},
		&bug20230601{},
		&story20230601{},
		&task20230601{},
		&project20230602{},
		&product20230602{},
	)
}

func (*addScopeConfigTables) Version() uint64 {
	return 20230602000001
}

func (*addScopeConfigTables) Name() string {
	return "zentao add scope config tables"
}
