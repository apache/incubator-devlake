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
	"github.com/apache/devlake/core/context"
	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/models/migrationscripts/archived"
	"github.com/apache/devlake/helpers/migrationhelper"
)

// The 20260708 migration embedded the credit-usage breakdown columns through an
// *unexported* anonymous struct (`creditUsageBreakdown20260708`). GORM's schema
// parser does not migrate the fields of an unexported embedded type, so the
// gross_/discount_/net_ quantity and amount columns (plus price_per_unit) were
// never created even though the runtime models declare them inline. This
// follow-up migration adds the missing columns. AutoMigrate only adds columns
// that do not yet exist, so it is a no-op on databases that somehow already
// have them.

type enterpriseAiCreditUsageBreakdown20260731 struct {
	GrossQuantity    float64 `gorm:"comment:Raw credits consumed"`
	DiscountQuantity float64 `gorm:"comment:Credits discounted"`
	NetQuantity      float64 `gorm:"comment:Credits after discount"`
	PricePerUnit     float64 `gorm:"comment:Price per credit unit"`
	GrossAmount      float64 `gorm:"comment:Gross cost before discount"`
	DiscountAmount   float64 `gorm:"comment:Discount amount"`
	NetAmount        float64 `gorm:"comment:Net cost after discount"`
	archived.NoPKModel
}

func (enterpriseAiCreditUsageBreakdown20260731) TableName() string {
	return "_tool_copilot_enterprise_ai_credit_usage"
}

type orgAiCreditUsageBreakdown20260731 struct {
	GrossQuantity    float64 `gorm:"comment:Raw credits consumed"`
	DiscountQuantity float64 `gorm:"comment:Credits discounted"`
	NetQuantity      float64 `gorm:"comment:Credits after discount"`
	PricePerUnit     float64 `gorm:"comment:Price per credit unit"`
	GrossAmount      float64 `gorm:"comment:Gross cost before discount"`
	DiscountAmount   float64 `gorm:"comment:Discount amount"`
	NetAmount        float64 `gorm:"comment:Net cost after discount"`
	archived.NoPKModel
}

func (orgAiCreditUsageBreakdown20260731) TableName() string {
	return "_tool_copilot_org_ai_credit_usage"
}

type userAiCreditUsageBreakdown20260731 struct {
	GrossQuantity    float64 `gorm:"comment:Raw credits consumed"`
	DiscountQuantity float64 `gorm:"comment:Credits discounted"`
	NetQuantity      float64 `gorm:"comment:Credits after discount"`
	PricePerUnit     float64 `gorm:"comment:Price per credit unit"`
	GrossAmount      float64 `gorm:"comment:Gross cost before discount"`
	DiscountAmount   float64 `gorm:"comment:Discount amount"`
	NetAmount        float64 `gorm:"comment:Net cost after discount"`
	archived.NoPKModel
}

func (userAiCreditUsageBreakdown20260731) TableName() string {
	return "_tool_copilot_user_ai_credit_usage"
}

type fixAiCreditUsageBreakdownColumns struct{}

func (u *fixAiCreditUsageBreakdownColumns) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&enterpriseAiCreditUsageBreakdown20260731{},
		&orgAiCreditUsageBreakdown20260731{},
		&userAiCreditUsageBreakdown20260731{},
	)
}

func (u *fixAiCreditUsageBreakdownColumns) Version() uint64 {
	return 20260731000000
}

func (u *fixAiCreditUsageBreakdownColumns) Name() string {
	return "add missing AI credit usage breakdown columns"
}
