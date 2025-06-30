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
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
)

type TestmoMilestone struct {
	ConnectionId uint64 `gorm:"primaryKey;type:BIGINT NOT NULL"`
	Id           uint64 `gorm:"primaryKey;type:BIGINT NOT NULL;autoIncrement:false" json:"id"`
	ProjectId    uint64 `gorm:"index;type:BIGINT  NOT NULL" json:"project_id"`
	Name         string `gorm:"type:varchar(255)" json:"name"`
	Description  string `gorm:"type:text" json:"description"`
	IsCompleted  bool   `json:"is_completed"`

	// Timestamps
	TestmoCreatedAt *time.Time `json:"created_at"`
	TestmoUpdatedAt *time.Time `json:"updated_at"`
	CompletedAt     *time.Time `json:"completed_at"`

	common.NoPKModel
}

func (TestmoMilestone) TableName() string {
	return "_tool_testmo_milestones"
}
