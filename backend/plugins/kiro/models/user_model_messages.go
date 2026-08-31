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

// KiroUserModelMessage holds the per-model message counts from the report CSV.
//
// These arrive as dynamic columns ("claude_opus_4.6_messages") whose names and
// count vary with the models a team uses, so they cannot live in a fixed struct.
// A narrow table also means a new model appears in the data without a migration.
//
// This table - not KiroChatLog.ModelId - is the authoritative source for model
// attribution: the interaction log carries ModelId on only 45% of records, so
// any share-of-usage computed from it is systematically skewed.
type KiroUserModelMessage struct {
	common.NoPKModel
	ConnectionId uint64    `gorm:"primaryKey"`
	ScopeId      string    `gorm:"primaryKey;type:varchar(255)"`
	UserId       string    `gorm:"primaryKey;type:varchar(64)"`
	Date         time.Time `gorm:"primaryKey;type:date"`
	ClientType   string    `gorm:"primaryKey;type:varchar(20)"`
	// ModelName keeps the CSV's underscore spelling ("claude_opus_4.6"). The
	// interaction log spells the same model with hyphens; values are stored as
	// each source emits them and reconciled at query time, so the stored value
	// always traces back to its source.
	//
	// Note that some values are routing modes rather than models ("auto",
	// "simple_task"). Upstream presents them through the same column pattern,
	// so no distinction is invented here.
	ModelName string `gorm:"primaryKey;type:varchar(100)"`

	MessageCount int `json:"messageCount"`
}

func (KiroUserModelMessage) TableName() string {
	return "_tool_kiro_user_model_messages"
}
