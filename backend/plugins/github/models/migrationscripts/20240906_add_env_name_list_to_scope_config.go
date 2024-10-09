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

var _ plugin.MigrationScript = (*addEnvNameListToScopeConfig)(nil)

type scopeConfig20240906 struct {
	EnvNameList []string `gorm:"type:json;serializer:json" json:"env_name_list" mapstructure:"env_name_list"`
}

func (scopeConfig20240906) TableName() string {
	return "_tool_github_scope_configs"
}

type addEnvNameListToScopeConfig struct{}

func (*addEnvNameListToScopeConfig) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if err := db.AutoMigrate(&scopeConfig20240906{}); err != nil {
		return err
	}
	return nil
}

func (*addEnvNameListToScopeConfig) Version() uint64 {
	return 20240906142100
}

func (*addEnvNameListToScopeConfig) Name() string {
	return "add env_name_list to _tool_github_scope_configs"
}
