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
)

type JiraIssueType struct {
	ConnectionId     uint64 `gorm:"primaryKey;autoIncrement:false"`
	Self             string `json:"self" gorm:"type:varchar(255)"`
	Id               string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Description      string `json:"description"`
	IconURL          string `json:"iconUrl" gorm:"type:varchar(255)"`
	Name             string `json:"name" gorm:"type:varchar(255)"`
	UntranslatedName string `json:"untranslatedName" gorm:"type:varchar(255)"`
	Subtask          bool   `json:"subtask"`
	AvatarID         uint64 `json:"avatarId"`
	HierarchyLevel   int    `json:"hierarchyLevel"`
	common.NoPKModel
}

func (JiraIssueType) TableName() string {
	return "_tool_jira_issue_types"
}
