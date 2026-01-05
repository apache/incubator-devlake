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
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/q_dev/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

var _ plugin.SubTaskEntryPoint = ExtractQDevS3Data

// ExtractQDevS3Data 从S3下载CSV数据并解析
func ExtractQDevS3Data(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*QDevTaskData)
	db := taskCtx.GetDal()

	// 查询未处理的文件元数据
	cursor, err := db.Cursor(
		dal.From(&models.QDevS3FileMeta{}),
		dal.Where("connection_id = ? AND processed = ?", data.Options.ConnectionId, false),
	)
	if err != nil {
		return errors.Default.Wrap(err, "failed to get file metadata cursor")
	}
	defer cursor.Close()

	taskCtx.SetProgress(0, -1)

	// 处理每个文件
	for cursor.Next() {
		fileMeta := &models.QDevS3FileMeta{}
		err = db.Fetch(cursor, fileMeta)
		if err != nil {
			return errors.Default.Wrap(err, "failed to fetch file metadata")
		}

		// 获取文件内容
		getInput := &s3.GetObjectInput{
			Bucket: aws.String(data.S3Client.Bucket),
			Key:    aws.String(fileMeta.S3Path),
		}

		getResult, err := data.S3Client.S3.GetObject(getInput)
		if err != nil {
			return errors.Convert(err)
		}

		// Use a transaction to process the file and update its status
		tx := db.Begin()
		csvErr := processCSVData(taskCtx, tx, getResult.Body, fileMeta)
		if csvErr != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				taskCtx.GetLogger().Error(rollbackErr, "failed to rollback transaction")
			}
			return errors.Default.Wrap(csvErr, fmt.Sprintf("failed to process CSV file %s", fileMeta.FileName))
		}

		// Update file processing status within the same transaction
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

		// Commit the transaction
		err = tx.Commit()
		if err != nil {
			return errors.Default.Wrap(err, "failed to commit transaction")
		}

		taskCtx.IncProgress(1)
	}

	return nil
}

// 处理CSV文件
func processCSVData(taskCtx plugin.SubTaskContext, db dal.Dal, reader io.ReadCloser, fileMeta *models.QDevS3FileMeta) errors.Error {
	defer reader.Close()

	// Get task data to access Identity Client
	data := taskCtx.GetData().(*QDevTaskData)

	csvReader := csv.NewReader(reader)
	// 使用默认的逗号分隔符，不需要设置 Comma
	csvReader.LazyQuotes = true    // 允许非标准引号处理
	csvReader.FieldsPerRecord = -1 // 允许每行字段数不同

	// 读取标头
	headers, err := csvReader.Read()
	taskCtx.GetLogger().Debug("CSV headers: %+v", headers)
	if err != nil {
		return errors.Convert(err)
	}

	// 逐行读取数据
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Convert(err)
		}

		// 创建用户数据对象 (updated to include display name resolution)
		userData, err := createUserDataWithDisplayName(taskCtx.GetLogger(), headers, record, fileMeta, data.IdentityClient)
		if err != nil {
			return errors.Default.Wrap(err, "failed to create user data")
		}

		// Save to database - no need to check for duplicates since we're processing each file only once
		err = db.Create(userData)
		if err != nil {
			return errors.Default.Wrap(err, "failed to save user data")
		}
	}

	return nil
}

// UserDisplayNameResolver interface for resolving user display names
type UserDisplayNameResolver interface {
	ResolveUserDisplayName(userId string) (string, error)
}

