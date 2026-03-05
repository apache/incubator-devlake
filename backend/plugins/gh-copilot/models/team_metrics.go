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
)

// GhCopilotTeamDailyMetrics stores daily aggregate team-level Copilot metrics.
type GhCopilotTeamDailyMetrics struct {
	ConnectionId uint64    `gorm:"primaryKey" json:"connectionId"`
	ScopeId      string    `gorm:"primaryKey;type:varchar(255)" json:"scopeId"`
	TeamSlug     string    `gorm:"primaryKey;type:varchar(191);index" json:"teamSlug"`
	Date         time.Time `gorm:"primaryKey;type:date" json:"date"`

	TotalActiveUsers             int `json:"totalActiveUsers"`
	TotalEngagedUsers            int `json:"totalEngagedUsers"`
	CompletionsTotalEngagedUsers int `json:"completionsTotalEngagedUsers"`
	IdeChatTotalEngagedUsers     int `json:"ideChatTotalEngagedUsers"`
	DotcomChatTotalEngagedUsers  int `json:"dotcomChatTotalEngagedUsers"`
	DotcomPrTotalEngagedUsers    int `json:"dotcomPrTotalEngagedUsers"`

	common.NoPKModel
}

func (GhCopilotTeamDailyMetrics) TableName() string {
	return "_tool_copilot_team_daily_metrics"
}

// GhCopilotTeamCompletions stores team-level IDE code completion metrics by editor/model/language.
type GhCopilotTeamCompletions struct {
	ConnectionId uint64    `gorm:"primaryKey" json:"connectionId"`
	ScopeId      string    `gorm:"primaryKey;type:varchar(255)" json:"scopeId"`
	TeamSlug     string    `gorm:"primaryKey;type:varchar(191);index" json:"teamSlug"`
	Date         time.Time `gorm:"primaryKey;type:date" json:"date"`
	Editor       string    `gorm:"primaryKey;type:varchar(50)" json:"editor"`
	Model        string    `gorm:"primaryKey;type:varchar(100)" json:"model"`
	Language     string    `gorm:"primaryKey;type:varchar(50)" json:"language"`

	TotalEngagedUsers       int        `json:"totalEngagedUsers"`
	TotalCodeSuggestions    int        `json:"totalCodeSuggestions"`
	TotalCodeAcceptances    int        `json:"totalCodeAcceptances"`
	TotalCodeLinesSuggested int        `json:"totalCodeLinesSuggested"`
	TotalCodeLinesAccepted  int        `json:"totalCodeLinesAccepted"`
	IsCustomModel           bool       `json:"isCustomModel"`
	CustomModelTrainingDate *time.Time `gorm:"type:date" json:"customModelTrainingDate"`

	common.NoPKModel
}

func (GhCopilotTeamCompletions) TableName() string {
	return "_tool_copilot_team_completions"
}

// GhCopilotTeamIdeChat stores team-level IDE chat metrics by editor/model.
type GhCopilotTeamIdeChat struct {
	ConnectionId uint64    `gorm:"primaryKey" json:"connectionId"`
	ScopeId      string    `gorm:"primaryKey;type:varchar(255)" json:"scopeId"`
	TeamSlug     string    `gorm:"primaryKey;type:varchar(191);index" json:"teamSlug"`
	Date         time.Time `gorm:"primaryKey;type:date" json:"date"`
	Editor       string    `gorm:"primaryKey;type:varchar(50)" json:"editor"`
	Model        string    `gorm:"primaryKey;type:varchar(100)" json:"model"`

	TotalEngagedUsers        int        `json:"totalEngagedUsers"`
	TotalChats               int        `json:"totalChats"`
	TotalChatInsertionEvents int        `json:"totalChatInsertionEvents"`
	TotalChatCopyEvents      int        `json:"totalChatCopyEvents"`
	IsCustomModel            bool       `json:"isCustomModel"`
	CustomModelTrainingDate  *time.Time `gorm:"type:date" json:"customModelTrainingDate"`

	common.NoPKModel
}

func (GhCopilotTeamIdeChat) TableName() string {
	return "_tool_copilot_team_ide_chat"
}

// GhCopilotTeamDotcomChat stores team-level dotcom chat metrics by model.
type GhCopilotTeamDotcomChat struct {
	ConnectionId uint64    `gorm:"primaryKey" json:"connectionId"`
	ScopeId      string    `gorm:"primaryKey;type:varchar(255)" json:"scopeId"`
	TeamSlug     string    `gorm:"primaryKey;type:varchar(191);index" json:"teamSlug"`
	Date         time.Time `gorm:"primaryKey;type:date" json:"date"`
	Model        string    `gorm:"primaryKey;type:varchar(100)" json:"model"`

	TotalEngagedUsers       int        `json:"totalEngagedUsers"`
	TotalChats              int        `json:"totalChats"`
	IsCustomModel           bool       `json:"isCustomModel"`
	CustomModelTrainingDate *time.Time `gorm:"type:date" json:"customModelTrainingDate"`

	common.NoPKModel
}

func (GhCopilotTeamDotcomChat) TableName() string {
	return "_tool_copilot_team_dotcom_chat"
}

// GhCopilotTeamDotcomPrs stores team-level dotcom pull request metrics by repository/model.
type GhCopilotTeamDotcomPrs struct {
	ConnectionId uint64    `gorm:"primaryKey" json:"connectionId"`
	ScopeId      string    `gorm:"primaryKey;type:varchar(255)" json:"scopeId"`
	TeamSlug     string    `gorm:"primaryKey;type:varchar(191);index" json:"teamSlug"`
	Date         time.Time `gorm:"primaryKey;type:date" json:"date"`
	Repository   string    `gorm:"primaryKey;type:varchar(191)" json:"repository"`
	Model        string    `gorm:"primaryKey;type:varchar(100)" json:"model"`

	TotalEngagedUsers       int        `json:"totalEngagedUsers"`
	TotalPrSummariesCreated int        `json:"totalPrSummariesCreated"`
	IsCustomModel           bool       `json:"isCustomModel"`
	CustomModelTrainingDate *time.Time `gorm:"type:date" json:"customModelTrainingDate"`

	common.NoPKModel
}

func (GhCopilotTeamDotcomPrs) TableName() string {
	return "_tool_copilot_team_dotcom_prs"
}
