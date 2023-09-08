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

package models

import (
	"github.com/apache/incubator-devlake/core/models/common"
)

type ZentaoScopeConfig struct {
	common.ScopeConfig  `mapstructure:",squash" json:",inline" gorm:"embedded"`
	ConnectionId        uint64            `mapstructure:"connectionId" json:"connectionId"`
	Name                string            `gorm:"type:varchar(255);index:idx_name_tapd,unique" validate:"required" mapstructure:"name" json:"name"`
	TypeMappings        map[string]string `mapstructure:"typeMappings,omitempty" json:"typeMappings" gorm:"serializer:json"`
	BugStatusMappings   map[string]string `mapstructure:"bugStatusMappings,omitempty" json:"bugStatusMappings" gorm:"serializer:json"`
	StoryStatusMappings map[string]string `mapstructure:"storyStatusMappings,omitempty" json:"storyStatusMappings" gorm:"serializer:json"`
	TaskStatusMappings  map[string]string `mapstructure:"taskStatusMappings,omitempty" json:"taskStatusMappings" gorm:"serializer:json"`
}

func (t ZentaoScopeConfig) TableName() string {
	return "_tool_zentao_scope_configs"
}
