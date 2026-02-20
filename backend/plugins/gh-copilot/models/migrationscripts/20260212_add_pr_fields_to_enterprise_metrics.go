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

type addPRFieldsToEnterpriseMetrics struct{}

// Properly exported embedded struct for activity metrics
type ActivityMetrics20260212v2 struct {
	UserInitiatedInteractionCount int
	CodeGenerationActivityCount   int
	CodeAcceptanceActivityCount   int
	LocSuggestedToAddSum          int
	LocSuggestedToDeleteSum       int
	LocAddedSum                   int
	LocDeletedSum                 int
}

type CodeMetrics20260212v2 struct {
	CodeGenerationActivityCount int
	CodeAcceptanceActivityCount int
	LocSuggestedToAddSum        int
	LocSuggestedToDeleteSum     int
	LocAddedSum                 int
	LocDeletedSum               int
}

// Enterprise daily metrics with all fields including PR stats
type enterpriseDailyMetrics20260212v2 struct {
	ConnectionId                  uint64    `gorm:"primaryKey"`
	ScopeId                       string    `gorm:"primaryKey;type:varchar(255)"`
	Day                           time.Time `gorm:"primaryKey;type:date"`
	EnterpriseId                  string    `gorm:"type:varchar(100)"`
	DailyActiveUsers              int
	WeeklyActiveUsers             int
	MonthlyActiveUsers            int
	MonthlyActiveChatUsers        int
	MonthlyActiveAgentUsers       int
	UserInitiatedInteractionCount int
	CodeGenerationActivityCount   int
	CodeAcceptanceActivityCount   int
	LocSuggestedToAddSum          int
	LocSuggestedToDeleteSum       int
	LocAddedSum                   int
	LocDeletedSum                 int
	PRTotalReviewed               int
	PRTotalCreated                int
	PRTotalCreatedByCopilot       int
	PRTotalReviewedByCopilot      int
	common.NoPKModel
}

func (enterpriseDailyMetrics20260212v2) TableName() string {
	return "_tool_copilot_enterprise_daily_metrics"
}

type metricsByIde20260212v2 struct {
	ConnectionId                  uint64    `gorm:"primaryKey"`
	ScopeId                       string    `gorm:"primaryKey;type:varchar(255)"`
	Day                           time.Time `gorm:"primaryKey;type:date"`
	Ide                           string    `gorm:"primaryKey;type:varchar(50)"`
	UserInitiatedInteractionCount int
	CodeGenerationActivityCount   int
	CodeAcceptanceActivityCount   int
	LocSuggestedToAddSum          int
	LocSuggestedToDeleteSum       int
	LocAddedSum                   int
	LocDeletedSum                 int
	common.NoPKModel
}

func (metricsByIde20260212v2) TableName() string {
	return "_tool_copilot_metrics_by_ide"
}

type metricsByFeature20260212v2 struct {
	ConnectionId                  uint64    `gorm:"primaryKey"`
	ScopeId                       string    `gorm:"primaryKey;type:varchar(255)"`
	Day                           time.Time `gorm:"primaryKey;type:date"`
	Feature                       string    `gorm:"primaryKey;type:varchar(100)"`
	UserInitiatedInteractionCount int
	CodeGenerationActivityCount   int
	CodeAcceptanceActivityCount   int
	LocSuggestedToAddSum          int
	LocSuggestedToDeleteSum       int
	LocAddedSum                   int
	LocDeletedSum                 int
	common.NoPKModel
}

func (metricsByFeature20260212v2) TableName() string {
	return "_tool_copilot_metrics_by_feature"
}

type metricsByLanguageFeature20260212v2 struct {
	ConnectionId            uint64    `gorm:"primaryKey"`
	ScopeId                 string    `gorm:"primaryKey;type:varchar(255)"`
	Day                     time.Time `gorm:"primaryKey;type:date"`
	Language                string    `gorm:"primaryKey;type:varchar(50)"`
	Feature                 string    `gorm:"primaryKey;type:varchar(100)"`
	CodeGenerationActivityCount int
	CodeAcceptanceActivityCount int
	LocSuggestedToAddSum        int
	LocSuggestedToDeleteSum     int
	LocAddedSum                 int
	LocDeletedSum               int
	common.NoPKModel
}

func (metricsByLanguageFeature20260212v2) TableName() string {
	return "_tool_copilot_metrics_by_language_feature"
}

