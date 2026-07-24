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
	"github.com/apache/incubator-devlake/core/models/common"
)

// GhCopilotUserAiCreditUsage tracks AI credit consumption at the individual user level.
// One row per time period per model per authenticated user.
type GhCopilotUserAiCreditUsage struct {
	ConnectionId uint64 `gorm:"primaryKey" json:"connectionId"`
	ScopeId      string `gorm:"primaryKey;type:varchar(191)" json:"scopeId"`
	Year         int    `gorm:"primaryKey" json:"year"`
	Month        int    `gorm:"primaryKey" json:"month"`
	Day          int    `gorm:"primaryKey" json:"day"`

	User  string `gorm:"primaryKey;type:varchar(191);comment:GitHub username" json:"user"`
	Model string `gorm:"primaryKey;type:varchar(191);comment:AI model name (e.g., gpt-4.1)" json:"model"`

	Product string `gorm:"type:varchar(32);comment:Product name (e.g., copilot)" json:"product"`

	// Credit usage breakdown
	GrossQuantity    float64 `json:"grossQuantity" gorm:"comment:Raw credits consumed"`
	DiscountQuantity float64 `json:"discountQuantity" gorm:"comment:Credits discounted"`
	NetQuantity      float64 `json:"netQuantity" gorm:"comment:Credits after discount"`
	PricePerUnit     float64 `json:"pricePerUnit" gorm:"comment:Price per credit unit"`
	GrossAmount      float64 `json:"grossAmount" gorm:"comment:Gross cost before discount"`
	DiscountAmount   float64 `json:"discountAmount" gorm:"comment:Discount amount"`
	NetAmount        float64 `json:"netAmount" gorm:"comment:Net cost after discount"`

	common.NoPKModel
}

func (GhCopilotUserAiCreditUsage) TableName() string {
	return "_tool_copilot_user_ai_credit_usage"
}
