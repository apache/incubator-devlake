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

type CircleciProject struct {
	common.Scope   `mapstructure:",squash"`
	Id             string `gorm:"primaryKey;type:varchar(100)" json:"id" mapstructure:"id"`
	Slug           string `gorm:"type:varchar(255)" json:"slug" mapstructure:"slug"`
	Name           string `gorm:"type:varchar(255)" json:"name" mapstructure:"name"`
	OrganizationId string `gorm:"type:varchar(100)" json:"organizationId" mapstructure:"organizationId"`
	// VcsInfo        CircleciVcsInfo `gorm:"serializer:json;type:text" json:"vcsInfo" mapstructure:"vcsInfo"`
}

type CircleciVcsInfo struct {
	VcsUrl        string `json:"vcsUrl"`
	Provider      string `json:"provider"`
	DefaultBranch string `json:"defaultBranch"`
}

func (CircleciProject) TableName() string {
	return "_tool_circleci_projects"
}

var _ plugin.ToolLayerScope = (*CircleciProject)(nil)

type CircleciApiParams struct {
	ConnectionId uint64
	ProjectSlug  string
}

// ScopeFullName implements plugin.ToolLayerScope.
func (c CircleciProject) ScopeFullName() string {
	return c.Slug
}

// ScopeId implements plugin.ToolLayerScope.
func (c CircleciProject) ScopeId() string {
	return c.Id
}

// ScopeName implements plugin.ToolLayerScope.
func (c CircleciProject) ScopeName() string {
	return c.Name
}

// ScopeParams implements plugin.ToolLayerScope.
func (c CircleciProject) ScopeParams() interface{} {
	return &CircleciApiParams{
		ConnectionId: c.ConnectionId,
		ProjectSlug:  c.Slug,
	}
}
