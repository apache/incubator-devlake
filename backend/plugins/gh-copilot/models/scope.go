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
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
)

// GhCopilotScope represents an organization-level collection scope.
type GhCopilotScope struct {
	common.Scope       `mapstructure:",squash"`
	Id                 string     `json:"id" mapstructure:"id" gorm:"primaryKey;type:varchar(255)"`
	Organization       string     `json:"organization" mapstructure:"organization" gorm:"type:varchar(255);not null"`
	Name               string     `json:"name" mapstructure:"name" gorm:"type:varchar(255)"`
	FullName           string     `json:"fullName" mapstructure:"fullName" gorm:"type:varchar(255)"`
	ImplementationDate *time.Time `json:"implementationDate" mapstructure:"implementationDate"`
	BaselinePeriodDays int        `json:"baselinePeriodDays" mapstructure:"baselinePeriodDays"`
	SeatsLastSyncedAt  *time.Time `json:"seatsLastSyncedAt" mapstructure:"seatsLastSyncedAt"`
}

func (GhCopilotScope) TableName() string {
	return "_tool_copilot_scopes"
}

func (s *GhCopilotScope) BeforeSave() error {
	// Populate Name and FullName from Organization and Id
	if s.Name == "" {
		s.Name = s.ScopeName()
	}
	if s.FullName == "" {
		s.FullName = s.ScopeFullName()
	}
	// Validate and normalize BaselinePeriodDays (7-365 range, default 90)
	if s.BaselinePeriodDays < 7 {
		s.BaselinePeriodDays = 90 // Default to 90 days
	} else if s.BaselinePeriodDays > 365 {
		s.BaselinePeriodDays = 365 // Cap at 1 year
	}
	return nil
}

func (s GhCopilotScope) ScopeId() string {
	return s.Id
}

func (s GhCopilotScope) ScopeName() string {
	if s.Id != "" {
		return s.Id
	}
	return s.Organization
}

func (s GhCopilotScope) ScopeFullName() string {
	return s.ScopeName()
}

func (s GhCopilotScope) ScopeParams() interface{} {
	return &GhCopilotScopeParams{
		ConnectionId: s.ConnectionId,
		ScopeId:      s.Id,
	}
}

// GhCopilotScopeParams is returned for blueprint configuration.
type GhCopilotScopeParams struct {
	ConnectionId uint64 `json:"connectionId"`
	ScopeId      string `json:"scopeId"`
}

var _ plugin.ToolLayerScope = (*GhCopilotScope)(nil)
