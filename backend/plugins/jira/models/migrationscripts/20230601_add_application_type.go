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
)

type scopeConfig20230601 struct {
	ApplicationType string `gorm:"type:varchar(255)"`
}

func (scopeConfig20230601) TableName() string {
	return "_tool_jira_scope_configs"
}

type addApplicationType struct{}

func (script *addApplicationType) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes, &scopeConfig20230601{})
}

func (*addApplicationType) Version() uint64 {
	return 20230601213738
}

func (*addApplicationType) Name() string {
	return "add application_type to _tool_jira_scope_configs"
}
