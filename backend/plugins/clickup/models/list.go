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

var _ plugin.ToolLayerScope = (*ClickUpList)(nil)

// ClickUpList is the data-source scope for the ClickUp plugin. A ClickUp List
// owns tasks (analogous to a Jira board), mapping cleanly to a DevLake
// domain-layer ticket.Board.
type ClickUpList struct {
	common.Scope `mapstructure:",squash"`
	ListId       string `json:"listId" mapstructure:"listId" gorm:"primaryKey;type:varchar(255)"`
	Name         string `json:"name" mapstructure:"name" gorm:"type:varchar(255)"`
	SpaceId      string `json:"spaceId" mapstructure:"spaceId" gorm:"type:varchar(255)"`
	SpaceName    string `json:"spaceName" mapstructure:"spaceName" gorm:"type:varchar(255)"`
}

func (l ClickUpList) ScopeId() string {
	return l.ListId
}

func (l ClickUpList) ScopeName() string {
	return l.Name
}

func (l ClickUpList) ScopeFullName() string {
	if l.SpaceName != "" {
		return l.SpaceName + "/" + l.Name
	}
	return l.Name
}

func (l ClickUpList) ScopeParams() interface{} {
	return &ClickUpApiParams{
		ConnectionId: l.ConnectionId,
		ListId:       l.ListId,
	}
}

func (ClickUpList) TableName() string {
	return "_tool_clickup_lists"
}

// ClickUpApiParams identifies the scope a raw row belongs to. It is stored in
// the `params` column of every _raw_clickup_* table.
type ClickUpApiParams struct {
	ConnectionId uint64
	ListId       string
}
