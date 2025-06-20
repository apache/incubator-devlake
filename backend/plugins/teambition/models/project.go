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
	"github.com/apache/incubator-devlake/core/plugin"
)

func (t TeambitionProject) ConvertApiScope() plugin.ToolLayerScope {
	return t
}

type TeambitionProject struct {
	common.Scope   `mapstructure:",squash"`
	Id             string                         `gorm:"primaryKey;type:varchar(100)" json:"id"`
	Name           string                         `gorm:"type:varchar(255)" json:"name"`
	Logo           string                         `gorm:"type:varchar(255)" json:"logo"`
	Description    string                         `gorm:"type:text" json:"description"`
	OrganizationId string                         `gorm:"type:varchar(255)" json:"organizationId"`
	Visibility     string                         `gorm:"typechar(100)" json:"visibility"`
	IsTemplate     bool                           `json:"isTemplate"`
	CreatorId      string                         `gorm:"type:varchar(100)" json:"creatorId"`
	IsArchived     bool                           `json:"isArchived"`
	IsSuspended    bool                           `json:"isSuspended"`
	UniqueIdPrefix string                         `gorm:"type:varchar(255)" json:"uniqueIdPrefix"`
	Created        *common.Iso8601Time            `json:"created"`
	Updated        *common.Iso8601Time            `json:"updated"`
	StartDate      *common.Iso8601Time            `json:"startDate"`
	EndDate        *common.Iso8601Time            `json:"endDate"`
	Customfields   []TeambitionProjectCustomField `gorm:"serializer:json;type:text" json:"customfields"`
}

// ScopeFullName implements plugin.ToolLayerScope.
func (t TeambitionProject) ScopeFullName() string {
	return t.Name
}

// ScopeId implements plugin.ToolLayerScope.
func (t TeambitionProject) ScopeId() string {
	return t.Id
}

// ScopeName implements plugin.ToolLayerScope.
func (t TeambitionProject) ScopeName() string {
	return t.Name
}

// ScopeParams implements plugin.ToolLayerScope.
func (t TeambitionProject) ScopeParams() interface{} {
	return &TeambitionApiParams{
		ConnectionId: t.ConnectionId,
		ProjectId:    t.Id,
	}
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

type TeambitionApiParams struct {
	ConnectionId   uint64
	OrganizationId string
	ProjectId      string
}

type TeambitionProjectsResponse struct {
	Status int `json:"status"`
	Data   []struct {
		TeambitionProject `json:"Workspace"`
	} `json:"data"`
	Info string `json:"info"`
}
