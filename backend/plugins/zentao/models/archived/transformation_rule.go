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

package archived

import (
	"encoding/json"

	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
)

type ZentaoTransformationRule struct {
	archived.Model      `mapstructure:"-"`
	ConnectionId        uint64          `mapstructure:"connectionId" json:"connectionId"`
	Name                string          `gorm:"type:varchar(255);index:idx_name_zentao,unique" validate:"required" mapstructure:"name" json:"name"`
	BugStatusMappings   json.RawMessage `mapstructure:"bugStatusMappings,omitempty" json:"bugStatusMappings"`
	StoryStatusMappings json.RawMessage `mapstructure:"storyStatusMappings,omitempty" json:"storyStatusMappings"`
	TaskStatusMappings  json.RawMessage `mapstructure:"taskStatusMappings,omitempty" json:"taskStatusMappings"`
}

func (t ZentaoTransformationRule) TableName() string {
	return "_tool_zentao_transformation_rules"
}
