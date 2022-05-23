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
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type TapdTask struct {
	ConnectionId    uint64          `gorm:"primaryKey"`
	ID              uint64          `gorm:"primaryKey;type:BIGINT" json:"id,string"`
	Name            string          `gorm:"type:varchar(255)" json:"name"`
	Description     string          `json:"description"`
	WorkspaceID     uint64          `json:"workspace_id,string"`
	Creator         string          `gorm:"type:varchar(255)" json:"creator"`
	Created         *helper.CSTTime `json:"created"`
	Modified        *helper.CSTTime `json:"modified" gorm:"index"`
	Status          string          `json:"status" gorm:"type:varchar(255)"`
	Owner           string          `json:"owner" gorm:"type:varchar(255)"`
	Cc              string          `json:"cc" gorm:"type:varchar(255)"`
	Begin           *helper.CSTTime `json:"begin"`
	Due             *helper.CSTTime `json:"due"`
	Priority        string          `gorm:"type:varchar(255)" json:"priority"`
	IterationID     uint64          `json:"iteration_id,string"`
	Completed       *helper.CSTTime `json:"completed"`
	Effort          float32         `json:"effort,string"`
	EffortCompleted float32         `json:"effort_completed,string"`
	Exceed          float32         `json:"exceed,string"`
	Remain          float32         `json:"remain,string"`
	StdStatus       string
	StdType         string
	Type            string
	StoryID         uint64 `json:"story_id,string"`
	Progress        int16  `json:"progress,string"`
	HasAttachment   string `gorm:"type:varchar(255)"`
	Url             string

	AttachmentCount  int16  `json:"attachment_count,string"`
	Follower         string `json:"follower" gorm:"type:varchar(255)"`
	CreatedFrom      string `json:"created_from" gorm:"type:varchar(255)"`
	PredecessorCount int16  `json:"predecessor_count,string"`
	SuccessorCount   int16  `json:"successor_count,string"`
	ReleaseId        uint64 `json:"release_id,string"`
	Label            string `json:"label" gorm:"type:varchar(255)"`
	NewStoryId       uint64 `json:"new_story_id,string"`
	common.NoPKModel
}

func (TapdTask) TableName() string {
	return "_tool_tapd_tasks"
}
