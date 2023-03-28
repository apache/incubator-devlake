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
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type TeambitionTask struct {
	ConnectionId   uint64                  `gorm:"primaryKey;type:BIGINT"`
	ProjectId      string                  `gorm:"primaryKey;type:varchar(100)" json:"projectId"`
	Id             string                  `gorm:"primaryKey;type:varchar(100)" json:"id"`
	Content        string                  `gorm:"type:varchar(255)" json:"content"`
	Note           string                  `gorm:"type:varchar(255)" json:"Content"`
	AncestorIds    []string                `gorm:"serializer:json;type:text" json:"ancestorIds"`
	ParentTaskId   string                  `gorm:"type:varchar(100)" json:"parentTaskId"`
	TfsId          string                  `gorm:"type:varchar(100)" json:"tfsId"`
	TasklistId     string                  `gorm:"type:varchar(100)" json:"tasklistId"`
	StageId        string                  `gorm:"type:varchar(100)" json:"stageId"`
	TagIds         []string                `gorm:"serializer:json;type:text" json:"tagIds"`
	CreatorId      string                  `gorm:"type:varchar(100)" json:"creatorId"`
	ExecutorId     string                  `gorm:"type:varchar(100)" json:"executorId"`
	InvolveMembers []string                `gorm:"serializer:json;type:text" json:"involveMembers"`
	Priority       int                     `json:"priority"`
	StoryPoint     string                  `gorm:"varchar(255)" json:"storyPoint"`
	Recurrence     []string                `gorm:"serializer:json;type:text" json:"recurrence"`
	IsDone         bool                    `json:"isDone"`
	IsArchived     bool                    `json:"isArchived"`
	Visible        string                  `gorm:"varchar(100)" json:"visible"`
	UniqueId       int64                   `json:"uniqueId"`
	StartDate      *api.Iso8601Time        `json:"startDate"`
	DueDate        *api.Iso8601Time        `json:"dueDate"`
	AccomplishTime *api.Iso8601Time        `json:"accomplishTime"`
	Created        *api.Iso8601Time        `json:"created"`
	Updated        *api.Iso8601Time        `json:"updated"`
	SfcId          string                  `gorm:"type:varchar(100)" json:"sfcId"`
	SprintId       string                  `gorm:"type:varchar(100)" json:"sprintId"`
	Customfields   []TeambitionCustomField `gorm:"serializer:json;type:text" json:"customfields"`

	StdType   string `gorm:"type:varchar(100)" json:"stdType"`
	StdStatus string `gorm:"type:varchar(100)" json:"stdStatus"`

	common.NoPKModel
}

type TeambitionCustomField struct {
	CfId string `gorm:"varchar(100)" json:"cfId"`
	Type string `gorm:"varchar(100)" json:"type"`
}

type TeambitionCustomFieldValue struct {
	Id         string `gorm:"varchar(100)" json:"id"`
	Title      string `gorm:"varchar(100)" json:"title"`
	MetaString string `gorm:"varchar(100)" json:"metaString"`
}

func (TeambitionTask) TableName() string {
	return "_tool_teambition_tasks"
}