// 从CSV记录创建用户数据对象 (enhanced with display name resolution)
func createUserDataWithDisplayName(logger interface {
	Debug(format string, a ...interface{})
}, headers []string, record []string, fileMeta *models.QDevS3FileMeta, identityClient UserDisplayNameResolver) (*models.QDevUserData, errors.Error) {
	userData := &models.QDevUserData{
		ConnectionId: fileMeta.ConnectionId,
		ScopeId:      fileMeta.ScopeId,
	}

	// 创建字段映射
	fieldMap := make(map[string]string)
	for i, header := range headers {
		if i < len(record) {
			logger.Debug("Mapping header[%d]: '%s' -> '%s'", i, header, record[i])
			fieldMap[header] = record[i]
			// 同时添加去除空格的版本
			trimmedHeader := strings.TrimSpace(header)
			if trimmedHeader != header {
				logger.Debug("Also adding trimmed header: '%s'", trimmedHeader)
				fieldMap[trimmedHeader] = record[i]
			}
		}
	}

	// 设置必要字段
	var err error
	var ok bool

	// 设置UserId
	userData.UserId, ok = fieldMap["UserId"]
	if !ok {
		return nil, errors.Default.New("UserId not found in CSV record")
	}

	// 设置DisplayName (new functionality)
	userData.DisplayName = resolveDisplayName(logger, userData.UserId, identityClient)

	// 设置Date
	dateStr, ok := fieldMap["Date"]
	if !ok {
		return nil, errors.Default.New("Date not found in CSV record")
	}

	userData.Date, err = parseDate(dateStr)
	if err != nil {
		return nil, errors.Default.Wrap(err, "failed to parse date")
	}

	// 设置所有指标字段
	userData.CodeReview_FindingsCount = parseInt(fieldMap, "CodeReview_FindingsCount")
	userData.CodeReview_SucceededEventCount = parseInt(fieldMap, "CodeReview_SucceededEventCount")
	userData.InlineChat_AcceptanceEventCount = parseInt(fieldMap, "InlineChat_AcceptanceEventCount")
	userData.InlineChat_AcceptedLineAdditions = parseInt(fieldMap, "InlineChat_AcceptedLineAdditions")
	userData.InlineChat_AcceptedLineDeletions = parseInt(fieldMap, "InlineChat_AcceptedLineDeletions")
	userData.InlineChat_DismissalEventCount = parseInt(fieldMap, "InlineChat_DismissalEventCount")
	userData.InlineChat_DismissedLineAdditions = parseInt(fieldMap, "InlineChat_DismissedLineAdditions")
	userData.InlineChat_DismissedLineDeletions = parseInt(fieldMap, "InlineChat_DismissedLineDeletions")
	userData.InlineChat_RejectedLineAdditions = parseInt(fieldMap, "InlineChat_RejectedLineAdditions")
	userData.InlineChat_RejectedLineDeletions = parseInt(fieldMap, "InlineChat_RejectedLineDeletions")
	userData.InlineChat_RejectionEventCount = parseInt(fieldMap, "InlineChat_RejectionEventCount")
	userData.InlineChat_TotalEventCount = parseInt(fieldMap, "InlineChat_TotalEventCount")
	userData.Inline_AICodeLines = parseInt(fieldMap, "Inline_AICodeLines")
	userData.Inline_AcceptanceCount = parseInt(fieldMap, "Inline_AcceptanceCount")
	userData.Inline_SuggestionsCount = parseInt(fieldMap, "Inline_SuggestionsCount")
	userData.Chat_AICodeLines = parseInt(fieldMap, "Chat_AICodeLines")
	userData.Chat_MessagesInteracted = parseInt(fieldMap, "Chat_MessagesInteracted")
	userData.Chat_MessagesSent = parseInt(fieldMap, "Chat_MessagesSent")
	userData.CodeFix_AcceptanceEventCount = parseInt(fieldMap, "CodeFix_AcceptanceEventCount")
	userData.CodeFix_AcceptedLines = parseInt(fieldMap, "CodeFix_AcceptedLines")
	userData.CodeFix_GeneratedLines = parseInt(fieldMap, "CodeFix_GeneratedLines")
	userData.CodeFix_GenerationEventCount = parseInt(fieldMap, "CodeFix_GenerationEventCount")
	userData.CodeReview_FailedEventCount = parseInt(fieldMap, "CodeReview_FailedEventCount")
	userData.Dev_AcceptanceEventCount = parseInt(fieldMap, "Dev_AcceptanceEventCount")
	userData.Dev_AcceptedLines = parseInt(fieldMap, "Dev_AcceptedLines")
	userData.Dev_GeneratedLines = parseInt(fieldMap, "Dev_GeneratedLines")
	userData.Dev_GenerationEventCount = parseInt(fieldMap, "Dev_GenerationEventCount")
	userData.DocGeneration_AcceptedFileUpdates = parseInt(fieldMap, "DocGeneration_AcceptedFileUpdates")
	userData.DocGeneration_AcceptedFilesCreations = parseInt(fieldMap, "DocGeneration_AcceptedFilesCreations")
	userData.DocGeneration_AcceptedLineAdditions = parseInt(fieldMap, "DocGeneration_AcceptedLineAdditions")
	userData.DocGeneration_AcceptedLineUpdates = parseInt(fieldMap, "DocGeneration_AcceptedLineUpdates")
	userData.DocGeneration_EventCount = parseInt(fieldMap, "DocGeneration_EventCount")
	userData.DocGeneration_RejectedFileCreations = parseInt(fieldMap, "DocGeneration_RejectedFileCreations")
	userData.DocGeneration_RejectedFileUpdates = parseInt(fieldMap, "DocGeneration_RejectedFileUpdates")
	userData.DocGeneration_RejectedLineAdditions = parseInt(fieldMap, "DocGeneration_RejectedLineAdditions")
	userData.DocGeneration_RejectedLineUpdates = parseInt(fieldMap, "DocGeneration_RejectedLineUpdates")
	userData.TestGeneration_AcceptedLines = parseInt(fieldMap, "TestGeneration_AcceptedLines")
	userData.TestGeneration_AcceptedTests = parseInt(fieldMap, "TestGeneration_AcceptedTests")
	userData.TestGeneration_EventCount = parseInt(fieldMap, "TestGeneration_EventCount")
	userData.TestGeneration_GeneratedLines = parseInt(fieldMap, "TestGeneration_GeneratedLines")
	userData.TestGeneration_GeneratedTests = parseInt(fieldMap, "TestGeneration_GeneratedTests")
	userData.Transformation_EventCount = parseInt(fieldMap, "Transformation_EventCount")
	userData.Transformation_LinesGenerated = parseInt(fieldMap, "Transformation_LinesGenerated")
	userData.Transformation_LinesIngested = parseInt(fieldMap, "Transformation_LinesIngested")

	return userData, nil
}

