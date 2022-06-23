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

	"github.com/apache/incubator-devlake/models/common"
)

type GitlabMergeRequest struct {
	ConnectionId     uint64 `gorm:"primaryKey"`
	GitlabId         int    `gorm:"primaryKey"`
	Iid              int    `gorm:"index"`
	ProjectId        int    `gorm:"index"`
	SourceProjectId  int
	TargetProjectId  int
	State            string `gorm:"type:varchar(255)"`
	Title            string
	WebUrl           string `gorm:"type:varchar(255)"`
	UserNotesCount   int
	WorkInProgress   bool
	SourceBranch     string `gorm:"type:varchar(255)"`
	TargetBranch     string `gorm:"type:varchar(255)"`
	MergeCommitSha   string `gorm:"type:varchar(255)"`
	MergedAt         *time.Time
	GitlabCreatedAt  time.Time
	ClosedAt         *time.Time
	Type             string `gorm:"type:varchar(255)"`
	MergedByUsername string `gorm:"type:varchar(255)"`
	Description      string
	AuthorUsername   string `gorm:"type:varchar(255)"`
	AuthorUserId     int
	Component        string     `gorm:"type:varchar(255)"`
	FirstCommentTime *time.Time `gorm:"comment:Time when the first comment occurred"`
	ReviewRounds     int        `gorm:"comment:How many rounds of review this MR went through"`
	common.NoPKModel
}

func (GitlabMergeRequest) TableName() string {
	return "_tool_gitlab_merge_requests"
}
