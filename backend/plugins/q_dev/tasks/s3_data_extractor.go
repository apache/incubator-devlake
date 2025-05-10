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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/q_dev/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"strconv"
	"strings"
	"time"
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

		// 处理CSV文件
		err = processCSVData(taskCtx, db, getResult.Body, fileMeta)
		if err != nil {
			return errors.Default.Wrap(err, fmt.Sprintf("failed to process CSV file %s", fileMeta.FileName))
		}

		// 更新文件处理状态
		fileMeta.Processed = true
		now := time.Now()
		fileMeta.ProcessedTime = &now
		err = db.Update(fileMeta)
		if err != nil {
			return errors.Default.Wrap(err, "failed to update file metadata")
		}

		taskCtx.IncProgress(1)
	}

	return nil
}

// 处理CSV文件
func processCSVData(taskCtx plugin.SubTaskContext, db dal.Dal, reader io.ReadCloser, fileMeta *models.QDevS3FileMeta) errors.Error {
	defer reader.Close()

	csvReader := csv.NewReader(reader)
	// 使用默认的逗号分隔符，不需要设置 Comma
	csvReader.LazyQuotes = true    // 允许非标准引号处理
	csvReader.FieldsPerRecord = -1 // 允许每行字段数不同

	// 读取标头
	headers, err := csvReader.Read()
	fmt.Printf("headers: %+v\n", headers)
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

		// 创建用户数据对象
		userData, err := createUserData(headers, record, fileMeta)
		if err != nil {
			return errors.Default.Wrap(err, "failed to create user data")
		}

		// 保存到数据库
		err = db.Create(userData)
		if err != nil {
			return errors.Default.Wrap(err, "failed to save user data")
		}
	}

	return nil
}

// 从CSV记录创建用户数据对象
func createUserData(headers []string, record []string, fileMeta *models.QDevS3FileMeta) (*models.QDevUserData, errors.Error) {
	userData := &models.QDevUserData{
		ConnectionId: fileMeta.ConnectionId,
	}

	// 创建字段映射
	fieldMap := make(map[string]string)
	for i, header := range headers {
		if i < len(record) {
			// 打印每个header和对应的值，帮助调试
			fmt.Printf("Mapping header[%d]: '%s' -> '%s'\n", i, header, record[i])
			fieldMap[header] = record[i]
			// 同时添加去除空格的版本
			trimmedHeader := strings.TrimSpace(header)
			if trimmedHeader != header {
				fmt.Printf("Also adding trimmed header: '%s'\n", trimmedHeader)
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

	// 设置Date
	dateStr, ok := fieldMap["Date"]
	if !ok {
		return nil, errors.Default.New("Date not found in CSV record")
	}

	userData.Date, err = parseDate(dateStr)
	if err != nil {
		return nil, errors.Default.Wrap(err, "failed to parse date")
	}

	// 设置指标字段
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

	return userData, nil
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
}
