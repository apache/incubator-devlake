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
	"time"

	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
)

type GitlabIssue struct {
	ConnectionId    uint64 `gorm:"primaryKey"`
	GitlabId        int    `gorm:"primaryKey"`
	ProjectId       int    `gorm:"index"`
	Number          int    `gorm:"index;comment:Used in API requests ex. api/repo/1/issue/<THIS_NUMBER>"`
	State           string `gorm:"type:varchar(255)"`
	Title           string
	Body            string
	Priority        string `gorm:"type:varchar(255)"`
	Type            string `gorm:"type:varchar(100)"`
	Status          string `gorm:"type:varchar(255)"`
	CreatorId       int
	CreatorName     string `gorm:"type:varchar(255)"`
	AssigneeId      int
	AssigneeName    string `gorm:"type:varchar(255)"`
	LeadTimeMinutes uint
	Url             string `gorm:"type:varchar(255)"`
	ClosedAt        *time.Time
	GitlabCreatedAt time.Time
	GitlabUpdatedAt time.Time `gorm:"index"`
	Severity        string    `gorm:"type:varchar(255)"`
	Component       string    `gorm:"type:varchar(255)"`
	TimeEstimate    int64
	TotalTimeSpent  int64
	archived.NoPKModel
}

func (GitlabIssue) TableName() string {
	return "_tool_gitlab_issues"
}

type GitlabAuthor struct {
	ConnectionId    uint64 `gorm:"primaryKey"`
	GitlabId        int    `gorm:"primaryKey" json:"id"`
	Username        string `gorm:"type:varchar(255)"`
	Email           string `gorm:"type:varchar(255)"`
	Name            string `gorm:"type:varchar(255)"`
	State           string `gorm:"type:varchar(255)"`
	MembershipState string `json:"membership_state" gorm:"type:varchar(255)"`
	AvatarUrl       string `json:"avatar_url" gorm:"type:varchar(255)"`
	WebUrl          string `json:"web_url" gorm:"type:varchar(255)"`

	archived.NoPKModel
}

func (GitlabAuthor) TableName() string {
	return "_tool_gitlab_accounts"
}

type GitlabAssignee struct {
	ConnectionId    uint64 `gorm:"primaryKey"`
	GitlabId        int    `gorm:"primaryKey" json:"id"`
	Username        string `gorm:"type:varchar(255)"`
	Email           string `gorm:"type:varchar(255)"`
	Name            string `gorm:"type:varchar(255)"`
	State           string `gorm:"type:varchar(255)"`
	MembershipState string `json:"membership_state" gorm:"type:varchar(255)"`
	AvatarUrl       string `json:"avatar_url" gorm:"type:varchar(255)"`
	WebUrl          string `json:"web_url" gorm:"type:varchar(255)"`

	archived.NoPKModel
}

func (GitlabAssignee) TableName() string {
	return "_tool_gitlab_accounts"
}
