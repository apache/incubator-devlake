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
	"encoding/json"
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type addCursorDailyUsage struct{}

type cursorDailyUsage20260714 struct {
	ConnectionId             uint64    `gorm:"primaryKey"`
	ScopeId                  string    `gorm:"primaryKey;type:varchar(255)"`
	UserId                   string    `gorm:"primaryKey;type:varchar(255)"`
	UsageDate                time.Time `gorm:"primaryKey"`
	Email                    string    `gorm:"type:varchar(255);index"`
	IsActive                 bool
	Completions              int
	PremiumRequests          int
	AgentRequests            int
	ChatRequests             int
	ComposerRequests         int
	TabsAccepted             int
	TabsShown                int
	UsageBasedReqs           int
	SubscriptionIncludedReqs int
	MostUsedModel            string `gorm:"type:varchar(255)"`
	ClientVersion            string `gorm:"type:varchar(100)"`
	LinesAdded               int
	LinesDeleted             int
	archived.NoPKModel
}

func (cursorDailyUsage20260714) TableName() string { return "_tool_cursor_daily_usage" }

type cursorRawDailyUsage20260714 struct {
	ID        uint64 `gorm:"primaryKey"`
	Params    string `gorm:"type:varchar(255);index"`
	Data      []byte
	Url       string
	Input     json.RawMessage `gorm:"type:json"`
	CreatedAt time.Time       `gorm:"index"`
}

func (cursorRawDailyUsage20260714) TableName() string { return "_raw_cursor_daily_usage" }

func (script *addCursorDailyUsage) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&cursorDailyUsage20260714{},
		&cursorRawDailyUsage20260714{},
	)
}

func (*addCursorDailyUsage) Version() uint64 { return 20260714000000 }

func (*addCursorDailyUsage) Name() string { return "cursor add daily usage tables" }
