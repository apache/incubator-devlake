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

type addCursorInitialTables struct{}

func (script *addCursorInitialTables) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&cursorConnection20260709{},
		&cursorScope20260709{},
		&cursorScopeConfig20260709{},
		&cursorUsageEvent20260709{},
		&cursorUserSpend20260709{},
		&cursorMember20260709{},
		&cursorRawUsageEvents20260709{},
		&cursorRawUserSpend20260709{},
		&cursorRawMembers20260709{},
	)
}

type cursorConnection20260709 struct {
	archived.Model
	Name             string `gorm:"type:varchar(100);uniqueIndex" json:"name"`
	Endpoint         string `gorm:"type:varchar(255)" json:"endpoint"`
	Proxy            string `gorm:"type:varchar(255)" json:"proxy"`
	RateLimitPerHour int    `json:"rateLimitPerHour"`
	Token            string `json:"token"`
}

func (cursorConnection20260709) TableName() string { return "_tool_cursor_connections" }

type cursorScope20260709 struct {
	archived.NoPKModel
	ConnectionId  uint64 `json:"connectionId" gorm:"primaryKey"`
	ScopeConfigId uint64 `json:"scopeConfigId,omitempty"`
	Id            string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	TeamId        string `json:"teamId" gorm:"type:varchar(255)"`
	Name          string `json:"name" gorm:"type:varchar(255)"`
	FullName      string `json:"fullName" gorm:"type:varchar(255)"`
}

func (cursorScope20260709) TableName() string { return "_tool_cursor_scopes" }

type cursorScopeConfig20260709 struct {
	archived.Model
	Entities     []string `gorm:"type:json;serializer:json" json:"entities" mapstructure:"entities"`
	ConnectionId uint64   `json:"connectionId" gorm:"index" validate:"required" mapstructure:"connectionId,omitempty"`
	Name         string   `mapstructure:"name" json:"name" gorm:"type:varchar(255);uniqueIndex" validate:"required"`
}

func (cursorScopeConfig20260709) TableName() string { return "_tool_cursor_scope_configs" }

type cursorUsageEvent20260709 struct {
	ConnectionId     uint64    `gorm:"primaryKey"`
	ScopeId          string    `gorm:"primaryKey;type:varchar(255)"`
	EventId          string    `gorm:"primaryKey;type:varchar(64)"`
	EventTime        time.Time `gorm:"index"`
	UserEmail        string    `gorm:"type:varchar(255);index"`
	Model            string    `gorm:"type:varchar(255);index"`
	Kind             string    `gorm:"type:varchar(100)"`
	ConversationId   string    `gorm:"type:varchar(64);index"`
	ChargedCents     float64
	RequestsCosts    int
	IsTokenBasedCall bool
	IsChargeable     bool
	MaxMode          bool
	IsHeadless       bool
	InputTokens      int
	OutputTokens     int
	CacheReadTokens  int
	CacheWriteTokens int
	TotalCents       float64
	CursorTokenFee   float64
	HostingType      string `gorm:"type:varchar(50)"`
	ServiceAccountId string `gorm:"type:varchar(255)"`
	archived.NoPKModel
}

func (cursorUsageEvent20260709) TableName() string { return "_tool_cursor_usage_events" }

type cursorUserSpend20260709 struct {
	ConnectionId             uint64    `gorm:"primaryKey"`
	ScopeId                  string    `gorm:"primaryKey;type:varchar(255)"`
	UserId                   string    `gorm:"primaryKey;type:varchar(255)"`
	BillingCycleStart        time.Time `gorm:"primaryKey"`
	CollectedAt              time.Time `gorm:"index"`
	Email                    string    `gorm:"type:varchar(255);index"`
	Name                     string    `gorm:"type:varchar(255)"`
	Role                     string    `gorm:"type:varchar(50)"`
	SpendCents               float64
	IncludedSpendCents       float64
	FastPremiumRequests      int
	MonthlyLimitDollars      float64
	HardLimitOverrideDollars float64
	archived.NoPKModel
}

func (cursorUserSpend20260709) TableName() string { return "_tool_cursor_user_spend" }

type cursorMember20260709 struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	ScopeId      string `gorm:"primaryKey;type:varchar(255)"`
	UserId       string `gorm:"primaryKey;type:varchar(255)"`
	Email        string `gorm:"type:varchar(255);index"`
	Name         string `gorm:"type:varchar(255)"`
	Role         string `gorm:"type:varchar(50)"`
	IsRemoved    bool
	archived.NoPKModel
}

func (cursorMember20260709) TableName() string { return "_tool_cursor_members" }

type cursorRawUsageEvents20260709 struct {
	ID        uint64 `gorm:"primaryKey"`
	Params    string `gorm:"type:varchar(255);index"`
	Data      []byte
	Url       string
	Input     json.RawMessage `gorm:"type:json"`
	CreatedAt time.Time       `gorm:"index"`
}

func (cursorRawUsageEvents20260709) TableName() string { return "_raw_cursor_usage_events" }

type cursorRawUserSpend20260709 struct {
	ID        uint64 `gorm:"primaryKey"`
	Params    string `gorm:"type:varchar(255);index"`
	Data      []byte
	Url       string
	Input     json.RawMessage `gorm:"type:json"`
	CreatedAt time.Time       `gorm:"index"`
}

func (cursorRawUserSpend20260709) TableName() string { return "_raw_cursor_user_spend" }

type cursorRawMembers20260709 struct {
	ID        uint64 `gorm:"primaryKey"`
	Params    string `gorm:"type:varchar(255);index"`
	Data      []byte
	Url       string
	Input     json.RawMessage `gorm:"type:json"`
	CreatedAt time.Time       `gorm:"index"`
}

func (cursorRawMembers20260709) TableName() string { return "_raw_cursor_members" }

func (*addCursorInitialTables) Version() uint64 { return 20260709000000 }

func (*addCursorInitialTables) Name() string { return "cursor init tables" }
