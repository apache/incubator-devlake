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

type TestmoTest struct {
	ConnectionId    uint64 `gorm:"primaryKey;type:BIGINT NOT NULL"`
	Id              uint64 `gorm:"primaryKey;type:BIGINT NOT NULL;autoIncrement:false" json:"id"`
	ProjectId       uint64 `gorm:"index;type:BIGINT NOT NULL" json:"project_id"`
	AutomationRunId uint64 `gorm:"index;type:BIGINT  NOT NULL" json:"automation_run_id"`
	ThreadId        uint64 `json:"thread_id"`
	Name            string `gorm:"type:varchar(500)" json:"name"`
	Key             string `gorm:"type:varchar(255)" json:"key"`
	Status          int32  `json:"status"`
	StatusName      string `gorm:"type:varchar(100)" json:"status_name"`
	Elapsed         *int64 `json:"elapsed"`
	Message         string `gorm:"type:text" json:"message"`

	// Test classification
	IsAcceptanceTest bool   `gorm:"index" json:"is_acceptance_test"`
	IsSmokeTest      bool   `gorm:"index" json:"is_smoke_test"`
	Team             string `gorm:"type:varchar(255);index" json:"team"`

	// Timestamps
	TestmoCreatedAt *time.Time `json:"created_at"`
	TestmoUpdatedAt *time.Time `json:"updated_at"`

	archived.NoPKModel
}

func (TestmoTest) TableName() string {
	return "_tool_testmo_tests"
}
