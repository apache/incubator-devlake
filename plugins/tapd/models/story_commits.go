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

type TapdStoryCommit struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	Id           uint64 `gorm:"primaryKey;type:BIGINT" json:"id,string"`

	UserId          string `json:"user_id" gorm:"type:varchar(255)"`
	HookUserName    string `json:"hook_user_name" gorm:"type:varchar(255)"`
	CommitId        string `json:"commit_id" gorm:"type:varchar(255)"`
	WorkspaceId     uint64 `json:"workspace_id,string" gorm:"type:BIGINT"`
	Message         string `json:"message" gorm:"type:text"`
	Path            string `json:"path" gorm:"type:varchar(255)"`
	WebURL          string `json:"web_url" gorm:"type:varchar(255)"`
	HookProjectName string `json:"hook_project_name" gorm:"type:varchar(255)"`

	Ref        string         `json:"ref" gorm:"type:varchar(255)"`
	RefStatus  string         `json:"ref_status" gorm:"type:varchar(255)"`
	GitEnv     string         `json:"git_env" gorm:"type:varchar(255)"`
	FileCommit string         `json:"file_commit"`
	CommitTime helper.CSTTime `json:"commit_time"`
	Created    helper.CSTTime `json:"created"`

	StoryId uint64
	common.NoPKModel
}

func (TapdStoryCommit) TableName() string {
	return "_tool_tapd_story_commits"
}
