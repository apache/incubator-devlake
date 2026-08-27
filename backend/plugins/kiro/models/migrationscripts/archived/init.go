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

// Package archived holds frozen copies of the model structs as they were when
// each migration was written. Migrations reference these instead of the live
// models so that editing a live model never changes what a historical migration
// does.
package archived

import (
	"time"

	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
)

type KiroConnection struct {
	archived.Model
	Name                string `gorm:"type:varchar(100);uniqueIndex" json:"name"`
	AccessKeyId         string `gorm:"type:varchar(255)" json:"accessKeyId"`
	SecretAccessKey     string `gorm:"type:varchar(255)" json:"secretAccessKey"`
	Region              string `gorm:"type:varchar(100)" json:"region"`
	Bucket              string `gorm:"type:varchar(255)" json:"bucket"`
	ReportPrefix        string `gorm:"type:varchar(512)" json:"reportPrefix"`
	PromptLogBucket     string `gorm:"type:varchar(255)" json:"promptLogBucket"`
	PromptLogPrefix     string `gorm:"type:varchar(512)" json:"promptLogPrefix"`
	IdentityStoreId     string `gorm:"type:varchar(255)" json:"identityStoreId"`
	IdentityStoreRegion string `gorm:"type:varchar(100)" json:"identityStoreRegion"`
	RateLimitPerHour    int    `json:"rateLimitPerHour"`
}

func (KiroConnection) TableName() string {
	return "_tool_kiro_connections"
}

type KiroS3Slice struct {
	archived.NoPKModel
	ConnectionId  uint64 `gorm:"primaryKey"`
	ScopeConfigId uint64 `json:"scopeConfigId,omitempty"`
	Id            string `gorm:"primaryKey;type:varchar(512)" json:"id"`
	BasePath      string `gorm:"type:varchar(512)" json:"basePath"`
	AccountId     string `gorm:"type:varchar(255);not null" json:"accountId"`
	Year          int    `gorm:"not null" json:"year"`
	Month         *int   `json:"month"`
}

func (KiroS3Slice) TableName() string {
	return "_tool_kiro_s3_slices"
}

type KiroS3FileMeta struct {
	archived.NoPKModel
	ConnectionId  uint64     `gorm:"primaryKey"`
	S3Path        string     `gorm:"primaryKey;type:varchar(512)"`
	FileName      string     `gorm:"type:varchar(255)" json:"fileName"`
	Bucket        string     `gorm:"type:varchar(255)" json:"bucket"`
	ScopeId       string     `gorm:"type:varchar(255);index" json:"scopeId"`
	FileType      string     `gorm:"type:varchar(32);index" json:"fileType"`
	Processed     bool       `gorm:"default:false" json:"processed"`
	ProcessedTime *time.Time `gorm:"default:null" json:"processedTime"`
	RecordCount   int        `gorm:"default:0" json:"recordCount"`
	ErrorMessage  string     `gorm:"type:text" json:"errorMessage"`
	AttemptCount  int        `gorm:"default:0" json:"attemptCount"`
}

func (KiroS3FileMeta) TableName() string {
	return "_tool_kiro_s3_file_meta"
}

type KiroUserReport struct {
	archived.NoPKModel
	ConnectionId       uint64    `gorm:"primaryKey"`
	ScopeId            string    `gorm:"primaryKey;type:varchar(255)"`
	UserId             string    `gorm:"primaryKey;type:varchar(64)"`
	Date               time.Time `gorm:"primaryKey;type:date"`
	ClientType         string    `gorm:"primaryKey;type:varchar(20)"`
	IdentityStoreId    string    `gorm:"type:varchar(32)" json:"identityStoreId"`
	UserEmail          *string   `gorm:"type:varchar(255);index" json:"userEmail"`
	DisplayName        *string   `gorm:"type:varchar(255)" json:"displayName"`
	ProfileId          string    `gorm:"type:varchar(512)" json:"profileId"`
	SubscriptionTier   string    `gorm:"type:varchar(50)" json:"subscriptionTier"`
	IsNewUser          *bool     `json:"isNewUser"`
	ChatConversations  int       `json:"chatConversations"`
	TotalMessages      int       `json:"totalMessages"`
	CreditsUsed        float64   `json:"creditsUsed"`
	OverageCap         float64   `json:"overageCap"`
	OverageCreditsUsed float64   `json:"overageCreditsUsed"`
	OverageEnabled     bool      `json:"overageEnabled"`
}

