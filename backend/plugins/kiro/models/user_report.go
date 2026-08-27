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

// KiroUserReport is one row of the daily per-user activity report.
//
// Kiro writes one CSV per client type per day, so the same person appears once
// per client type they used - ClientType is therefore part of the key, not a
// descriptive column.
//
// Several fields are nullable because the report schema has grown over time:
// across the sampled history there are 19 distinct header layouts, and
// User_Email and New_User were both introduced partway through. Rows collected
// from earlier months legitimately lack them.
type KiroUserReport struct {
	common.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey"`
	ScopeId      string `gorm:"primaryKey;type:varchar(255)"`
	// UserId is the bare UUID with any identity-store prefix stripped. Logs
	// always carry the prefix and some report files do too, so both paths
	// normalize - otherwise one person yields two distinct ids and every
	// per-person aggregate is quietly wrong.
	UserId string    `gorm:"primaryKey;type:varchar(64)"`
	Date   time.Time `gorm:"primaryKey;type:date"`
	// ClientType is KIRO_IDE, KIRO_CLI, KIRO_WEB or PLUGIN. KIRO_WEB appears in
	// real exports but not in the published docs, so this is not validated
	// against a fixed set.
	ClientType string `gorm:"primaryKey;type:varchar(20)"`

	// IdentityStoreId is set only when the source value carried a prefix.
	IdentityStoreId string `gorm:"type:varchar(32)" json:"identityStoreId"`
	// UserEmail is the join key to git identity. Nullable: absent in early
	// report versions.
	UserEmail *string `gorm:"type:varchar(255);index" json:"userEmail"`
	// DisplayName comes from IAM Identity Center and is for display only.
	DisplayName *string `gorm:"type:varchar(255)" json:"displayName"`
	// ProfileId is a full ARN, not a short id.
	ProfileId string `gorm:"type:varchar(512)" json:"profileId"`
	// SubscriptionTier is normalized to UPPER_SNAKE (POWER, PRO_PLUS).
	SubscriptionTier string `gorm:"type:varchar(50)" json:"subscriptionTier"`
	// IsNewUser marks a subscription activated on the report date, which gives
	// an adoption timestamp directly instead of inferring one from first usage.
	// Nullable: absent in early report versions.
	IsNewUser *bool `json:"isNewUser"`

	ChatConversations int     `json:"chatConversations"`
	TotalMessages     int     `json:"totalMessages"`
	CreditsUsed       float64 `json:"creditsUsed"`
	OverageCap        float64 `json:"overageCap"`
	// OverageCreditsUsed is genuinely non-zero in real data.
	OverageCreditsUsed float64 `json:"overageCreditsUsed"`
	OverageEnabled     bool    `json:"overageEnabled"`
}

func (KiroUserReport) TableName() string {
	return "_tool_kiro_user_report"
}
