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

type renameTr2ScopeConfig struct {
}

type scopeConfig20230529 struct {
	Entities []string `gorm:"type:json" json:"entities"`
}

func (scopeConfig20230529) TableName() string {
	return "_tool_gitlab_scope_configs"
}

func (u *renameTr2ScopeConfig) Up(baseRes context.BasicRes) errors.Error {
	db := baseRes.GetDal()
	err := db.RenameColumn("_tool_gitlab_projects", "transformation_rule_id", "scope_config_id")
	if err != nil {
		return err
	}
	err = db.RenameTable("_tool_gitlab_transformation_rules", "_tool_gitlab_scope_configs")
	if err != nil {
		return err
	}
	return migrationhelper.AutoMigrateTables(baseRes, &scopeConfig20230529{})
}

func (*renameTr2ScopeConfig) Version() uint64 {
	return 20230529173435
}

func (*renameTr2ScopeConfig) Name() string {
	return "rename transformation rule to scope config for gitlab"
}
