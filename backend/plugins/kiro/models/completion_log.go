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

// KiroCompletionLog is one GenerateCompletions (inline suggestion) request.
//
// The counters are named Returned*, not Accepted*: a record is written when the
// suggestion is requested, not when it is accepted. Empty completions arrays are
// common, so these numbers measure what Kiro offered rather than what the
// developer took.
//
// Note also that FileName is a bare file name with no directory path, so it
// cannot be resolved to a unique repository file.
type KiroCompletionLog struct {
	common.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey"`
	ScopeId      string `gorm:"primaryKey;type:varchar(255)"`
	RequestId    string `gorm:"primaryKey;type:varchar(64)"`

	UserId          string    `gorm:"type:varchar(64);index" json:"userId"`
	IdentityStoreId string    `gorm:"type:varchar(32)" json:"identityStoreId"`
	Timestamp       time.Time `gorm:"type:datetime(6);index" json:"timestamp"`

	// FileName has no path component; see the type comment.
	FileName      string `gorm:"type:varchar(255);index" json:"fileName"`
	FileExtension string `gorm:"type:varchar(50)" json:"fileExtension"`
	// HasCustomization records whether a customization ARN was attached. Present
	// on completion records but never on chat records.
	HasCustomization bool `json:"hasCustomization"`

	// CompletionsCount is often 0 - the request was logged, no suggestion came
	// back. Such records are still stored.
	CompletionsCount   int `json:"completionsCount"`
	ReturnedCharCount  int `json:"returnedCharCount"`
	ReturnedLineCount  int `json:"returnedLineCount"`
	LeftContextLength  int `json:"leftContextLength"`
	RightContextLength int `json:"rightContextLength"`
}

func (KiroCompletionLog) TableName() string {
	return "_tool_kiro_completion_log"
}