func (KiroUserReport) TableName() string {
	return "_tool_kiro_user_report"
}

type KiroUserModelMessage struct {
	archived.NoPKModel
	ConnectionId uint64    `gorm:"primaryKey"`
	ScopeId      string    `gorm:"primaryKey;type:varchar(255)"`
	UserId       string    `gorm:"primaryKey;type:varchar(64)"`
	Date         time.Time `gorm:"primaryKey;type:date"`
	ClientType   string    `gorm:"primaryKey;type:varchar(20)"`
	ModelName    string    `gorm:"primaryKey;type:varchar(100)"`
	MessageCount int       `json:"messageCount"`
}

func (KiroUserModelMessage) TableName() string {
	return "_tool_kiro_user_model_messages"
}

type KiroChatLog struct {
	archived.NoPKModel
	ConnectionId       uint64    `gorm:"primaryKey"`
	ScopeId            string    `gorm:"primaryKey;type:varchar(255)"`
	RequestId          string    `gorm:"primaryKey;type:varchar(64)"`
	UserId             string    `gorm:"type:varchar(64);index" json:"userId"`
	IdentityStoreId    string    `gorm:"type:varchar(32)" json:"identityStoreId"`
	Timestamp          time.Time `gorm:"type:datetime(6);index" json:"timestamp"`
	ChatTriggerType    string    `gorm:"type:varchar(20)" json:"chatTriggerType"`
	ModelId            *string   `gorm:"type:varchar(100)" json:"modelId"`
	HasPrompt          bool      `gorm:"index" json:"hasPrompt"`
	PromptLength       int       `json:"promptLength"`
	PromptSha256       *string   `gorm:"type:char(64);index" json:"promptSha256"`
	ResponseLength     int       `json:"responseLength"`
	HasFollowupPrompts bool      `json:"hasFollowupPrompts"`
	HasSteering        *bool     `json:"hasSteering"`
	IsSpecMode         *bool     `json:"isSpecMode"`
	ConversationId     *string   `gorm:"type:varchar(64)" json:"conversationId"`
	UtteranceId        *string   `gorm:"type:varchar(64)" json:"utteranceId"`
}

func (KiroChatLog) TableName() string {
	return "_tool_kiro_chat_log"
}

type KiroCompletionLog struct {
	archived.NoPKModel
	ConnectionId       uint64    `gorm:"primaryKey"`
	ScopeId            string    `gorm:"primaryKey;type:varchar(255)"`
	RequestId          string    `gorm:"primaryKey;type:varchar(64)"`
	UserId             string    `gorm:"type:varchar(64);index" json:"userId"`
	IdentityStoreId    string    `gorm:"type:varchar(32)" json:"identityStoreId"`
	Timestamp          time.Time `gorm:"type:datetime(6);index" json:"timestamp"`
	FileName           string    `gorm:"type:varchar(255);index" json:"fileName"`
	FileExtension      string    `gorm:"type:varchar(50)" json:"fileExtension"`
	HasCustomization   bool      `json:"hasCustomization"`
	CompletionsCount   int       `json:"completionsCount"`
	ReturnedCharCount  int       `json:"returnedCharCount"`
	ReturnedLineCount  int       `json:"returnedLineCount"`
	LeftContextLength  int       `json:"leftContextLength"`
	RightContextLength int       `json:"rightContextLength"`
}

func (KiroCompletionLog) TableName() string {
	return "_tool_kiro_completion_log"
}
