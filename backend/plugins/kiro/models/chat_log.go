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

// KiroChatLog is one GenerateAssistantResponse interaction.
//
// Neither the prompt nor the assistant response text is stored. Derived
// features are computed from them during extraction and the originals are
// discarded, which keeps proprietary code and personal content out of the
// warehouse while preserving the signal.
//
// RequestId is the key rather than (UserId, Timestamp) because log timestamps
// carry nanosecond precision that MySQL DATETIME(6) truncates - two records
// within the same microsecond would collide.
type KiroChatLog struct {
	common.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey"`
	ScopeId      string `gorm:"primaryKey;type:varchar(255)"`
	RequestId    string `gorm:"primaryKey;type:varchar(64)"`

	// UserId is the bare UUID; log records always carry an identity-store
	// prefix, which is stripped so this joins against the report table.
	UserId string `gorm:"type:varchar(64);index" json:"userId"`
	// IdentityStoreId is the stripped prefix, retained for auditability.
	IdentityStoreId string    `gorm:"type:varchar(32)" json:"identityStoreId"`
	Timestamp       time.Time `gorm:"type:datetime(6);index" json:"timestamp"`
	// ChatTriggerType is MANUAL or INLINE_CHAT per the docs; only MANUAL has
	// been observed. Not validated against a fixed set.
	ChatTriggerType string `gorm:"type:varchar(20)" json:"chatTriggerType"`
	// ModelId is present on only about 45% of records, hence nullable. Do not
	// use it for model attribution; see KiroUserModelMessage.
	ModelId *string `gorm:"type:varchar(100)" json:"modelId"`

	// HasPrompt distinguishes a user turn from an agent self-continuation.
	//
	// An empty prompt means the agent continued on its own after a tool call;
	// a non-empty one means the user spoke. The ratio between them measures how
	// much of the traffic the agent generates versus the human - roughly 71% of
	// sampled records have no prompt.
	HasPrompt    bool `gorm:"index" json:"hasPrompt"`
	PromptLength int  `json:"promptLength"`
	// PromptSha256 identifies the same prompt being resubmitted, which is a
	// rework signal. NULL when the prompt is empty: hashing the empty string
	// would give ~71% of rows one shared hash and destroy the signal.
	PromptSha256 *string `gorm:"type:char(64);index" json:"promptSha256"`

	ResponseLength     int  `json:"responseLength"`
	HasFollowupPrompts bool `json:"hasFollowupPrompts"`

	// HasSteering and IsSpecMode are heuristics derived from prompt text, so
	// they are meaningful only when HasPrompt is true and NULL otherwise -
	// storing false would be indistinguishable from a real negative.
	HasSteering *bool `json:"hasSteering"`
	IsSpecMode  *bool `json:"isSpecMode"`

	// ConversationId and UtteranceId are documented but absent from every
	// sampled record. Kept as nullable columns in case they appear later; their
	// absence is why the S3 logs cannot group interactions into sessions.
	ConversationId *string `gorm:"type:varchar(64)" json:"conversationId"`
	UtteranceId    *string `gorm:"type:varchar(64)" json:"utteranceId"`
}

func (KiroChatLog) TableName() string {
	return "_tool_kiro_chat_log"
}
