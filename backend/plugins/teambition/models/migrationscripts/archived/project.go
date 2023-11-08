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

type TeambitionProject struct {
	ConnectionId   uint64 `gorm:"primaryKey;type:BIGINT"`
	Id             string `gorm:"primaryKey;type:varchar(100)"`
	Name           string `gorm:"type:varchar(255)"`
	Logo           string `gorm:"type:varchar(255)"`
	Description    string `gorm:"type:text"`
	OrganizationId string `gorm:"type:varchar(255)"`
	Visibility     string `gorm:"typechar(100)"`
	IsTemplate     bool
	CreatorId      string `gorm:"type:varchar(100)"`
	IsArchived     bool
	IsSuspended    bool
	UniqueIdPrefix string `gorm:"type:varchar(255)"`
	Created        *time.Time
	Updated        *time.Time
	StartDate      *time.Time
	EndDate        *time.Time
	Customfields   []TeambitionProjectCustomField `gorm:"serializer:json;type:text"`

	archived.NoPKModel
}

type TeambitionProjectCustomField struct {
	CustomfieldId string                       `json:"customfieldId"`
	Type          string                       `json:"type"`
	Value         []TeambitionCustomFieldValue `json:"value"`
}

type TeambitionProjectCustomFieldValue struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	MetaString string `json:"metaString"`
}

func (TeambitionProject) TableName() string {
	return "_tool_teambition_projects"
}
