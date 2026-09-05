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

// CursorUserSpend stores per-user billing cycle spend from /teams/spend.
type CursorUserSpend struct {
	ConnectionId      uint64    `gorm:"primaryKey" json:"connectionId"`
	ScopeId           string    `gorm:"primaryKey;type:varchar(255)" json:"scopeId"`
	UserId            string    `gorm:"primaryKey;type:varchar(255)" json:"userId"`
	BillingCycleStart time.Time `gorm:"primaryKey" json:"billingCycleStart"`
	CollectedAt       time.Time `gorm:"index" json:"collectedAt"`

	Email                    string  `json:"email" gorm:"type:varchar(255);index"`
	Name                     string  `json:"name" gorm:"type:varchar(255)"`
	Role                     string  `json:"role" gorm:"type:varchar(50)"`
	SpendCents               float64 `json:"spendCents"`
	IncludedSpendCents       float64 `json:"includedSpendCents"`
	FastPremiumRequests      int     `json:"fastPremiumRequests"`
	MonthlyLimitDollars      float64 `json:"monthlyLimitDollars"`
	HardLimitOverrideDollars float64 `json:"hardLimitOverrideDollars"`

	common.NoPKModel
}

func (CursorUserSpend) TableName() string {
	return "_tool_cursor_user_spend"
}
