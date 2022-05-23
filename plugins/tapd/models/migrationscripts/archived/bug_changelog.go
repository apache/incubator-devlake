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

type TapdBugChangelog struct {
	ConnectionId uint64          `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	WorkspaceID  uint64          `gorm:"type:BIGINT  NOT NULL"`
	ID           uint64          `gorm:"primaryKey;type:BIGINT  NOT NULL" json:"id,string"`
	BugID        uint64          `json:"bug_id,string"`
	Author       string          `json:"author"`
	Field        string          `gorm:"primaryKey;type:varchar(255)" json:"field"`
	OldValue     string          `json:"old_value"`
	NewValue     string          `json:"new_value"`
	Memo         string          `json:"memo"`
	Created      *helper.CSTTime `json:"created"`
	common.NoPKModel
}

type TapdBugChangelogItem struct {
	ConnectionId      uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	ChangelogId       uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	Field             string `json:"field" gorm:"primaryKey;foreignKey:ChangelogId;references:ID"`
	ValueBeforeParsed string `json:"value_before_parsed"`
	ValueAfterParsed  string `json:"value_after_parsed"`
	IterationIdFrom   uint64
	IterationIdTo     uint64
	common.NoPKModel
}

func (TapdBugChangelog) TableName() string {
	return "_tool_tapd_bug_changelogs"
}
func (TapdBugChangelogItem) TableName() string {
	return "_tool_tapd_bug_changelog_items"
}
