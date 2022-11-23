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

type TapdBug struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	Id           uint64 `gorm:"primaryKey;type:BIGINT" json:"id,string"`
	EpicKey      string
	Title        string          `json:"title" gorm:"type:varchar(255)"`
	Description  string          `json:"description"`
	WorkspaceId  uint64          `json:"workspace_id,string"`
	Created      *helper.CSTTime `json:"created"`
	Modified     *helper.CSTTime `json:"modified" gorm:"index"`
	Status       string          `json:"status" gorm:"type:varchar(255)"`
	Cc           string          `json:"cc" gorm:"type:varchar(255)"`
	Begin        *helper.CSTTime `json:"begin"`
	Due          *helper.CSTTime `json:"due"`
	Priority     string          `json:"priority" gorm:"type:varchar(255)"`
	IterationId  uint64          `json:"iteration_id,string"`
	Source       string          `json:"source" gorm:"type:varchar(255)"`
	Module       string          `json:"module" gorm:"type:varchar(255)"`
	ReleaseId    uint64          `json:"release_id,string"`
	CreatedFrom  string          `json:"created_from" gorm:"type:varchar(255)"`
	Feature      string          `json:"feature" gorm:"type:varchar(255)"`
	common.NoPKModel

	Severity         string          `json:"severity" gorm:"type:varchar(255)"`
	Reporter         string          `json:"reporter" gorm:"type:varchar(255)"`
	Resolved         *helper.CSTTime `json:"resolved"`
	Closed           *helper.CSTTime `json:"closed"`
	Lastmodify       string          `json:"lastmodify" gorm:"type:varchar(255)"`
	Auditer          string          `json:"auditer" gorm:"type:varchar(255)"`
	De               string          `json:"De" gorm:"comment:developer;type:varchar(255)"`
	Fixer            string          `json:"fixer" gorm:"type:varchar(255)"`
	VersionTest      string          `json:"version_test" gorm:"type:varchar(255)"`
	VersionReport    string          `json:"version_report" gorm:"type:varchar(255)"`
	VersionClose     string          `json:"version_close" gorm:"type:varchar(255)"`
	VersionFix       string          `json:"version_fix" gorm:"type:varchar(255)"`
	BaselineFind     string          `json:"baseline_find" gorm:"type:varchar(255)"`
	BaselineJoin     string          `json:"baseline_join" gorm:"type:varchar(255)"`
	BaselineClose    string          `json:"baseline_close" gorm:"type:varchar(255)"`
	BaselineTest     string          `json:"baseline_test" gorm:"type:varchar(255)"`
	Sourcephase      string          `json:"sourcephase" gorm:"type:varchar(255)"`
	Te               string          `json:"te" gorm:"type:varchar(255)"`
	CurrentOwner     string          `json:"current_owner" gorm:"type:varchar(255)"`
	Resolution       string          `json:"resolution" gorm:"type:varchar(255)"`
	Originphase      string          `json:"originphase" gorm:"type:varchar(255)"`
	Confirmer        string          `json:"confirmer" gorm:"type:varchar(255)"`
	Participator     string          `json:"participator" gorm:"type:varchar(255)"`
	Closer           string          `json:"closer" gorm:"type:varchar(50)"`
	Platform         string          `json:"platform" gorm:"type:varchar(50)"`
	Os               string          `json:"os" gorm:"type:varchar(50)"`
	Testtype         string          `json:"testtype" gorm:"type:varchar(20)"`
	Testphase        string          `json:"testphase" gorm:"type:varchar(255)"`
	Frequency        string          `json:"frequency" gorm:"type:varchar(255)"`
	RegressionNumber string          `json:"regression_number" gorm:"type:varchar(20)"`
	Flows            string          `json:"flows" gorm:"type:varchar(255)"`
	Testmode         string          `json:"testmode" gorm:"type:varchar(50)"`
	IssueId          uint64          `json:"issue_id,string"`
	VerifyTime       *helper.CSTTime `json:"verify_time"`
	RejectTime       *helper.CSTTime `json:"reject_time"`
	ReopenTime       *helper.CSTTime `json:"reopen_time"`
	AuditTime        *helper.CSTTime `json:"audit_time"`
	SuspendTime      *helper.CSTTime `json:"suspend_time"`
	Deadline         *helper.CSTTime `json:"deadline"`
	InProgressTime   *helper.CSTTime `json:"in_progress_time"`
	AssignedTime     *helper.CSTTime `json:"assigned_time"`
	TemplateId       uint64          `json:"template_id,string"`
	StoryId          uint64          `json:"story_id,string"`
	StdStatus        string          `gorm:"type:varchar(20)"`
	StdType          string          `gorm:"type:varchar(20)"`
	Type             string          `gorm:"type:varchar(20)"`
	Url              string          `gorm:"type:varchar(255)"`

	SupportId       uint64  `json:"support_id,string"`
	SupportForumId  uint64  `json:"support_forum_id,string"`
	TicketId        uint64  `json:"ticket_id,string"`
	Follower        string  `json:"follower" gorm:"type:varchar(255)"`
	SyncType        string  `json:"sync_type" gorm:"type:text"`
	Label           string  `json:"label" gorm:"type:varchar(255)"`
	Effort          float32 `json:"effort,string"`
	EffortCompleted float32 `json:"effort_completed,string"`
	Exceed          float32 `json:"exceed,string"`
	Remain          float32 `json:"remain,string"`
	Progress        string  `json:"progress" gorm:"type:varchar(255)"`
	Estimate        float32 `json:"estimate,string"`
	Bugtype         string  `json:"bugtype" gorm:"type:varchar(20)"`

	Milestone        string `json:"milestone" gorm:"type:varchar(255)"`
	CustomFieldOne   string `json:"custom_field_one" gorm:"type:text"`
	CustomFieldTwo   string `json:"custom_field_two" gorm:"type:text"`
	CustomFieldThree string `json:"custom_field_three" gorm:"type:text"`
	CustomFieldFour  string `json:"custom_field_four" gorm:"type:text"`
	CustomFieldFive  string `json:"custom_field_five" gorm:"type:text"`
	CustomField6     string `json:"custom_field_6" gorm:"type:text"`
	CustomField7     string `json:"custom_field_7" gorm:"type:text"`
	CustomField8     string `json:"custom_field_8" gorm:"type:text"`
	CustomField9     string `json:"custom_field_9" gorm:"type:text"`
	CustomField10    string `json:"custom_field_10" gorm:"type:text"`
	CustomField11    string `json:"custom_field_11" gorm:"type:text"`
	CustomField12    string `json:"custom_field_12" gorm:"type:text"`
	CustomField13    string `json:"custom_field_13" gorm:"type:text"`
	CustomField14    string `json:"custom_field_14" gorm:"type:text"`
	CustomField15    string `json:"custom_field_15" gorm:"type:text"`
	CustomField16    string `json:"custom_field_16" gorm:"type:text"`
	CustomField17    string `json:"custom_field_17" gorm:"type:text"`
	CustomField18    string `json:"custom_field_18" gorm:"type:text"`
	CustomField19    string `json:"custom_field_19" gorm:"type:text"`
	CustomField20    string `json:"custom_field_20" gorm:"type:text"`
	CustomField21    string `json:"custom_field_21" gorm:"type:text"`
	CustomField22    string `json:"custom_field_22" gorm:"type:text"`
	CustomField23    string `json:"custom_field_23" gorm:"type:text"`
	CustomField24    string `json:"custom_field_24" gorm:"type:text"`
	CustomField25    string `json:"custom_field_25" gorm:"type:text"`
	CustomField26    string `json:"custom_field_26" gorm:"type:text"`
	CustomField27    string `json:"custom_field_27" gorm:"type:text"`
	CustomField28    string `json:"custom_field_28" gorm:"type:text"`
	CustomField29    string `json:"custom_field_29" gorm:"type:text"`
	CustomField30    string `json:"custom_field_30" gorm:"type:text"`
	CustomField31    string `json:"custom_field_31" gorm:"type:text"`
	CustomField32    string `json:"custom_field_32" gorm:"type:text"`
	CustomField33    string `json:"custom_field_33" gorm:"type:text"`
	CustomField34    string `json:"custom_field_34" gorm:"type:text"`
	CustomField35    string `json:"custom_field_35" gorm:"type:text"`
	CustomField36    string `json:"custom_field_36" gorm:"type:text"`
	CustomField37    string `json:"custom_field_37" gorm:"type:text"`
	CustomField38    string `json:"custom_field_38" gorm:"type:text"`
	CustomField39    string `json:"custom_field_39" gorm:"type:text"`
	CustomField40    string `json:"custom_field_40" gorm:"type:text"`
	CustomField41    string `json:"custom_field_41" gorm:"type:text"`
	CustomField42    string `json:"custom_field_42" gorm:"type:text"`
	CustomField43    string `json:"custom_field_43" gorm:"type:text"`
	CustomField44    string `json:"custom_field_44" gorm:"type:text"`
	CustomField45    string `json:"custom_field_45" gorm:"type:text"`
	CustomField46    string `json:"custom_field_46" gorm:"type:text"`
	CustomField47    string `json:"custom_field_47" gorm:"type:text"`
	CustomField48    string `json:"custom_field_48" gorm:"type:text"`
	CustomField49    string `json:"custom_field_49" gorm:"type:text"`
	CustomField50    string `json:"custom_field_50" gorm:"type:text"`
}

func (TapdBug) TableName() string {
	return "_tool_tapd_bugs"
}
