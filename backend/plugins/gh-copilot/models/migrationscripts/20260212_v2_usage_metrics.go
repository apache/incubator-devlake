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

package migrationscripts

import (
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type migrateToUsageMetricsV2 struct{}

// --- Snapshot structs for migration (avoid importing models package to prevent drift) ---

// Connection: add Enterprise column
type connection20260212 struct {
	Enterprise string `gorm:"type:varchar(100)"`
}

func (connection20260212) TableName() string {
	return "_tool_copilot_connections"
}

// Scope: add Enterprise column
type scope20260212 struct {
	Enterprise string `gorm:"type:varchar(100)"`
}

func (scope20260212) TableName() string {
	return "_tool_copilot_scopes"
}

// --- Common embedded structs ---

type activityMetrics20260212 struct {
	UserInitiatedInteractionCount int `gorm:"comment:Chat messages and inline prompts initiated by user"`
	CodeGenerationActivityCount   int
	CodeAcceptanceActivityCount   int
	LocSuggestedToAddSum          int
	LocSuggestedToDeleteSum       int
	LocAddedSum                   int
	LocDeletedSum                 int
}

type codeMetrics20260212 struct {
	CodeGenerationActivityCount int
	CodeAcceptanceActivityCount int
	LocSuggestedToAddSum        int
	LocSuggestedToDeleteSum     int
	LocAddedSum                 int
	LocDeletedSum               int
}

// --- Enterprise metrics tables ---

type enterpriseDailyMetrics20260212 struct {
	ConnectionId            uint64    `gorm:"primaryKey"`
	ScopeId                 string    `gorm:"primaryKey;type:varchar(255)"`
	Day                     time.Time `gorm:"primaryKey;type:date"`
	EnterpriseId            string    `gorm:"type:varchar(100)"`
	DailyActiveUsers        int
	WeeklyActiveUsers       int
	MonthlyActiveUsers      int
	MonthlyActiveChatUsers  int
	MonthlyActiveAgentUsers int
	activityMetrics20260212 `gorm:"embedded"`
	common.NoPKModel
}

func (enterpriseDailyMetrics20260212) TableName() string {
	return "_tool_copilot_enterprise_daily_metrics"
}

type metricsByIde20260212 struct {
	ConnectionId            uint64    `gorm:"primaryKey"`
	ScopeId                 string    `gorm:"primaryKey;type:varchar(255)"`
	Day                     time.Time `gorm:"primaryKey;type:date"`
	Ide                     string    `gorm:"primaryKey;type:varchar(50)"`
	activityMetrics20260212 `gorm:"embedded"`
	common.NoPKModel
}

func (metricsByIde20260212) TableName() string {
	return "_tool_copilot_metrics_by_ide"
}

type metricsByFeature20260212 struct {
	ConnectionId            uint64    `gorm:"primaryKey"`
	ScopeId                 string    `gorm:"primaryKey;type:varchar(255)"`
	Day                     time.Time `gorm:"primaryKey;type:date"`
	Feature                 string    `gorm:"primaryKey;type:varchar(100)"`
	activityMetrics20260212 `gorm:"embedded"`
	common.NoPKModel
}

func (metricsByFeature20260212) TableName() string {
	return "_tool_copilot_metrics_by_feature"
}

type metricsByLanguageFeature20260212 struct {
	ConnectionId         uint64    `gorm:"primaryKey"`
	ScopeId              string    `gorm:"primaryKey;type:varchar(255)"`
	Day                  time.Time `gorm:"primaryKey;type:date"`
	Language             string    `gorm:"primaryKey;type:varchar(50)"`
	Feature              string    `gorm:"primaryKey;type:varchar(100)"`
	codeMetrics20260212  `gorm:"embedded"`
	common.NoPKModel
}

func (metricsByLanguageFeature20260212) TableName() string {
	return "_tool_copilot_metrics_by_language_feature"
}

type metricsByLanguageModel20260212 struct {
	ConnectionId         uint64    `gorm:"primaryKey"`
	ScopeId              string    `gorm:"primaryKey;type:varchar(255)"`
	Day                  time.Time `gorm:"primaryKey;type:date"`
	Language             string    `gorm:"primaryKey;type:varchar(50)"`
	Model                string    `gorm:"primaryKey;type:varchar(100)"`
	codeMetrics20260212  `gorm:"embedded"`
	common.NoPKModel
}

func (metricsByLanguageModel20260212) TableName() string {
	return "_tool_copilot_metrics_by_language_model"
}

type metricsByModelFeature20260212 struct {
	ConnectionId            uint64    `gorm:"primaryKey"`
	ScopeId                 string    `gorm:"primaryKey;type:varchar(255)"`
	Day                     time.Time `gorm:"primaryKey;type:date"`
	Model                   string    `gorm:"primaryKey;type:varchar(100)"`
	Feature                 string    `gorm:"primaryKey;type:varchar(100)"`
	activityMetrics20260212 `gorm:"embedded"`
	common.NoPKModel
}

func (metricsByModelFeature20260212) TableName() string {
	return "_tool_copilot_metrics_by_model_feature"
}

// --- User metrics tables ---

type userDailyMetrics20260212 struct {
	ConnectionId            uint64    `gorm:"primaryKey"`
	ScopeId                 string    `gorm:"primaryKey;type:varchar(255)"`
	Day                     time.Time `gorm:"primaryKey;type:date"`
	UserId                  int64     `gorm:"primaryKey"`
	EnterpriseId            string    `gorm:"type:varchar(100)"`
	UserLogin               string    `gorm:"type:varchar(255);index"`
	UsedAgent               bool
	UsedChat                bool
	activityMetrics20260212 `gorm:"embedded"`
	common.NoPKModel
}

func (userDailyMetrics20260212) TableName() string {
	return "_tool_copilot_user_daily_metrics"
}

type userMetricsByIde20260212 struct {
	ConnectionId            uint64    `gorm:"primaryKey"`
	ScopeId                 string    `gorm:"primaryKey;type:varchar(255)"`
	Day                     time.Time `gorm:"primaryKey;type:date"`
	UserId                  int64     `gorm:"primaryKey"`
	Ide                     string    `gorm:"primaryKey;type:varchar(50)"`
	LastKnownPluginName     string    `gorm:"type:varchar(100)"`
	LastKnownPluginVersion  string    `gorm:"type:varchar(50)"`
	LastKnownIdeVersion     string    `gorm:"type:varchar(50)"`
	activityMetrics20260212 `gorm:"embedded"`
	common.NoPKModel
}

func (userMetricsByIde20260212) TableName() string {
	return "_tool_copilot_user_metrics_by_ide"
}

type userMetricsByFeature20260212 struct {
	ConnectionId            uint64    `gorm:"primaryKey"`
	ScopeId                 string    `gorm:"primaryKey;type:varchar(255)"`
	Day                     time.Time `gorm:"primaryKey;type:date"`
	UserId                  int64     `gorm:"primaryKey"`
	Feature                 string    `gorm:"primaryKey;type:varchar(100)"`
	activityMetrics20260212 `gorm:"embedded"`
	common.NoPKModel
}

func (userMetricsByFeature20260212) TableName() string {
	return "_tool_copilot_user_metrics_by_feature"
}

type userMetricsByLanguageFeature20260212 struct {
	ConnectionId        uint64    `gorm:"primaryKey"`
	ScopeId             string    `gorm:"primaryKey;type:varchar(255)"`
	Day                 time.Time `gorm:"primaryKey;type:date"`
	UserId              int64     `gorm:"primaryKey"`
	Language            string    `gorm:"primaryKey;type:varchar(50)"`
	Feature             string    `gorm:"primaryKey;type:varchar(100)"`
	codeMetrics20260212 `gorm:"embedded"`
	common.NoPKModel
}

func (userMetricsByLanguageFeature20260212) TableName() string {
	return "_tool_copilot_user_metrics_by_language_feature"
}

type userMetricsByLanguageModel20260212 struct {
	ConnectionId        uint64    `gorm:"primaryKey"`
	ScopeId             string    `gorm:"primaryKey;type:varchar(255)"`
	Day                 time.Time `gorm:"primaryKey;type:date"`
	UserId              int64     `gorm:"primaryKey"`
	Language            string    `gorm:"primaryKey;type:varchar(50)"`
	Model               string    `gorm:"primaryKey;type:varchar(100)"`
	codeMetrics20260212 `gorm:"embedded"`
	common.NoPKModel
}

func (userMetricsByLanguageModel20260212) TableName() string {
	return "_tool_copilot_user_metrics_by_language_model"
}

type userMetricsByModelFeature20260212 struct {
	ConnectionId            uint64    `gorm:"primaryKey"`
	ScopeId                 string    `gorm:"primaryKey;type:varchar(255)"`
	Day                     time.Time `gorm:"primaryKey;type:date"`
	UserId                  int64     `gorm:"primaryKey"`
	Model                   string    `gorm:"primaryKey;type:varchar(100)"`
	Feature                 string    `gorm:"primaryKey;type:varchar(100)"`
	activityMetrics20260212 `gorm:"embedded"`
	common.NoPKModel
}

func (userMetricsByModelFeature20260212) TableName() string {
	return "_tool_copilot_user_metrics_by_model_feature"
}

// --- New org metrics tables (replacing old ones) ---

type orgDailyMetrics20260212 struct {
	ConnectionId             uint64    `gorm:"primaryKey"`
	ScopeId                  string    `gorm:"primaryKey;type:varchar(255)"`
	Date                     time.Time `gorm:"primaryKey;type:date"`
	TotalActiveUsers         int
	TotalEngagedUsers        int
	CompletionSuggestions    int
	CompletionAcceptances    int
	CompletionLinesSuggested int
	CompletionLinesAccepted  int
	IdeChats                 int
	IdeChatCopyEvents        int
	IdeChatInsertionEvents   int
	IdeChatEngagedUsers      int
	DotcomChats              int
	DotcomChatEngagedUsers   int
	PRSummariesCreated       int
	PREngagedUsers           int
	SeatActiveCount          int
	SeatTotal                int
	common.NoPKModel
}

func (orgDailyMetrics20260212) TableName() string {
	return "_tool_copilot_org_daily_metrics"
}

type orgLanguageMetrics20260212 struct {
	ConnectionId uint64    `gorm:"primaryKey"`
	ScopeId      string    `gorm:"primaryKey;type:varchar(255)"`
	Date         time.Time `gorm:"primaryKey;type:date"`
	Editor       string    `gorm:"primaryKey;type:varchar(50)"`
	Language     string    `gorm:"primaryKey;type:varchar(50)"`
	EngagedUsers   int
	Suggestions    int
	Acceptances    int
	LinesSuggested int
	LinesAccepted  int
	common.NoPKModel
}

func (orgLanguageMetrics20260212) TableName() string {
	return "_tool_copilot_org_language_metrics"
}

func (script *migrateToUsageMetricsV2) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()

	// Drop legacy metrics tables (data must be re-collected from new API)
	if err := db.DropTables(
		"_tool_copilot_org_metrics",
		"_tool_copilot_language_metrics",
	); err != nil {
		basicRes.GetLogger().Warn(err, "Failed to drop legacy copilot tables (may not exist)")
	}

	// Add Enterprise column to connections and scopes
	if err := migrationhelper.AutoMigrateTables(basicRes,
		&connection20260212{},
		&scope20260212{},
	); err != nil {
		return err
	}

	// Create all new tables
	return migrationhelper.AutoMigrateTables(basicRes,
		// New org metrics (replacing dropped tables)
		&orgDailyMetrics20260212{},
		&orgLanguageMetrics20260212{},
		// Enterprise metrics
		&enterpriseDailyMetrics20260212{},
		&metricsByIde20260212{},
		&metricsByFeature20260212{},
		&metricsByLanguageFeature20260212{},
		&metricsByLanguageModel20260212{},
		&metricsByModelFeature20260212{},
		// User metrics
		&userDailyMetrics20260212{},
		&userMetricsByIde20260212{},
		&userMetricsByFeature20260212{},
		&userMetricsByLanguageFeature20260212{},
		&userMetricsByLanguageModel20260212{},
		&userMetricsByModelFeature20260212{},
	)
}

func (*migrateToUsageMetricsV2) Version() uint64 {
	return 20260212000000
}

func (*migrateToUsageMetricsV2) Name() string {
	return "Migrate GitHub Copilot to Usage Metrics Report API v2"
}
