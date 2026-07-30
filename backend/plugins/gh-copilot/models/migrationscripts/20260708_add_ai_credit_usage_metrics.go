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
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type addAiCreditUsageMetrics struct{}

// --- Snapshot structs for migration (avoid importing models package to prevent drift) ---

type creditUsageBreakdown20260708 struct {
	GrossQuantity    float64
	DiscountQuantity float64
	NetQuantity      float64
	PricePerUnit     float64
	GrossAmount      float64
	DiscountAmount   float64
	NetAmount        float64
}

type enterpriseAiCreditUsage20260708 struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	ScopeId      string `gorm:"primaryKey;type:varchar(191)"`
	Year         int    `gorm:"primaryKey"`
	Month        int    `gorm:"primaryKey"`
	Day          int    `gorm:"primaryKey"`

	Enterprise   string `gorm:"primaryKey;type:varchar(191)"`
	Model        string `gorm:"primaryKey;type:varchar(191)"`
	Organization string `gorm:"index;type:varchar(255)"`
	User         string `gorm:"index;type:varchar(255)"`

	Product        string `gorm:"type:varchar(32)"`
	CostCenterId   string `gorm:"index;type:varchar(255)"`
	CostCenterName string `gorm:"type:varchar(255)"`

	creditUsageBreakdown20260708 `gorm:"embedded"`
	archived.NoPKModel
}

func (enterpriseAiCreditUsage20260708) TableName() string {
	return "_tool_copilot_enterprise_ai_credit_usage"
}

type orgAiCreditUsage20260708 struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	ScopeId      string `gorm:"primaryKey;type:varchar(191)"`
	Year         int    `gorm:"primaryKey"`
	Month        int    `gorm:"primaryKey"`
	Day          int    `gorm:"primaryKey"`

	Organization string `gorm:"primaryKey;type:varchar(191)"`
	Model        string `gorm:"primaryKey;type:varchar(191)"`
	User         string `gorm:"index;type:varchar(255)"`

	Product string `gorm:"type:varchar(32)"`

	creditUsageBreakdown20260708 `gorm:"embedded"`
	archived.NoPKModel
}

func (orgAiCreditUsage20260708) TableName() string {
	return "_tool_copilot_org_ai_credit_usage"
}

type userAiCreditUsage20260708 struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	ScopeId      string `gorm:"primaryKey;type:varchar(191)"`
	Year         int    `gorm:"primaryKey"`
	Month        int    `gorm:"primaryKey"`
	Day          int    `gorm:"primaryKey"`

	User  string `gorm:"primaryKey;type:varchar(191)"`
	Model string `gorm:"primaryKey;type:varchar(191)"`

	Product string `gorm:"type:varchar(32)"`

	creditUsageBreakdown20260708 `gorm:"embedded"`
	archived.NoPKModel
}

func (userAiCreditUsage20260708) TableName() string {
	return "_tool_copilot_user_ai_credit_usage"
}

// userDailyMetrics20260708 adds AI credit, CLI and code-review columns to the existing table.
type userDailyMetrics20260708 struct {
	ConnectionId uint64    `gorm:"primaryKey"`
	ScopeId      string    `gorm:"primaryKey;type:varchar(255)"`
	Day          time.Time `gorm:"primaryKey;type:date"`
	UserId       int64     `gorm:"primaryKey"`

	OrganizationId               string `gorm:"type:varchar(100)"`
	EnterpriseId                 string `gorm:"type:varchar(100)"`
	UserLogin                    string `gorm:"type:varchar(255);index"`
	UsedAgent                    bool
	UsedChat                     bool
	UsedCli                      bool    `gorm:"comment:Whether user used Copilot CLI"`
	UsedCopilotCodeReviewActive  bool    `gorm:"comment:Whether user actively used code review"`
	UsedCopilotCodeReviewPassive bool    `gorm:"comment:Whether user passively used code review"`
	AiCreditsUsed                float64 `gorm:"comment:AI credits consumed on this day"`

	UserInitiatedInteractionCount int
	CodeGenerationActivityCount   int
	CodeAcceptanceActivityCount   int
	LocSuggestedToAddSum          int
	LocSuggestedToDeleteSum       int
	LocAddedSum                   int
	LocDeletedSum                 int

	CliSessionCount   int
	CliRequestCount   int
	CliPromptCount    int
	CliOutputTokenSum int
	CliPromptTokenSum int

	archived.NoPKModel
}

func (userDailyMetrics20260708) TableName() string {
	return "_tool_copilot_user_daily_metrics"
}

func (u *addAiCreditUsageMetrics) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&enterpriseAiCreditUsage20260708{},
		&orgAiCreditUsage20260708{},
		&userAiCreditUsage20260708{},
		&userDailyMetrics20260708{},
	)
}

func (u *addAiCreditUsageMetrics) Version() uint64 {
	return 20260708000000
}

func (u *addAiCreditUsageMetrics) Name() string {
	return "add AI credit usage billing tables"
}
