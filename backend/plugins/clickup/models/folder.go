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

var _ plugin.ToolLayerScope = (*ClickUpFolder)(nil)

// ClickUpFolder is the data-source scope for the ClickUp plugin. A ClickUp
// folder (e.g. a team's "Dev Team" / "Sprint Folder") owns the backlog and the
// rolling sprint lists, so it maps cleanly onto a DevLake domain-layer
// ticket.Board — analogous to a Jira board. Selecting the folder (rather than
// an individual list) means new sprints are picked up automatically on each
// sync and ephemeral/archived sprint lists never need re-scoping.
type ClickUpFolder struct {
	common.Scope `mapstructure:",squash"`
	FolderId     string `json:"folderId" mapstructure:"folderId" gorm:"primaryKey;type:varchar(255)"`
	Name         string `json:"name" mapstructure:"name" gorm:"type:varchar(255)"`
	SpaceId      string `json:"spaceId" mapstructure:"spaceId" gorm:"type:varchar(255)"`
	SpaceName    string `json:"spaceName" mapstructure:"spaceName" gorm:"type:varchar(255)"`
}

func (f ClickUpFolder) ScopeId() string {
	return f.FolderId
}

func (f ClickUpFolder) ScopeName() string {
	return f.Name
}

func (f ClickUpFolder) ScopeFullName() string {
	if f.SpaceName != "" {
		return f.SpaceName + "/" + f.Name
	}
	return f.Name
}

func (f ClickUpFolder) ScopeParams() interface{} {
	return &ClickUpApiParams{
		ConnectionId: f.ConnectionId,
		FolderId:     f.FolderId,
	}
}

func (ClickUpFolder) TableName() string {
	return "_tool_clickup_folders"
}

// ClickUpApiParams identifies the scope a raw row belongs to. It is stored in
// the `params` column of every _raw_clickup_* table.
type ClickUpApiParams struct {
	ConnectionId uint64
	FolderId     string
}
