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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/q_dev/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"strings"
)

var _ plugin.SubTaskEntryPoint = CollectQDevS3Files

// CollectQDevS3Files 收集S3文件元数据
func CollectQDevS3Files(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*QDevTaskData)
	db := taskCtx.GetDal()

	// 列出指定前缀下的所有对象
	var continuationToken *string
	prefix := data.Options.S3Prefix
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix = prefix + "/"
	}

	taskCtx.SetProgress(0, -1)

	// 清空以前的元数据记录
	err := db.Delete(&models.QDevS3FileMeta{}, dal.Where("connection_id = ?", data.Options.ConnectionId))
	if err != nil {
		return errors.Default.Wrap(err, "failed to clean previous file metadata")
	}

	for {
		input := &s3.ListObjectsV2Input{
			Bucket:            aws.String(data.S3Client.Bucket),
			Prefix:            aws.String(prefix),
			ContinuationToken: continuationToken,
		}

		result, err := data.S3Client.S3.ListObjectsV2(input)
		if err != nil {
			return errors.Convert(err)
		}

		// 处理每个CSV文件
		for _, object := range result.Contents {
			// 只处理CSV文件
			if !strings.HasSuffix(*object.Key, ".csv") {
				continue
			}

			// 保存文件元数据
			fileMeta := &models.QDevS3FileMeta{
				ConnectionId: data.Options.ConnectionId,
				FileName:     *object.Key,
				S3Path:       *object.Key,
				Processed:    false,
			}

			err = db.Create(fileMeta)
			if err != nil {
				return errors.Default.Wrap(err, "failed to create file metadata")
			}

			taskCtx.IncProgress(1)
		}

		// 如果没有更多对象，退出循环
		if !*result.IsTruncated {
			break
		}

		continuationToken = result.NextContinuationToken
	}

	return nil
}

var CollectQDevS3FilesMeta = plugin.SubTaskMeta{
	Name:             "collectQDevS3Files",
	EntryPoint:       CollectQDevS3Files,
	EnabledByDefault: true,
	Description:      "Collect S3 file metadata from AWS S3 bucket",
}
