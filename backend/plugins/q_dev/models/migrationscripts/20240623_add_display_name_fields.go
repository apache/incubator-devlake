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
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

var _ plugin.MigrationScript = (*addDisplayNameFields)(nil)

type addDisplayNameFields struct{}

func (*addDisplayNameFields) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&QDevUserDataWithDisplayName{},
		&QDevUserMetricsWithDisplayName{},
	)
}

func (*addDisplayNameFields) Version() uint64 {
	return 20240623000001
}

func (*addDisplayNameFields) Name() string {
	return "add display_name fields to user tables"
}

// Archived models for migration with display_name field added
type QDevUserDataWithDisplayName struct {
	archived.Model
	ConnectionId                      uint64    `gorm:"primaryKey"`
	UserId                            string    `gorm:"index" json:"userId"`
	Date                              time.Time `gorm:"index" json:"date"`
	DisplayName                       string    `gorm:"type:varchar(255)" json:"displayName"` // New field
	CodeReview_FindingsCount          int
	CodeReview_SucceededEventCount    int
	InlineChat_AcceptanceEventCount   int
	InlineChat_AcceptedLineAdditions  int
	InlineChat_AcceptedLineDeletions  int
	InlineChat_DismissalEventCount    int
	InlineChat_DismissedLineAdditions int
	InlineChat_DismissedLineDeletions int
	InlineChat_RejectedLineAdditions  int
	InlineChat_RejectedLineDeletions  int
	InlineChat_RejectionEventCount    int
	InlineChat_TotalEventCount        int
	Inline_AICodeLines                int
	Inline_AcceptanceCount            int
	Inline_SuggestionsCount           int
}

func (QDevUserDataWithDisplayName) TableName() string {
	return "_tool_q_dev_user_data"
}

type QDevUserMetricsWithDisplayName struct {
	archived.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey"`
	UserId       string `gorm:"primaryKey"`
	DisplayName  string `gorm:"type:varchar(255)" json:"displayName"` // New field
	FirstDate    time.Time
	LastDate     time.Time
	TotalDays    int

	// 聚合指标
	TotalCodeReview_FindingsCount          int
	TotalCodeReview_SucceededEventCount    int
	TotalInlineChat_AcceptanceEventCount   int
	TotalInlineChat_AcceptedLineAdditions  int
	TotalInlineChat_AcceptedLineDeletions  int
	TotalInlineChat_DismissalEventCount    int
	TotalInlineChat_DismissedLineAdditions int
	TotalInlineChat_DismissedLineDeletions int
	TotalInlineChat_RejectedLineAdditions  int
	TotalInlineChat_RejectedLineDeletions  int
	TotalInlineChat_RejectionEventCount    int
	TotalInlineChat_TotalEventCount        int
	TotalInline_AICodeLines                int
	TotalInline_AcceptanceCount            int
	TotalInline_SuggestionsCount           int

	// 平均指标
	AvgCodeReview_FindingsCount        float64
	AvgCodeReview_SucceededEventCount  float64
	AvgInlineChat_AcceptanceEventCount float64
	AvgInlineChat_TotalEventCount      float64
	AvgInline_AICodeLines              float64
	AvgInline_AcceptanceCount          float64
	AvgInline_SuggestionsCount         float64

	// 接受率指标
	AcceptanceRate float64
}

func (QDevUserMetricsWithDisplayName) TableName() string {
	return "_tool_q_dev_user_metrics"
}
