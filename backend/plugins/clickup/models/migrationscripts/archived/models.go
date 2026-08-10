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

// Package archived holds frozen snapshots of the tool-layer models as they
// existed at each migration. The live models in plugins/clickup/models may
// evolve; these snapshots keep historical migrations stable.
package archived

import (
	"time"

	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
)

type ClickUpConnection struct {
	Name string `gorm:"type:varchar(100);uniqueIndex" json:"name"`
	archived.Model
	Endpoint         string `mapstructure:"endpoint" json:"endpoint"`
	Proxy            string `mapstructure:"proxy" json:"proxy"`
	RateLimitPerHour int    `json:"rateLimitPerHour"`
	Token            string `mapstructure:"token" json:"token" gorm:"serializer:encdec"`
}

func (ClickUpConnection) TableName() string { return "_tool_clickup_connections" }

type ClickUpList struct {
	archived.NoPKModel
	ConnectionId  uint64 `json:"connectionId" gorm:"primaryKey"`
	ScopeConfigId uint64 `json:"scopeConfigId,omitempty"`
	ListId        string `json:"listId" gorm:"primaryKey;type:varchar(255)"`
	Name          string `json:"name" gorm:"type:varchar(255)"`
	SpaceId       string `json:"spaceId" gorm:"type:varchar(255)"`
	SpaceName     string `json:"spaceName" gorm:"type:varchar(255)"`
}

func (ClickUpList) TableName() string { return "_tool_clickup_lists" }

type ClickUpScopeConfig struct {
	archived.ScopeConfig
	ConnectionId          uint64   `json:"connectionId" gorm:"index"`
	Name                  string   `gorm:"type:varchar(255);uniqueIndex" json:"name"`
	IssueStatusTodo       []string `json:"issueStatusTodo" gorm:"type:json;serializer:json"`
	IssueStatusInProgress []string `json:"issueStatusInProgress" gorm:"type:json;serializer:json"`
	IssueStatusDone       []string `json:"issueStatusDone" gorm:"type:json;serializer:json"`
	IssueTypeRequirement  string   `json:"issueTypeRequirement" gorm:"type:varchar(255)"`
	IssueTypeBug          string   `json:"issueTypeBug" gorm:"type:varchar(255)"`
	IssueTypeIncident     string   `json:"issueTypeIncident" gorm:"type:varchar(255)"`
}

func (ClickUpScopeConfig) TableName() string { return "_tool_clickup_scope_configs" }

type ClickUpUser struct {
	ConnectionId   uint64 `gorm:"primaryKey"`
	Id             string `gorm:"primaryKey;type:varchar(255)"`
	Username       string `gorm:"type:varchar(255)"`
	Email          string `gorm:"type:varchar(255)"`
	Color          string `gorm:"type:varchar(50)"`
	ProfilePicture string `gorm:"type:varchar(255)"`
	archived.NoPKModel
}

func (ClickUpUser) TableName() string { return "_tool_clickup_users" }

type ClickUpTask struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	Id           string `gorm:"primaryKey;type:varchar(255)"`
	ListId       string `gorm:"index;type:varchar(255)"`
	SpaceId      string `gorm:"type:varchar(255)"`
	CustomId     string `gorm:"type:varchar(255)"`
	Name         string
	Description  string
	Status       string `gorm:"type:varchar(255)"`
	StatusType   string `gorm:"type:varchar(100)"`
	Type         string `gorm:"type:varchar(100)"`
	Priority     string `gorm:"type:varchar(100)"`
	Url          string `gorm:"type:varchar(255)"`
	CreatorId    string `gorm:"type:varchar(255)"`
	AssigneeId   string `gorm:"type:varchar(255)"`
	AssigneeName string `gorm:"type:varchar(255)"`
	ParentId     string `gorm:"type:varchar(255)"`
	CreatedDate  *time.Time
	UpdatedDate  *time.Time `gorm:"index"`
	ClosedDate   *time.Time
	archived.NoPKModel
}

func (ClickUpTask) TableName() string { return "_tool_clickup_tasks" }

type ClickUpTaskComment struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	Id           string `gorm:"primaryKey;type:varchar(255)"`
	TaskId       string `gorm:"index;type:varchar(255)"`
	Body         string
	UserId       string `gorm:"type:varchar(255)"`
	CreatedDate  *time.Time
	archived.NoPKModel
}

func (ClickUpTaskComment) TableName() string { return "_tool_clickup_task_comments" }
