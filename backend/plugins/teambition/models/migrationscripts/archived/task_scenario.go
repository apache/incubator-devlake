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
	"time"

	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
)

type TeambitionTaskScenario struct {
	ConnectionId      uint64                    `gorm:"primaryKey;type:BIGINT"`
	ProjectId         string                    `gorm:"primaryKey;type:varchar(100)"`
	Id                string                    `gorm:"primaryKey;type:varchar(100)"`
	Name              string                    `gorm:"type:varchar(100)"`
	Icon              string                    `gorm:"type:varchar(100)"`
	Type              string                    `gorm:"type:varchar(100)"`
	Scenariofields    []TeambitionScenarioField `gorm:"serializer:json;type:text"`
	CreatorId         string                    `gorm:"type:varchar(100)"`
	Source            string                    `gorm:"type:varchar(100)"`
	BoundToObjectId   string                    `gorm:"type:varchar(100)"`
	BoundToObjectType string                    `gorm:"type:varchar(100)"`
	TaskflowId        string                    `gorm:"type:varchar(100)"`
	IsArchived        bool
	Created           *time.Time
	Updated           *time.Time

	archived.NoPKModel
}

type TeambitionScenarioField struct {
	CustomfieldId string
	FieldType     string
	Required      bool
}

func (TeambitionTaskScenario) TableName() string {
	return "_tool_teambition_task_scenarios"
}
