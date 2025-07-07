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

type TestmoAutomationRun struct {
	ConnectionId uint64 `gorm:"primaryKey;type:BIGINT NOT NULL"`
	Id           uint64 `gorm:"primaryKey;type:BIGINT NOT NULL;autoIncrement:false" json:"id"`
	ProjectId    uint64 `gorm:"index;type:BIGINT NOT NULL" json:"project_id"`
	SourceId     uint64 `json:"source_id"`
	Name         string `gorm:"type:varchar(255)" json:"name"`
	Status       int32  `json:"status"`
	ConfigId     uint64 `json:"config_id"`
	MilestoneId  uint64 `json:"milestone_id"`
	Elapsed      *int64 `json:"elapsed"`
	IsCompleted  bool   `json:"is_completed"`

	// Test counts by status
	UntestedCount uint64 `json:"untested_count"`
	Status1Count  uint64 `json:"status1_count"`
	Status2Count  uint64 `json:"status2_count"`
	Status3Count  uint64 `json:"status3_count"`
	Status4Count  uint64 `json:"status4_count"`
	Status5Count  uint64 `json:"status5_count"`
	Status6Count  uint64 `json:"status6_count"`
	Status7Count  uint64 `json:"status7_count"`
	Status8Count  uint64 `json:"status8_count"`
	Status9Count  uint64 `json:"status9_count"`
	Status10Count uint64 `json:"status10_count"`
	Status11Count uint64 `json:"status11_count"`
	Status12Count uint64 `json:"status12_count"`
	Status13Count uint64 `json:"status13_count"`
	Status14Count uint64 `json:"status14_count"`
	Status15Count uint64 `json:"status15_count"`
	Status16Count uint64 `json:"status16_count"`
	Status17Count uint64 `json:"status17_count"`
	Status18Count uint64 `json:"status18_count"`
	Status19Count uint64 `json:"status19_count"`
	Status20Count uint64 `json:"status20_count"`
	Status21Count uint64 `json:"status21_count"`
	Status22Count uint64 `json:"status22_count"`
	Status23Count uint64 `json:"status23_count"`
	Status24Count uint64 `json:"status24_count"`

	// Aggregate counts
	SuccessCount   uint64 `json:"success_count"`
	FailureCount   uint64 `json:"failure_count"`
	CompletedCount uint64 `json:"completed_count"`
	TotalCount     uint64 `json:"total_count"`

	// Thread counts
	ThreadCount          uint64 `json:"thread_count"`
	ThreadActiveCount    uint64 `json:"thread_active_count"`
	ThreadCompletedCount uint64 `json:"thread_completed_count"`

	// Timestamps
	TestmoCreatedAt *time.Time `json:"created_at"`
	CreatedBy       uint64     `json:"created_by"`
	TestmoUpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy       *uint64    `json:"updated_by"`
	CompletedAt     *time.Time `json:"completed_at"`
	CompletedBy     *uint64    `json:"completed_by"`

	archived.NoPKModel
}

func (TestmoAutomationRun) TableName() string {
	return "_tool_testmo_automation_runs"
}
