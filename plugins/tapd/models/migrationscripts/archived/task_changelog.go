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

type TapdTaskChangelog struct {
	ConnectionId   uint64         `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Id             uint64         `gorm:"primaryKey;type:BIGINT  NOT NULL" json:"id,string"`
	WorkspaceId    uint64         `json:"workspace_id,string"`
	WorkitemTypeId uint64         `json:"workitem_type_id,string"`
	Creator        string         `json:"creator" gorm:"type:varchar(255)"`
	Created        helper.CSTTime `json:"created"`
	ChangeSummary  string         `json:"change_summary" gorm:"type:varchar(255)"`
	Comment        string         `json:"comment"`
	EntityType     string         `json:"entity_type" gorm:"type:varchar(255)"`
	ChangeType     string         `json:"change_type" gorm:"type:varchar(255)"`
	ChangeTypeText string         `json:"change_type_text" gorm:"type:varchar(255)"`
	TaskId         uint64         `json:"task_id,string"`
	common.NoPKModel
	FieldChanges []TapdTaskChangelogItem `json:"field_changes" gorm:"-"`
}

type TapdTaskChangelogItem struct {
	ConnectionId      uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ChangelogId       uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Field             string `json:"field" gorm:"primaryKey;type:varchar(255)"`
	ValueBeforeParsed string `json:"value_before_parsed"`
	ValueAfterParsed  string `json:"value_after_parsed"`
	IterationIdFrom   uint64
	IterationIdTo     uint64
	common.NoPKModel
}

func (TapdTaskChangelog) TableName() string {
	return "_tool_tapd_task_changelogs"
}
func (TapdTaskChangelogItem) TableName() string {
	return "_tool_tapd_task_changelog_items"
}