type metricsByLanguageModel20260212v2 struct {
	ConnectionId            uint64    `gorm:"primaryKey"`
	ScopeId                 string    `gorm:"primaryKey;type:varchar(255)"`
	Day                     time.Time `gorm:"primaryKey;type:date"`
	Language                string    `gorm:"primaryKey;type:varchar(50)"`
	Model                   string    `gorm:"primaryKey;type:varchar(100)"`
	CodeGenerationActivityCount int
	CodeAcceptanceActivityCount int
	LocSuggestedToAddSum        int
	LocSuggestedToDeleteSum     int
	LocAddedSum                 int
	LocDeletedSum               int
	common.NoPKModel
}

func (metricsByLanguageModel20260212v2) TableName() string {
	return "_tool_copilot_metrics_by_language_model"
}

type metricsByModelFeature20260212v2 struct {
	ConnectionId                  uint64    `gorm:"primaryKey"`
	ScopeId                       string    `gorm:"primaryKey;type:varchar(255)"`
	Day                           time.Time `gorm:"primaryKey;type:date"`
	Model                         string    `gorm:"primaryKey;type:varchar(100)"`
	Feature                       string    `gorm:"primaryKey;type:varchar(100)"`
	UserInitiatedInteractionCount int
	CodeGenerationActivityCount   int
	CodeAcceptanceActivityCount   int
	LocSuggestedToAddSum          int
	LocSuggestedToDeleteSum       int
	LocAddedSum                   int
	LocDeletedSum                 int
	common.NoPKModel
}

func (metricsByModelFeature20260212v2) TableName() string {
	return "_tool_copilot_metrics_by_model_feature"
}

// User metrics tables
type userDailyMetrics20260212v2 struct {
	ConnectionId                  uint64    `gorm:"primaryKey"`
	ScopeId                       string    `gorm:"primaryKey;type:varchar(255)"`
	Day                           time.Time `gorm:"primaryKey;type:date"`
	UserId                        int64     `gorm:"primaryKey"`
	EnterpriseId                  string    `gorm:"type:varchar(100)"`
	UserLogin                     string    `gorm:"type:varchar(255);index"`
	UsedAgent                     bool
	UsedChat                      bool
	UserInitiatedInteractionCount int
	CodeGenerationActivityCount   int
	CodeAcceptanceActivityCount   int
	LocSuggestedToAddSum          int
	LocSuggestedToDeleteSum       int
	LocAddedSum                   int
	LocDeletedSum                 int
	common.NoPKModel
}

func (userDailyMetrics20260212v2) TableName() string {
	return "_tool_copilot_user_daily_metrics"
}

type userMetricsByIde20260212v2 struct {
	ConnectionId                  uint64    `gorm:"primaryKey"`
	ScopeId                       string    `gorm:"primaryKey;type:varchar(255)"`
	Day                           time.Time `gorm:"primaryKey;type:date"`
	UserId                        int64     `gorm:"primaryKey"`
	Ide                           string    `gorm:"primaryKey;type:varchar(50)"`
	LastKnownPluginName           string    `gorm:"type:varchar(100)"`
	LastKnownPluginVersion        string    `gorm:"type:varchar(50)"`
	LastKnownIdeVersion           string    `gorm:"type:varchar(50)"`
	UserInitiatedInteractionCount int
	CodeGenerationActivityCount   int
	CodeAcceptanceActivityCount   int
	LocSuggestedToAddSum          int
	LocSuggestedToDeleteSum       int
	LocAddedSum                   int
	LocDeletedSum                 int
	common.NoPKModel
}

func (userMetricsByIde20260212v2) TableName() string {
	return "_tool_copilot_user_metrics_by_ide"
}

type userMetricsByFeature20260212v2 struct {
	ConnectionId                  uint64    `gorm:"primaryKey"`
	ScopeId                       string    `gorm:"primaryKey;type:varchar(255)"`
	Day                           time.Time `gorm:"primaryKey;type:date"`
	UserId                        int64     `gorm:"primaryKey"`
	Feature                       string    `gorm:"primaryKey;type:varchar(100)"`
	UserInitiatedInteractionCount int
	CodeGenerationActivityCount   int
	CodeAcceptanceActivityCount   int
	LocSuggestedToAddSum          int
	LocSuggestedToDeleteSum       int
	LocAddedSum                   int
	LocDeletedSum                 int
	common.NoPKModel
}

func (userMetricsByFeature20260212v2) TableName() string {
	return "_tool_copilot_user_metrics_by_feature"
}

