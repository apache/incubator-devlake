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

// File types, derived from the S3 prefix a file was found under. Storing the
// type lets extractors select their input without re-parsing paths.
const (
	FileTypeReport        = "report"
	FileTypeChatLog       = "chat_log"
	FileTypeCompletionLog = "completion_log"
)

// MaxAttempts caps retries of a failing file.
//
// Without a cap, a permanently malformed object is retried on every run: the
// log fills with the same error and the scope never reaches a finished state.
const MaxAttempts = 3

// KiroS3FileMeta is the incremental collection cursor - one row per S3 object.
//
// The primary key is (ConnectionId, S3Path), NOT the file name. Existence
// checks query by full path, so the key must match: with the path unindexed,
// each check degenerates into a scan of every row for the connection, and at a
// few million rows the extractor simply never finishes. The failure mode is "the
// task doesn't complete" rather than an error, which makes it hard to diagnose.
type KiroS3FileMeta struct {
	common.NoPKModel
	ConnectionId uint64 `gorm:"primaryKey"`
	// S3Path is the full object key. 512 chars because real keys already run
	// to ~155 and the prefix is user-configurable.
	S3Path string `gorm:"primaryKey;type:varchar(512)"`
	// FileName is the basename, for display and log messages only.
	FileName string `gorm:"type:varchar(255)" json:"fileName"`
	// Bucket records which bucket the object came from, since reports and logs
	// may live in different ones.
	Bucket  string `gorm:"type:varchar(255)" json:"bucket"`
	ScopeId string `gorm:"type:varchar(255);index" json:"scopeId"`
	// FileType is one of the FileType* constants above.
	FileType      string     `gorm:"type:varchar(32);index" json:"fileType"`
	Processed     bool       `gorm:"default:false" json:"processed"`
	ProcessedTime *time.Time `gorm:"default:null" json:"processedTime"`
	// RecordCount is how many rows the file yielded. Together with
	// ErrorMessage this turns "why is this month short?" from guesswork into a
	// single query.
	RecordCount int `gorm:"default:0" json:"recordCount"`
	// ErrorMessage holds the last parse failure reason, if any.
	ErrorMessage string `gorm:"type:text" json:"errorMessage"`
	// AttemptCount guards against infinite retries; see MaxAttempts.
	AttemptCount int `gorm:"default:0" json:"attemptCount"`
}

func (KiroS3FileMeta) TableName() string {
	return "_tool_kiro_s3_file_meta"
}
