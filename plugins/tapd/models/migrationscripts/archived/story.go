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
	"github.com/apache/incubator-devlake/models/common"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type TapdStory struct {
	ConnectionId    uint64          `gorm:"primaryKey"`
	ID              uint64          `gorm:"primaryKey;type:BIGINT" json:"id,string"`
	WorkitemTypeID  uint64          `json:"workitem_type_id,string"`
	Name            string          `gorm:"type:varchar(255)" json:"name"`
	Description     string          `json:"description"`
	WorkspaceID     uint64          `json:"workspace_id,string"`
	Creator         string          `gorm:"type:varchar(255)"`
	Created         *helper.CSTTime `json:"created"`
	Modified        *helper.CSTTime `json:"modified" gorm:"index"`
	Status          string          `json:"status" gorm:"type:varchar(255)"`
	Owner           string          `json:"owner" gorm:"type:varchar(255)"`
	Cc              string          `json:"cc" gorm:"type:varchar(255)"`
	Begin           *helper.CSTTime `json:"begin"`
	Due             *helper.CSTTime `json:"due"`
	Size            int16           `json:"size,string"`
	Priority        string          `gorm:"type:varchar(255)" json:"priority"`
	Developer       string          `gorm:"type:varchar(255)" json:"developer"`
	IterationID     uint64          `json:"iteration_id,string"`
	TestFocus       string          `json:"test_focus" gorm:"type:varchar(255)"`
	Type            string          `json:"type" gorm:"type:varchar(255)"`
	Source          string          `json:"source" gorm:"type:varchar(255)"`
	Module          string          `json:"module" gorm:"type:varchar(255)"`
	Version         string          `json:"version" gorm:"type:varchar(255)"`
	Completed       *helper.CSTTime `json:"completed"`
	CategoryID      uint64          `json:"category_id,string"`
	Path            string          `gorm:"type:varchar(255)" json:"path"`
	ParentID        uint64          `json:"parent_id,string"`
	ChildrenID      string          `gorm:"type:varchar(255)" json:"children_id"`
	AncestorID      uint64          `json:"ancestor_id,string"`
	BusinessValue   string          `gorm:"type:varchar(255)" json:"business_value"`
	Effort          float32         `json:"effort,string"`
	EffortCompleted float32         `json:"effort_completed,string"`
	Exceed          float32         `json:"exceed,string"`
	Remain          float32         `json:"remain,string"`
	ReleaseID       uint64          `json:"release_id,string"`
	Confidential    string          `gorm:"type:varchar(255)" json:"confidential"`
	TemplatedID     uint64          `json:"templated_id,string"`
	CreatedFrom     string          `gorm:"type:varchar(255)" json:"created_from"`
	Feature         string          `gorm:"type:varchar(255)" json:"feature"`
	StdStatus       string
	StdType         string
	Url             string

	AttachmentCount  int16  `json:"attachment_count,string"`
	HasAttachment    string `json:"has_attachment" gorm:"type:varchar(255)"`
	BugID            uint64 `json:"bug_id,string"`
	Follower         string `json:"follower" gorm:"type:varchar(255)"`
	SyncType         string `json:"sync_type" gorm:"type:varchar(255)"`
	PredecessorCount int16  `json:"predecessor_count,string"`
	IsArchived       string `json:"is_archived" gorm:"type:varchar(255)"`
	Modifier         string `json:"modifier" gorm:"type:varchar(255)"`
	ProgressManual   string `json:"progress_manual" gorm:"type:varchar(255)"`
	SuccessorCount   int16  `json:"successor_count,string"`
	Label            string `json:"label" gorm:"type:varchar(255)"`
	common.NoPKModel
}

func (TapdStory) TableName() string {
	return "_tool_tapd_stories"
}
