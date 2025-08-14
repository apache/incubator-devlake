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
	"strings"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/q_dev/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// normalizeS3Prefix ensures the prefix ends with "/" if it's not empty
func normalizeS3Prefix(prefix string) string {
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		return prefix + "/"
	}
	return prefix
}

// isCSVFile checks if the given S3 object key represents a CSV file
func isCSVFile(key string) bool {
	return strings.HasSuffix(key, ".csv")
}

// S3ListObjectsFunc defines the callback for listing S3 objects
type S3ListObjectsFunc func(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error)

// FindFileMetaFunc defines the callback for finding existing file metadata
type FindFileMetaFunc func(connectionId uint64, s3Path string) (*models.QDevS3FileMeta, error)

// SaveFileMetaFunc defines the callback for saving file metadata
type SaveFileMetaFunc func(fileMeta *models.QDevS3FileMeta) error

// ProgressFunc defines the callback for progress tracking
type ProgressFunc func(increment int)

// LogFunc defines the callback for logging
type LogFunc func(format string, args ...interface{})

// collectS3FilesCore contains the core logic for collecting S3 files
func collectS3FilesCore(
	bucket, prefix string,
	connectionId uint64,
	listObjects S3ListObjectsFunc,
	findFileMeta FindFileMetaFunc,
	saveFileMeta SaveFileMetaFunc,
	progress ProgressFunc,
	logDebug LogFunc,
) error {
	// List all objects under the specified prefix
	var continuationToken *string
	normalizedPrefix := normalizeS3Prefix(prefix)
	csvFilesFound := 0

	for {
		input := &s3.ListObjectsV2Input{
			Bucket:            aws.String(bucket),
			Prefix:            aws.String(normalizedPrefix),
			ContinuationToken: continuationToken,
		}

		result, err := listObjects(input)
		if err != nil {
			return err
		}

		// Process each CSV file
		for _, object := range result.Contents {
			// Only process CSV files
			if !isCSVFile(*object.Key) {
				logDebug("Skipping non-CSV file: %s", *object.Key)
				continue
			}

			csvFilesFound++

			// Check if this file already exists in our database
			existingFile, err := findFileMeta(connectionId, *object.Key)
			if err == nil {
				// File already exists in database, skip it if it's already processed
				if existingFile.Processed {
					logDebug("Skipping already processed file: %s", *object.Key)
					continue
				}
				// Otherwise, we'll keep the existing record (which is still marked as unprocessed)
				logDebug("Found existing unprocessed file: %s", *object.Key)
				continue
			}

			// This is a new file, save its metadata
			fileMeta := &models.QDevS3FileMeta{
				ConnectionId: connectionId,
				FileName:     *object.Key,
				S3Path:       *object.Key,
				Processed:    false,
			}

			if err := saveFileMeta(fileMeta); err != nil {
				return err
			}

			progress(1)
		}

		// If there are no more objects, exit the loop
		if !*result.IsTruncated {
			break
		}

		continuationToken = result.NextContinuationToken
	}

	// Check if no CSV files were found
	if csvFilesFound == 0 {
		return errors.BadInput.New("no CSV files found in S3 path. Please verify the S3 bucket and prefix configuration")
	}

	return nil
}



var _ plugin.SubTaskEntryPoint = CollectQDevS3Files

// CollectQDevS3Files 收集S3文件元数据
func CollectQDevS3Files(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*QDevTaskData)
	db := taskCtx.GetDal()

	taskCtx.SetProgress(0, -1)

	// Define callback functions
	listObjects := func(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
		return data.S3Client.S3.ListObjectsV2(input)
	}

	findFileMeta := func(connectionId uint64, s3Path string) (*models.QDevS3FileMeta, error) {
		existingFile := &models.QDevS3FileMeta{}
		err := db.First(existingFile, dal.Where("connection_id = ? AND s3_path = ?", connectionId, s3Path))
		if err != nil {
			if db.IsErrorNotFound(err) {
				return nil, err
			}
			return nil, errors.Default.Wrap(err, "failed to query existing file metadata")
		}
		return existingFile, nil
	}

	saveFileMeta := func(fileMeta *models.QDevS3FileMeta) error {
		err := db.Create(fileMeta)
		if err != nil {
			return errors.Default.Wrap(err, "failed to create file metadata")
		}
		return nil
	}

	progress := func(increment int) {
		taskCtx.IncProgress(increment)
	}

	logDebug := func(format string, args ...interface{}) {
		taskCtx.GetLogger().Debug(format, args...)
	}

	// Call the core function
	err := collectS3FilesCore(
		data.S3Client.Bucket,
		data.Options.S3Prefix,
		data.Options.ConnectionId,
		listObjects,
		findFileMeta,
		saveFileMeta,
		progress,
		logDebug,
	)

	return errors.Convert(err)
}

var CollectQDevS3FilesMeta = plugin.SubTaskMeta{
	Name:             "collectQDevS3Files",
	EntryPoint:       CollectQDevS3Files,
	EnabledByDefault: true,
	Description:      "Collect S3 file metadata from AWS S3 bucket",
}
