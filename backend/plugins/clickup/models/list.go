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

	"github.com/apache/devlake/core/models/common"
)

// ClickUpList is a list inside a scoped folder. It is NOT itself a scope (the
// folder is — see ClickUpFolder). A list is one of two things:
//
//   - a sprint list (IsSprint) — a rolling sprint, converted to ticket.Sprint;
//     tasks in it become sprint_issues. ClickUp teams encode sprints as lists
//     named e.g. "v4.3.0 Sprint 40 (7/6/26 - 7/19/26)"; the sprint number and
//     start/end dates are parsed from that name (there are no list date fields).
//   - a regular list (Backlog / Bug Tracking / DevOps) — its tasks are plain
//     board issues, no sprint.
type ClickUpList struct {
	ConnectionId uint64     `gorm:"primaryKey" json:"connectionId"`
	ListId       string     `gorm:"primaryKey;type:varchar(255)" json:"listId"`
	FolderId     string     `gorm:"index;type:varchar(255)" json:"folderId"`
	SpaceId      string     `gorm:"type:varchar(255)" json:"spaceId"`
	Name         string     `gorm:"type:varchar(255)" json:"name"`
	Archived     bool       `json:"archived"`
	IsSprint     bool       `json:"isSprint"`
	SprintName   string     `gorm:"type:varchar(255)" json:"sprintName"`
	StartDate    *time.Time `json:"startDate"`
	EndDate      *time.Time `json:"endDate"`
	common.NoPKModel
}

func (ClickUpList) TableName() string {
	return "_tool_clickup_lists"
}
