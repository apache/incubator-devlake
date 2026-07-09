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
	"strings"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"gorm.io/gorm"
)

const (
	// DefaultScopeID is the single team-level scope for Cursor Admin API data.
	DefaultScopeID = "team"
)

// CursorScope represents a team-level collection scope.
type CursorScope struct {
	common.Scope `mapstructure:",squash"`
	Id           string `json:"id" mapstructure:"id" gorm:"primaryKey;type:varchar(255)"`
	TeamId       string `json:"teamId" mapstructure:"teamId" gorm:"type:varchar(255)"`
	Name         string `json:"name" mapstructure:"name" gorm:"type:varchar(255)"`
	FullName     string `json:"fullName" mapstructure:"fullName" gorm:"type:varchar(255)"`
}

func (CursorScope) TableName() string {
	return "_tool_cursor_scopes"
}

func (s *CursorScope) BeforeSave(tx *gorm.DB) error {
	if s == nil {
		return nil
	}

	s.Id = strings.TrimSpace(s.Id)
	s.TeamId = strings.TrimSpace(s.TeamId)
	s.Name = strings.TrimSpace(s.Name)
	s.FullName = strings.TrimSpace(s.FullName)

	if s.Id == "" {
		s.Id = DefaultScopeID
	}
	if s.Name == "" {
		s.Name = s.ScopeName()
	}
	if s.FullName == "" {
		s.FullName = s.ScopeFullName()
	}

	return nil
}

func (s CursorScope) ScopeId() string {
	return s.Id
}

func (s CursorScope) ScopeName() string {
	if s.Name != "" {
		return s.Name
	}
	return DefaultScopeID
}

func (s CursorScope) ScopeFullName() string {
	if s.FullName != "" {
		return s.FullName
	}
	return s.ScopeName()
}

func (s CursorScope) ScopeParams() interface{} {
	return &CursorScopeParams{
		ConnectionId: s.ConnectionId,
		ScopeId:      s.Id,
	}
}

// CursorScopeParams is returned for blueprint configuration.
type CursorScopeParams struct {
	ConnectionId uint64 `json:"connectionId"`
	ScopeId      string `json:"scopeId"`
}

var _ plugin.ToolLayerScope = (*CursorScope)(nil)
