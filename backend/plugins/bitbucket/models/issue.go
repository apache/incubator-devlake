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

	"github.com/apache/incubator-devlake/core/models/common"
)

type BitbucketIssue struct {
	ConnectionId       uint64 `gorm:"primaryKey"`
	RepoId             string `gorm:"primaryKey;type:varchar(255)"`
	BitbucketId        int    `gorm:"primaryKey"`
	Number             int    `gorm:"index;comment:Used in API requests ex. api/issues/<THIS_NUMBER>"`
	State              string `gorm:"type:varchar(255)"`
	StdState           string `gorm:"type:varchar(255)"`
	Title              string `gorm:"type:varchar(255)"`
	Body               string
	Priority           string `gorm:"type:varchar(255)"`
	Type               string `gorm:"type:varchar(100)"`
	AuthorId           string `gorm:"type:varchar(255)"`
	AuthorName         string `gorm:"type:varchar(255)"`
	AssigneeId         string `gorm:"type:varchar(255)"`
	AssigneeName       string `gorm:"type:varchar(255)"`
	MilestoneId        int    `gorm:"index"`
	LeadTimeMinutes    *uint
	Url                string `gorm:"type:varchar(255)"`
	ClosedAt           *time.Time
	BitbucketCreatedAt time.Time
	BitbucketUpdatedAt time.Time `gorm:"index"`
	Severity           string    `gorm:"type:varchar(255)"`
	Component          string    `gorm:"type:text"`
	common.NoPKModel
}

func (BitbucketIssue) TableName() string {
	return "_tool_bitbucket_issues"
}