type userMetricsByLanguageFeature20260212v2 struct {
	ConnectionId            uint64    `gorm:"primaryKey"`
	ScopeId                 string    `gorm:"primaryKey;type:varchar(255)"`
	Day                     time.Time `gorm:"primaryKey;type:date"`
	UserId                  int64     `gorm:"primaryKey"`
	Language                string    `gorm:"primaryKey;type:varchar(50)"`
	Feature                 string    `gorm:"primaryKey;type:varchar(100)"`
	CodeGenerationActivityCount int
	CodeAcceptanceActivityCount int
	LocSuggestedToAddSum        int
	LocSuggestedToDeleteSum     int
	LocAddedSum                 int
	LocDeletedSum               int
	common.NoPKModel
}

func (userMetricsByLanguageFeature20260212v2) TableName() string {
	return "_tool_copilot_user_metrics_by_language_feature"
}

type userMetricsByLanguageModel20260212v2 struct {
	ConnectionId            uint64    `gorm:"primaryKey"`
	ScopeId                 string    `gorm:"primaryKey;type:varchar(255)"`
	Day                     time.Time `gorm:"primaryKey;type:date"`
	UserId                  int64     `gorm:"primaryKey"`
	Language                string    `gorm:"primaryKey;type:varchar(50)"`
	Model                   string    `gorm:"primaryKey;type:varchar(100)"`
	CodeGenerationActivityCount int
	CodeAcceptanceActivityCount int
	LocSuggestedToAddSum        int
	LocSuggestedToDeleteSum     int
	LocAddedSum                 int
	LocDeletedSum               int
	common.NoPKModel
}

func (userMetricsByLanguageModel20260212v2) TableName() string {
	return "_tool_copilot_user_metrics_by_language_model"
}

type userMetricsByModelFeature20260212v2 struct {
	ConnectionId                  uint64    `gorm:"primaryKey"`
	ScopeId                       string    `gorm:"primaryKey;type:varchar(255)"`
	Day                           time.Time `gorm:"primaryKey;type:date"`
	UserId                        int64     `gorm:"primaryKey"`
	Model                         string    `gorm:"primaryKey;type:varchar(100)"`
	Feature                       string    `gorm:"primaryKey;type:varchar(100)"`
	UserInitiatedInteractionCount int
	CodeGenerationActivityCount   int
	CodeAcceptanceActivityCount   int
	LocSuggestedToAddSum          int
	LocSuggestedToDeleteSum       int
	LocAddedSum                   int
	LocDeletedSum                 int
	common.NoPKModel
}

func (userMetricsByModelFeature20260212v2) TableName() string {
	return "_tool_copilot_user_metrics_by_model_feature"
}

func (script *addPRFieldsToEnterpriseMetrics) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()

	// Drop all metric tables created by v2 migration (they have missing columns
	// due to unexported embedded struct bug in GORM). Data will be re-collected.
	if err := db.DropTables(
		"_tool_copilot_enterprise_daily_metrics",
		"_tool_copilot_metrics_by_ide",
		"_tool_copilot_metrics_by_feature",
		"_tool_copilot_metrics_by_language_feature",
		"_tool_copilot_metrics_by_language_model",
		"_tool_copilot_metrics_by_model_feature",
		"_tool_copilot_user_daily_metrics",
		"_tool_copilot_user_metrics_by_ide",
		"_tool_copilot_user_metrics_by_feature",
		"_tool_copilot_user_metrics_by_language_feature",
		"_tool_copilot_user_metrics_by_language_model",
		"_tool_copilot_user_metrics_by_model_feature",
	); err != nil {
		basicRes.GetLogger().Warn(err, "Failed to drop tables for recreation")
	}

	// Recreate with all columns properly defined (no embedded structs)
	return migrationhelper.AutoMigrateTables(basicRes,
		&enterpriseDailyMetrics20260212v2{},
		&metricsByIde20260212v2{},
		&metricsByFeature20260212v2{},
		&metricsByLanguageFeature20260212v2{},
		&metricsByLanguageModel20260212v2{},
		&metricsByModelFeature20260212v2{},
		&userDailyMetrics20260212v2{},
		&userMetricsByIde20260212v2{},
		&userMetricsByFeature20260212v2{},
		&userMetricsByLanguageFeature20260212v2{},
		&userMetricsByLanguageModel20260212v2{},
		&userMetricsByModelFeature20260212v2{},
	)
}

func (*addPRFieldsToEnterpriseMetrics) Version() uint64 {
	return 20260212100000
}

func (*addPRFieldsToEnterpriseMetrics) Name() string {
	return "Recreate metric tables with all columns and PR fields"
}
