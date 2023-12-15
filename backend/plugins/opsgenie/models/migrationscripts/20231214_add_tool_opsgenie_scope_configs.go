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
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*addOpsenieScopeConfig20231214)(nil)

type OpsenieScopeConfig20231214 struct {
	archived.ScopeConfig `mapstructure:",squash" json:",inline" gorm:"embedded"`
	Name                 string `gorm:"type:varchar(255);index:idx_name_opsgenie,unique" validate:"required" mapstructure:"name" json:"name"`
}

func (o OpsenieScopeConfig20231214) TableName() string {
	return "_tool_opsgenie_scope_configs"
}

type addOpsenieScopeConfig20231214 struct{}

func (script *addOpsenieScopeConfig20231214) Up(basicRes context.BasicRes) errors.Error {
	return basicRes.GetDal().AutoMigrate(&OpsenieScopeConfig20231214{})
}

func (*addOpsenieScopeConfig20231214) Version() uint64 {
	return 20231214160000
}

func (script *addOpsenieScopeConfig20231214) Name() string {
	return "init table _tool_opsgenie_scope_configs"
}
