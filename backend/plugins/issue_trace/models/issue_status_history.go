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

// IssueStatusHistory records issue status history (status original value)
// end_date of current status is set to now() to avoid false assumption of future status.
// handled by ConvertIssueStatusHistory task
type IssueStatusHistory struct {
	common.NoPKModel
	IssueId           string     `gorm:"primaryKey;type:varchar(255)"`
	Status            string     `gorm:"type:varchar(100)"`
	OriginalStatus    string     `gorm:"primaryKey;type:varchar(255)"`
	StartDate         time.Time  `gorm:"primaryKey"`
	EndDate           *time.Time `gorm:"type:timestamp"`
	IsCurrentStatus   bool       `gorm:"type:boolean"`
	IsFirstStatus     bool       `gorm:"type:boolean"`
	StatusTimeMinutes int32      `gorm:"type:integer"`
}

func (IssueStatusHistory) TableName() string {
	return "issue_status_history"
}
