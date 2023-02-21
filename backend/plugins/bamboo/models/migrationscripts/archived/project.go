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
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
)

type BambooProject struct {
	ConnectionId         uint64 `json:"connectionId" mapstructure:"connectionId" gorm:"primaryKey"`
	ProjectKey           string `json:"projectKey" gorm:"primaryKey;type:varchar(256)"`
	TransformationRuleId uint64 `json:"transformationRuleId,omitempty" mapstructure:"transformationRuleId"`
	Name                 string `json:"name" gorm:"index;type:varchar(256)"`
	Description          string `json:"description"`
	Href                 string `json:"link"`
	Rel                  string `json:"rel" gorm:"type:varchar(100)"`
	archived.NoPKModel   `json:"-" mapstructure:"-"`
}

func (BambooProject) TableName() string {
	return "_tool_bamboo_projects"
}
