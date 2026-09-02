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

// ClickUpTask is the tool-layer representation of a ClickUp task.
type ClickUpTask struct {
	ConnectionId uint64     `gorm:"primaryKey"`
	Id           string     `gorm:"primaryKey;type:varchar(255)" json:"id"`
	ListId       string     `gorm:"index;type:varchar(255)" json:"listId"`
	FolderId     string     `gorm:"index;type:varchar(255)" json:"folderId"`
	SpaceId      string     `gorm:"type:varchar(255)" json:"spaceId"`
	CustomId     string     `gorm:"type:varchar(255)" json:"customId"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Status       string     `gorm:"type:varchar(255)" json:"status"`
	StatusType   string     `gorm:"type:varchar(100)" json:"statusType"`
	Type         string     `gorm:"type:varchar(100)" json:"type"`
	Priority     string     `gorm:"type:varchar(100)" json:"priority"`
	Url          string     `gorm:"type:varchar(255)" json:"url"`
	CreatorId    string     `gorm:"type:varchar(255)" json:"creatorId"`
	AssigneeId   string     `gorm:"type:varchar(255)" json:"assigneeId"`
	AssigneeName string     `gorm:"type:varchar(255)" json:"assigneeName"`
	ParentId     string     `gorm:"type:varchar(255)" json:"parentId"`
	StoryPoint   *float64   `json:"storyPoint"`
	CreatedDate  *time.Time `json:"createdDate"`
	UpdatedDate  *time.Time `gorm:"index" json:"updatedDate"`
	ClosedDate   *time.Time `json:"closedDate"`
	common.NoPKModel
}

func (ClickUpTask) TableName() string {
	return "_tool_clickup_tasks"
}
