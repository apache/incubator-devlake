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

	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
)

type TeambitionTask struct {
	ConnectionId   uint64   `gorm:"primaryKey;type:BIGINT"`
	ProjectId      string   `gorm:"primaryKey;type:varchar(100)"`
	Id             string   `gorm:"primaryKey;type:varchar(100)"`
	Content        string   `gorm:"type:varchar(255)"`
	Note           string   `gorm:"type:varchar(255)"`
	AncestorIds    []string `gorm:"serializer:json;type:text"`
	ParentTaskId   string   `gorm:"type:varchar(100)"`
	TfsId          string   `gorm:"type:varchar(100)"`
	TasklistId     string   `gorm:"type:varchar(100)"`
	StageId        string   `gorm:"type:varchar(100)"`
	TagIds         []string `gorm:"serializer:json;type:text"`
	CreatorId      string   `gorm:"type:varchar(100)"`
	ExecutorId     string   `gorm:"type:varchar(100)"`
	InvolveMembers []string `gorm:"serializer:json;type:text"`
	Priority       int
	StoryPoint     string   `gorm:"varchar(255)"`
	Recurrence     []string `gorm:"serializer:json;type:text"`
	IsDone         bool
	IsArchived     bool
	Visible        string `gorm:"varchar(100)"`
	UniqueId       int64
	StartDate      *time.Time
	DueDate        *time.Time
	AccomplishTime *time.Time
	Created        *time.Time
	Updated        *time.Time
	SfcId          string                  `gorm:"type:varchar(100)"`
	SprintId       string                  `gorm:"type:varchar(100)"`
	Customfields   []TeambitionCustomField `gorm:"serializer:json;type:text"`

	StdType   string `gorm:"type:varchar(100)"`
	StdStatus string `gorm:"type:varchar(100)"`

	archived.NoPKModel
}

type TeambitionCustomField struct {
	CfId string `gorm:"varchar(100)"`
	Type string `gorm:"varchar(100)"`
}

type TeambitionCustomFieldValue struct {
	Id         string `gorm:"varchar(100)"`
	Title      string `gorm:"varchar(100)"`
	MetaString string `gorm:"varchar(100)"`
}

func (TeambitionTask) TableName() string {
	return "_tool_teambition_tasks"
}
