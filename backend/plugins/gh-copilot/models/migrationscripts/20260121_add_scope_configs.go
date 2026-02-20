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
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type addScopeConfig20260121 struct{}

type scopeConfig20260121 struct {
	common.ScopeConfig `mapstructure:",squash" json:",inline" gorm:"embedded"`
	ImplementationDate *time.Time `json:"implementationDate" mapstructure:"implementationDate" gorm:"type:datetime"`
	BaselinePeriodDays int        `json:"baselinePeriodDays" mapstructure:"baselinePeriodDays" gorm:"default:90"`
}

func (scopeConfig20260121) TableName() string {
	return "_tool_copilot_scope_configs"
}

func (*addScopeConfig20260121) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes, &scopeConfig20260121{})
}

func (*addScopeConfig20260121) Version() uint64 {
	return 20260121000000
}

func (*addScopeConfig20260121) Name() string {
	return "Add scope_configs table for GitHub Copilot impact analysis"
}
