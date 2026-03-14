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

package tasks

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/q_dev/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

var _ plugin.SubTaskEntryPoint = ExtractQDevLoggingData

// ExtractQDevLoggingData extracts logging data from S3 JSON.gz files
func ExtractQDevLoggingData(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*QDevTaskData)
	db := taskCtx.GetDal()

	cursor, err := db.Cursor(
		dal.From(&models.QDevS3FileMeta{}),
		dal.Where("connection_id = ? AND processed = ? AND file_name LIKE ?",
			data.Options.ConnectionId, false, "%.json.gz"),
	)
	if err != nil {
		return errors.Default.Wrap(err, "failed to get logging file metadata cursor")
	}
	defer cursor.Close()

	taskCtx.SetProgress(0, -1)

	for cursor.Next() {
		fileMeta := &models.QDevS3FileMeta{}
		err = db.Fetch(cursor, fileMeta)
		if err != nil {
			return errors.Default.Wrap(err, "failed to fetch file metadata")
		}

		getInput := &s3.GetObjectInput{
			Bucket: aws.String(data.S3Client.Bucket),
			Key:    aws.String(fileMeta.S3Path),
		}

		getResult, err := data.S3Client.S3.GetObject(getInput)
		if err != nil {
			return errors.Convert(err)
		}

		tx := db.Begin()
		processErr := processLoggingData(taskCtx, tx, getResult.Body, fileMeta)
		if processErr != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				taskCtx.GetLogger().Error(rollbackErr, "failed to rollback transaction")
			}
			return errors.Default.Wrap(processErr, fmt.Sprintf("failed to process logging file %s", fileMeta.FileName))
		}

		fileMeta.Processed = true
		now := time.Now()
		fileMeta.ProcessedTime = &now
		err = tx.Update(fileMeta)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				taskCtx.GetLogger().Error(rollbackErr, "failed to rollback transaction")
			}
			return errors.Default.Wrap(err, "failed to update file metadata")
		}

		err = tx.Commit()
		if err != nil {
			return errors.Default.Wrap(err, "failed to commit transaction")
		}

		taskCtx.IncProgress(1)
	}

	return nil
}

// JSON structures for logging data

type loggingFile struct {
	Records []json.RawMessage `json:"records"`
}

type chatLogRecord struct {
	Request  *chatLogRequest  `json:"generateAssistantResponseEventRequest"`
	Response *chatLogResponse `json:"generateAssistantResponseEventResponse"`
}

type chatLogRequest struct {
	UserID           string  `json:"userId"`
	Timestamp        string  `json:"timeStamp"`
	ChatTriggerType  string  `json:"chatTriggerType"`
	CustomizationArn *string `json:"customizationArn"`
	ModelID          string  `json:"modelId"`
	Prompt           string  `json:"prompt"`
}

type chatLogResponse struct {
	RequestID         string `json:"requestId"`
	AssistantResponse string `json:"assistantResponse"`
	MessageMetadata   struct {
		ConversationID *string `json:"conversationId"`
		UtteranceID    *string `json:"utteranceId"`
	} `json:"messageMetadata"`
}

type completionLogRecord struct {
	Request  *completionLogRequest  `json:"generateCompletionsEventRequest"`
	Response *completionLogResponse `json:"generateCompletionsEventResponse"`
}

type completionLogRequest struct {
	UserID           string  `json:"userId"`
	Timestamp        string  `json:"timeStamp"`
	FileName         string  `json:"fileName"`
	CustomizationArn *string `json:"customizationArn"`
}

type completionLogResponse struct {
	RequestID   string            `json:"requestId"`
	Completions []json.RawMessage `json:"completions"`
}

func processLoggingData(taskCtx plugin.SubTaskContext, db dal.Dal, reader io.ReadCloser, fileMeta *models.QDevS3FileMeta) errors.Error {
	defer reader.Close()

	data := taskCtx.GetData().(*QDevTaskData)

	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return errors.Convert(err)
	}
	defer gzReader.Close()

	var logFile loggingFile
	decoder := json.NewDecoder(gzReader)
	if err := decoder.Decode(&logFile); err != nil {
		return errors.Convert(err)
	}

	isChatLog := strings.Contains(fileMeta.S3Path, "GenerateAssistantResponse")

	for _, rawRecord := range logFile.Records {
		if isChatLog {
			if err := processChatRecord(taskCtx, db, rawRecord, fileMeta, data.IdentityClient); err != nil {
				return err
			}
		} else {
			if err := processCompletionRecord(taskCtx, db, rawRecord, fileMeta, data.IdentityClient); err != nil {
				return err
			}
		}
	}

	return nil
}

