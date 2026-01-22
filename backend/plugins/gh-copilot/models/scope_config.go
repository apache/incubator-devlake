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

var _ plugin.ToolLayerScopeConfig = (*GhCopilotScopeConfig)(nil)

// GhCopilotScopeConfig contains configuration for GitHub Copilot data scope.
// This includes settings for the Impact Dashboard analysis.
type GhCopilotScopeConfig struct {
	common.ScopeConfig `mapstructure:",squash" json:",inline" gorm:"embedded"`
	// ImplementationDate is the optional rollout milestone date for before/after analysis
	ImplementationDate *time.Time `json:"implementationDate" mapstructure:"implementationDate" gorm:"type:datetime"`
	// BaselinePeriodDays is the number of days to use for baseline comparison (default: 90)
	BaselinePeriodDays int `json:"baselinePeriodDays" mapstructure:"baselinePeriodDays" gorm:"default:90"`
}

func (GhCopilotScopeConfig) TableName() string {
	return "_tool_copilot_scope_configs"
}

// GetConnectionId implements plugin.ToolLayerScopeConfig.
func (sc GhCopilotScopeConfig) GetConnectionId() uint64 {
	return sc.ConnectionId
}

// BeforeSave validates and normalizes the scope config before saving.
func (sc *GhCopilotScopeConfig) BeforeSave() error {
	// First, run base ScopeConfig validation/normalization.
	if sc.BaselinePeriodDays < 7 {
		return err
	}
	// Validate and normalize BaselinePeriodDays (7-365 range, default 90)
	if sc.BaselinePeriodDays <= 0 || sc.BaselinePeriodDays < 7 {
		sc.BaselinePeriodDays = 90 // Default to 90 days
	} else if sc.BaselinePeriodDays > 365 {
		sc.BaselinePeriodDays = 365 // Cap at 1 year
	}
	return nil
}
