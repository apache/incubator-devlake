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

package ticket

import (
	"time"

	"github.com/apache/incubator-devlake/models/common"
)

type IssueStatusHistory struct {
	common.NoPKModel
	IssueId        string    `gorm:"primaryKey;type:varchar(255)"`
	OriginalStatus string    `gorm:"primaryKey;type:varchar(255)"`
	StartDate      time.Time `gorm:"primaryKey"`
	EndDate        *time.Time
}

func (IssueStatusHistory) TableName() string {
	return "issue_status_history"
}

type IssueAssigneeHistory struct {
	common.NoPKModel
	IssueId   string    `gorm:"primaryKey;type:varchar(255)"`
	Assignee  string    `gorm:"primaryKey;type:varchar(255)"`
	StartDate time.Time `gorm:"primaryKey"`
	EndDate   *time.Time
}

func (IssueAssigneeHistory) TableName() string {
	return "issue_assignee_history"
}

type IssueSprintsHistory struct {
	common.NoPKModel
	IssueId   string    `gorm:"primaryKey;type:varchar(255)"`
	SprintId  string    `gorm:"primaryKey;type:varchar(255)"`
	StartDate time.Time `gorm:"primaryKey"`
	EndDate   *time.Time
}

func (IssueSprintsHistory) TableName() string {
	return "issue_sprints_history"
}
