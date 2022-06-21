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
	ConnectionId     uint64          `gorm:"primaryKey"`
	Id               uint64          `gorm:"primaryKey;type:BIGINT" json:"id,string"`
	WorkitemTypeId   uint64          `json:"workitem_type_id,string"`
	Name             string          `gorm:"type:varchar(255)" json:"name"`
	Description      string          `json:"description"`
	WorkspaceId      uint64          `json:"workspace_id,string"`
	Creator          string          `gorm:"type:varchar(255)"`
	Created          *helper.CSTTime `json:"created"`
	Modified         *helper.CSTTime `json:"modified" gorm:"index"`
	Status           string          `json:"status" gorm:"type:varchar(255)"`
	Owner            string          `json:"owner" gorm:"type:varchar(255)"`
	Cc               string          `json:"cc" gorm:"type:varchar(255)"`
	Begin            *helper.CSTTime `json:"begin"`
	Due              *helper.CSTTime `json:"due"`
	Size             int16           `json:"size,string"`
	Priority         string          `gorm:"type:varchar(255)" json:"priority"`
	Developer        string          `gorm:"type:varchar(255)" json:"developer"`
	IterationId      uint64          `json:"iteration_id,string"`
	TestFocus        string          `json:"test_focus" gorm:"type:varchar(255)"`
	Type             string          `json:"type" gorm:"type:varchar(20)"`
	Source           string          `json:"source" gorm:"type:varchar(255)"`
	Module           string          `json:"module" gorm:"type:varchar(255)"`
	Version          string          `json:"version" gorm:"type:varchar(255)"`
	Completed        *helper.CSTTime `json:"completed"`
	CategoryId       int64           `json:"category_id,string"`
	Path             string          `gorm:"type:varchar(255)" json:"path"`
	ParentId         uint64          `json:"parent_id,string"`
	ChildrenId       string          `gorm:"type:text" json:"children_id"`
	AncestorId       uint64          `json:"ancestor_id,string"`
	BusinessValue    string          `gorm:"type:varchar(255)" json:"business_value"`
	Effort           float32         `json:"effort,string"`
	EffortCompleted  float32         `json:"effort_completed,string"`
	Exceed           float32         `json:"exceed,string"`
	Remain           float32         `json:"remain,string"`
	ReleaseId        uint64          `json:"release_id,string"`
	Confidential     string          `gorm:"type:varchar(255)" json:"confidential"`
	TemplatedId      uint64          `json:"templated_id,string"`
	CreatedFrom      string          `gorm:"type:varchar(255)" json:"created_from"`
	Feature          string          `gorm:"type:varchar(255)" json:"feature"`
	StdStatus        string          `gorm:"type:varchar(20)"`
	StdType          string          `gorm:"type:varchar(20)"`
	Url              string          `gorm:"type:varchar(255)"`
	
	AttachmentCount  int16           `json:"attachment_count,string"`
	HasAttachment    string          `json:"has_attachment" gorm:"type:varchar(255)"`
	BugId            uint64          `json:"bug_id,string"`
	Follower         string          `json:"follower" gorm:"type:varchar(255)"`
	SyncType         string          `json:"sync_type" gorm:"type:text"`
	PredecessorCount int16           `json:"predecessor_count,string"`
	IsArchived       string          `json:"is_archived" gorm:"type:varchar(255)"`
	Modifier         string          `json:"modifier" gorm:"type:varchar(255)"`
	ProgressManual   string          `json:"progress_manual" gorm:"type:varchar(255)"`
	SuccessorCount   int16           `json:"successor_count,string"`
	Label            string          `json:"label" gorm:"type:varchar(255)"`
	CustomFieldOne   string          `json:"custom_field_one" gorm:"type:text"`
	CustomFieldTwo   string          `json:"custom_field_two" gorm:"type:text"`
	CustomFieldThree string          `json:"custom_field_three" gorm:"type:text"`
	CustomFieldFour  string          `json:"custom_field_four" gorm:"type:text"`
	CustomFieldFive  string          `json:"custom_field_five" gorm:"type:text"`
	CustomField6     string          `json:"custom_field_6" gorm:"type:text"`
	CustomField7     string          `json:"custom_field_7" gorm:"type:text"`
	CustomField8     string          `json:"custom_field_8" gorm:"type:text"`
	CustomField9     string          `json:"custom_field_9" gorm:"type:text"`
	CustomField10    string          `json:"custom_field_10" gorm:"type:text"`
	CustomField11    string          `json:"custom_field_11" gorm:"type:text"`
	CustomField12    string          `json:"custom_field_12" gorm:"type:text"`
	CustomField13    string          `json:"custom_field_13" gorm:"type:text"`
	CustomField14    string          `json:"custom_field_14" gorm:"type:text"`
	CustomField15    string          `json:"custom_field_15" gorm:"type:text"`
	CustomField16    string          `json:"custom_field_16" gorm:"type:text"`
	CustomField17    string          `json:"custom_field_17" gorm:"type:text"`
	CustomField18    string          `json:"custom_field_18" gorm:"type:text"`
	CustomField19    string          `json:"custom_field_19" gorm:"type:text"`
	CustomField20    string          `json:"custom_field_20" gorm:"type:text"`
	CustomField21    string          `json:"custom_field_21" gorm:"type:text"`
	CustomField22    string          `json:"custom_field_22" gorm:"type:text"`
	CustomField23    string          `json:"custom_field_23" gorm:"type:text"`
	CustomField24    string          `json:"custom_field_24" gorm:"type:text"`
	CustomField25    string          `json:"custom_field_25" gorm:"type:text"`
	CustomField26    string          `json:"custom_field_26" gorm:"type:text"`
	CustomField27    string          `json:"custom_field_27" gorm:"type:text"`
	CustomField28    string          `json:"custom_field_28" gorm:"type:text"`
	CustomField29    string          `json:"custom_field_29" gorm:"type:text"`
	CustomField30    string          `json:"custom_field_30" gorm:"type:text"`
	CustomField31    string          `json:"custom_field_31" gorm:"type:text"`
	CustomField32    string          `json:"custom_field_32" gorm:"type:text"`
	CustomField33    string          `json:"custom_field_33" gorm:"type:text"`
	CustomField34    string          `json:"custom_field_34" gorm:"type:text"`
	CustomField35    string          `json:"custom_field_35" gorm:"type:text"`
	CustomField36    string          `json:"custom_field_36" gorm:"type:text"`
	CustomField37    string          `json:"custom_field_37" gorm:"type:text"`
	CustomField38    string          `json:"custom_field_38" gorm:"type:text"`
	CustomField39    string          `json:"custom_field_39" gorm:"type:text"`
	CustomField40    string          `json:"custom_field_40" gorm:"type:text"`
	CustomField41    string          `json:"custom_field_41" gorm:"type:text"`
	CustomField42    string          `json:"custom_field_42" gorm:"type:text"`
	CustomField43    string          `json:"custom_field_43" gorm:"type:text"`
	CustomField44    string          `json:"custom_field_44" gorm:"type:text"`
	CustomField45    string          `json:"custom_field_45" gorm:"type:text"`
	CustomField46    string          `json:"custom_field_46" gorm:"type:text"`
	CustomField47    string          `json:"custom_field_47" gorm:"type:text"`
	CustomField48    string          `json:"custom_field_48" gorm:"type:text"`
	CustomField49    string          `json:"custom_field_49" gorm:"type:text"`
	CustomField50    string          `json:"custom_field_50" gorm:"type:text"`

	common.NoPKModel
}

func (TapdStory) TableName() string {
	return "_tool_tapd_stories"
}