func processChatRecord(taskCtx plugin.SubTaskContext, db dal.Dal, raw json.RawMessage, fileMeta *models.QDevS3FileMeta, identityClient UserDisplayNameResolver) errors.Error {
	var record chatLogRecord
	if err := json.Unmarshal(raw, &record); err != nil {
		return errors.Convert(err)
	}

	if record.Request == nil || record.Response == nil {
		return nil
	}

	ts, err := time.Parse(time.RFC3339Nano, record.Request.Timestamp)
	if err != nil {
		ts = time.Now()
	}

	chatLog := &models.QDevChatLog{
		ConnectionId:     fileMeta.ConnectionId,
		ScopeId:          fileMeta.ScopeId,
		RequestId:        record.Response.RequestID,
		UserId:           record.Request.UserID,
		DisplayName:      resolveDisplayName(taskCtx.GetLogger(), record.Request.UserID, identityClient),
		Timestamp:        ts,
		ChatTriggerType:  record.Request.ChatTriggerType,
		HasCustomization: record.Request.CustomizationArn != nil && *record.Request.CustomizationArn != "",
		ModelId:          record.Request.ModelID,
		PromptLength:     len(record.Request.Prompt),
		ResponseLength:   len(record.Response.AssistantResponse),
	}

	// Parse structured info from prompt
	prompt := record.Request.Prompt
	chatLog.OpenFileCount = countOpenFiles(prompt)
	chatLog.ActiveFileName, chatLog.ActiveFileExtension = parseActiveFile(prompt)
	chatLog.HasSteering = strings.Contains(prompt, ".kiro/steering")
	chatLog.IsSpecMode = strings.Contains(prompt, "implicit-rules")

	if record.Response.MessageMetadata.ConversationID != nil {
		chatLog.ConversationId = *record.Response.MessageMetadata.ConversationID
	}
	if record.Response.MessageMetadata.UtteranceID != nil {
		chatLog.UtteranceId = *record.Response.MessageMetadata.UtteranceID
	}

	return errors.Default.Wrap(db.CreateOrUpdate(chatLog), "failed to save chat log")
}

// countOpenFiles counts <file name="..."> tags within <OPEN-EDITOR-FILES> block
func countOpenFiles(prompt string) int {
	start := strings.Index(prompt, "<OPEN-EDITOR-FILES>")
	if start == -1 {
		return 0
	}
	end := strings.Index(prompt, "</OPEN-EDITOR-FILES>")
	if end == -1 {
		return 0
	}
	block := prompt[start:end]
	return strings.Count(block, "<file name=")
}

// parseActiveFile extracts the active file name and extension from prompt
func parseActiveFile(prompt string) (string, string) {
	start := strings.Index(prompt, "<ACTIVE-EDITOR-FILE>")
	if start == -1 {
		return "", ""
	}
	end := strings.Index(prompt[start:], "</ACTIVE-EDITOR-FILE>")
	if end == -1 {
		return "", ""
	}
	block := prompt[start : start+end]
	// Find <file name="..." />
	nameStart := strings.Index(block, "name=\"")
	if nameStart == -1 {
		return "", ""
	}
	nameStart += len("name=\"")
	nameEnd := strings.Index(block[nameStart:], "\"")
	if nameEnd == -1 {
		return "", ""
	}
	fileName := block[nameStart : nameStart+nameEnd]
	ext := filepath.Ext(fileName)
	return fileName, ext
}

func processCompletionRecord(taskCtx plugin.SubTaskContext, db dal.Dal, raw json.RawMessage, fileMeta *models.QDevS3FileMeta, identityClient UserDisplayNameResolver) errors.Error {
	var record completionLogRecord
	if err := json.Unmarshal(raw, &record); err != nil {
		return errors.Convert(err)
	}

	if record.Request == nil || record.Response == nil {
		return nil
	}

	ts, err := time.Parse(time.RFC3339Nano, record.Request.Timestamp)
	if err != nil {
		ts = time.Now()
	}

	completionLog := &models.QDevCompletionLog{
		ConnectionId:     fileMeta.ConnectionId,
		ScopeId:          fileMeta.ScopeId,
		RequestId:        record.Response.RequestID,
		UserId:           record.Request.UserID,
		DisplayName:      resolveDisplayName(taskCtx.GetLogger(), record.Request.UserID, identityClient),
		Timestamp:        ts,
		FileName:         record.Request.FileName,
		FileExtension:    filepath.Ext(record.Request.FileName),
		HasCustomization: record.Request.CustomizationArn != nil && *record.Request.CustomizationArn != "",
		CompletionsCount: len(record.Response.Completions),
	}

	return errors.Default.Wrap(db.CreateOrUpdate(completionLog), "failed to save completion log")
}

var ExtractQDevLoggingDataMeta = plugin.SubTaskMeta{
	Name:             "extractQDevLoggingData",
	EntryPoint:       ExtractQDevLoggingData,
	EnabledByDefault: true,
	Description:      "Extract logging data from S3 JSON.gz files (chat and completion events)",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
	Dependencies:     []*plugin.SubTaskMeta{&CollectQDevS3FilesMeta},
}
