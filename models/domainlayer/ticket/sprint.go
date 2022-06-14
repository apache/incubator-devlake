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
	"github.com/apache/incubator-devlake/models/common"
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer"
)

var (
	BeforeSprint = "BEFORE_SPRINT"
	DuringSprint = "DURING_SPRINT"
	AfterSprint  = "AFTER_SPRINT"
)

type Sprint struct {
	domainlayer.DomainEntity
	Name            string `gorm:"type:varchar(255)"`
	Url             string `gorm:"type:varchar(255)"`
	Status          string `gorm:"type:varchar(100)"`
	StartedDate     *time.Time
	EndedDate       *time.Time
	CompletedDate   *time.Time
	OriginalBoardID string `gorm:"type:varchar(255)"`
}

func (Sprint) TableName() string {
	return "sprints"
}

type SprintIssue struct {
	common.NoPKModel
	SprintId string `gorm:"primaryKey;type:varchar(255)"`
	IssueId  string `gorm:"primaryKey;type:varchar(255)"`
}

func (SprintIssue) TableName() string {
	return "sprint_issues"
}
