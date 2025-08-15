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
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
)

// Define the actual table structure directly in the migration script
type issueLeadTimeMetric struct {
	ProjectName             string `gorm:"primaryKey;type:varchar(255)"`
	IssueId                 string `gorm:"primaryKey;type:varchar(255)"`
	InProgressDate          *time.Time
	DoneDate                *time.Time
	InProgressToDoneMinutes *int64
}

// TableName specifies the table name
func (issueLeadTimeMetric) TableName() string {
	return "_tool_dora_issue_lead_time_metrics"
}

type addIssueLeadTimeMetricsTable struct{}

func (script *addIssueLeadTimeMetricsTable) Up(baseRes context.BasicRes) errors.Error {
	db := baseRes.GetDal()
	// Use our directly defined model instead of importing from models
	return db.AutoMigrate(&issueLeadTimeMetric{})
}

func (*addIssueLeadTimeMetricsTable) Version() uint64 {
	return 2025042401
}

func (*addIssueLeadTimeMetricsTable) Name() string {
	return "dora add _tool_dora_issue_lead_time_metrics table"
}