// resolveDisplayName resolves user ID to display name using Identity Client
func resolveDisplayName(logger interface {
	Debug(format string, a ...interface{})
}, userId string, identityClient UserDisplayNameResolver) string {
	// If no identity client available, use userId as fallback
	if identityClient == nil {
		return userId
	}

	// Try to resolve display name
	displayName, err := identityClient.ResolveUserDisplayName(userId)
	if err != nil {
		// Log error but continue with userId as fallback
		logger.Debug("Failed to resolve display name for user %s: %v", userId, err)
		return userId
	}

	// If display name is empty, use userId as fallback
	if displayName == "" {
		return userId
	}

	return displayName
}

// 解析日期
func parseDate(dateStr string) (time.Time, errors.Error) {
	// 尝试常见的日期格式
	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"01/02/2006",
		"01-02-2006",
		time.RFC3339,
	}

	for _, format := range formats {
		date, err := time.Parse(format, dateStr)
		if err == nil {
			return date, nil
		}
	}

	return time.Time{}, errors.Default.New(fmt.Sprintf("failed to parse date: %s", dateStr))
}

// 解析整数
func parseInt(fieldMap map[string]string, field string) int {
	value, ok := fieldMap[field]
	if !ok {
		return 0
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return intValue
}

var ExtractQDevS3DataMeta = plugin.SubTaskMeta{
	Name:             "extractQDevS3Data",
	EntryPoint:       ExtractQDevS3Data,
	EnabledByDefault: true,
	Description:      "Extract data from S3 CSV files and save to database",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
	Dependencies:     []*plugin.SubTaskMeta{&CollectQDevS3FilesMeta},
}
