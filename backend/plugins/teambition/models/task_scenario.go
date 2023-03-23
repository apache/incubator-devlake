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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type TeambitionTaskScenario struct {
	ConnectionId      uint64                    `gorm:"primaryKey;type:BIGINT"`
	ProjectId         string                    `gorm:"primaryKey;type:varchar(100)" json:"projectId"`
	Id                string                    `gorm:"primaryKey;type:varchar(100)" json:"id"`
	Name              string                    `gorm:"type:varchar(100)" json:"name"`
	Icon              string                    `gorm:"type:varchar(100)" json:"icon"`
	Type              string                    `gorm:"type:varchar(100)" json:"type"`
	Scenariofields    []TeambitionScenarioField `gorm:"serializer:json;type:text" json:"scenariofields"`
	CreatorId         string                    `gorm:"type:varchar(100)" json:"creatorId"`
	Source            string                    `gorm:"type:varchar(100)" json:"source"`
	BoundToObjectId   string                    `gorm:"type:varchar(100)" json:"boundToObjectId"`
	BoundToObjectType string                    `gorm:"type:varchar(100)" json:"boundToObjectType"`
	TaskflowId        string                    `gorm:"type:varchar(100)" json:"taskflowId"`
	IsArchived        bool                      `json:"isArchived"`
	Created           *api.Iso8601Time          `json:"created"`
	Updated           *api.Iso8601Time          `json:"updated"`

	common.NoPKModel
}

type TeambitionScenarioField struct {
	CustomfieldId string `json:"customfieldId"`
	FieldType     string `json:"fieldType"`
	Required      bool   `json:"required"`
}

func (TeambitionTaskScenario) TableName() string {
	return "_tool_teambition_task_scenarios"
}
