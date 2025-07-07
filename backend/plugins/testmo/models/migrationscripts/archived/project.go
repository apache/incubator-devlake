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

type TestmoProject struct {
	ConnectionId uint64 `gorm:"primaryKey;type:BIGINT NOT NULL"`
	Id           uint64 `gorm:"primaryKey;type:BIGINT NOT NULL;autoIncrement:false" json:"id"`
	Name         string `gorm:"type:varchar(255)" json:"name"`
	IsCompleted  bool   `json:"is_completed"`

	// Counts for metrics
	MilestoneCount               uint64 `json:"milestone_count"`
	MilestoneActiveCount         uint64 `json:"milestone_active_count"`
	MilestoneCompletedCount      uint64 `json:"milestone_completed_count"`
	RunCount                     uint64 `json:"run_count"`
	RunActiveCount               uint64 `json:"run_active_count"`
	RunClosedCount               uint64 `json:"run_closed_count"`
	AutomationSourceCount        uint64 `json:"automation_source_count"`
	AutomationSourceActiveCount  uint64 `json:"automation_source_active_count"`
	AutomationSourceRetiredCount uint64 `json:"automation_source_retired_count"`
	AutomationRunCount           uint64 `json:"automation_run_count"`
	AutomationRunActiveCount     uint64 `json:"automation_run_active_count"`
	AutomationRunCompletedCount  uint64 `json:"automation_run_completed_count"`

	archived.NoPKModel
}

func (TestmoProject) TableName() string {
	return "_tool_testmo_projects"
}
