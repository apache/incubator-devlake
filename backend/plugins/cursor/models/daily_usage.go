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

	"github.com/apache/devlake/core/models/common"
)

// CursorDailyUsage stores per-user per-day adoption metrics from POST /teams/daily-usage-data.
// Line fields map from acceptedLinesAdded/totalLinesAdded (etc.); LinesAdded/LinesDeleted
// duplicate accepted_lines_added/accepted_lines_deleted for backward compatibility.
type CursorDailyUsage struct {
	ConnectionId uint64    `gorm:"primaryKey" json:"connectionId"`
	ScopeId      string    `gorm:"primaryKey;type:varchar(255)" json:"scopeId"`
	UserId       string    `gorm:"primaryKey;type:varchar(255)" json:"userId"`
	UsageDate    time.Time `gorm:"primaryKey" json:"usageDate"`

	Email                    string `json:"email" gorm:"type:varchar(255);index"`
	IsActive                 bool   `json:"isActive"`
	Completions              int    `json:"completions"`
	PremiumRequests          int    `json:"premiumRequests"`
	AgentRequests            int    `json:"agentRequests"`
	ChatRequests             int    `json:"chatRequests"`
	ComposerRequests         int    `json:"composerRequests"`
	TabsAccepted             int    `json:"tabsAccepted"`
	TabsShown                int    `json:"tabsShown"`
	UsageBasedReqs           int    `json:"usageBasedReqs"`
	SubscriptionIncludedReqs int    `json:"subscriptionIncludedReqs"`
	MostUsedModel            string `json:"mostUsedModel" gorm:"type:varchar(255)"`
	ClientVersion            string `json:"clientVersion" gorm:"type:varchar(100)"`
	AcceptedLinesAdded       int    `json:"acceptedLinesAdded"`
	AcceptedLinesDeleted     int    `json:"acceptedLinesDeleted"`
	TotalLinesAdded          int    `json:"totalLinesAdded"`
	TotalLinesDeleted        int    `json:"totalLinesDeleted"`
	TotalApplies             int    `json:"totalApplies"`
	TotalAccepts             int    `json:"totalAccepts"`
	TotalRejects             int    `json:"totalRejects"`
	LinesAdded               int    `json:"linesAdded"`
	LinesDeleted             int    `json:"linesDeleted"`

	common.NoPKModel
}

func (CursorDailyUsage) TableName() string {
	return "_tool_cursor_daily_usage"
}
