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

var _ plugin.MigrationScript = (*addFilterFieldsToPagerDutyScopeConfig20251003)(nil)

type PagerDutyScopeConfig20251003 struct {
	PriorityFilter []string `mapstructure:"priorityFilter" json:"priorityFilter" gorm:"type:text;serializer:json"`
	UrgencyFilter  []string `mapstructure:"urgencyFilter" json:"urgencyFilter" gorm:"type:text;serializer:json"`
}

func (o PagerDutyScopeConfig20251003) TableName() string {
	return "_tool_pagerduty_scope_configs"
}

type addFilterFieldsToPagerDutyScopeConfig20251003 struct{}

func (script *addFilterFieldsToPagerDutyScopeConfig20251003) Up(basicRes context.BasicRes) errors.Error {
	return basicRes.GetDal().AutoMigrate(&PagerDutyScopeConfig20251003{})
}

func (*addFilterFieldsToPagerDutyScopeConfig20251003) Version() uint64 {
	return 20251003000000
}

func (script *addFilterFieldsToPagerDutyScopeConfig20251003) Name() string {
	return "add priority_filter and urgency_filter fields to table _tool_pagerduty_scope_configs"
}
