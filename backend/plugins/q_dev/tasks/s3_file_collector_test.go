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
	"errors"
	"testing"

	"github.com/apache/incubator-devlake/plugins/q_dev/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeS3Prefix(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"prefix", "prefix/"},
		{"prefix/", "prefix/"},
		{"path/to/folder", "path/to/folder/"},
		{"path/to/folder/", "path/to/folder/"},
	}

	for _, test := range tests {
		result := normalizeS3Prefix(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestIsCSVFile(t *testing.T) {
	tests := []struct {
		key      string
		expected bool
	}{
		{"file.csv", true},
		{"data.CSV", false},
		{"report.csv", true},
		{"document.txt", false},
		{"path/to/file.csv", true},
		{"file.csv.backup", false},
		{"", false},
	}

	for _, test := range tests {
		result := isCSVFile(test.key)
		assert.Equal(t, test.expected, result)
	}
}

func TestCollectS3FilesCore_Success(t *testing.T) {
	// Mock functions
	listObjects := func(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
		return &s3.ListObjectsV2Output{
			Contents: []*s3.Object{
				{Key: aws.String("file1.csv")},
				{Key: aws.String("file2.txt")},
				{Key: aws.String("data.csv")},
			},
			IsTruncated: aws.Bool(false),
		}, nil
	}

	findFileMeta := func(connectionId uint64, s3Path string) (*models.QDevS3FileMeta, error) {
		return nil, errors.New("not found")
	}

	createdFiles := []string{}
	saveFileMeta := func(fileMeta *models.QDevS3FileMeta) error {
		createdFiles = append(createdFiles, fileMeta.S3Path)
		return nil
	}

	progressCount := 0
	progress := func(increment int) {
		progressCount += increment
	}

	logMessages := []string{}
	logDebug := func(format string, args ...interface{}) {
		logMessages = append(logMessages, format)
	}

	err := collectS3FilesCore("bucket", "prefix", 1, listObjects, findFileMeta, saveFileMeta, progress, logDebug)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(createdFiles))
	assert.Contains(t, createdFiles, "file1.csv")
	assert.Contains(t, createdFiles, "data.csv")
	assert.Equal(t, 2, progressCount)
	assert.Contains(t, logMessages, "Skipping non-CSV file: %s")
}

func TestCollectS3FilesCore_NoCSVFiles(t *testing.T) {
	listObjects := func(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
		return &s3.ListObjectsV2Output{
			Contents: []*s3.Object{
				{Key: aws.String("file1.txt")},
				{Key: aws.String("file2.json")},
			},
			IsTruncated: aws.Bool(false),
		}, nil
	}

	findFileMeta := func(connectionId uint64, s3Path string) (*models.QDevS3FileMeta, error) {
		return nil, errors.New("not found")
	}

	saveFileMeta := func(fileMeta *models.QDevS3FileMeta) error {
		return nil
	}

	progress := func(increment int) {}
	logDebug := func(format string, args ...interface{}) {}

	err := collectS3FilesCore("bucket", "prefix", 1, listObjects, findFileMeta, saveFileMeta, progress, logDebug)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no CSV files found")
}

func TestCollectS3FilesCore_SkipProcessedFiles(t *testing.T) {
	listObjects := func(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
		return &s3.ListObjectsV2Output{
			Contents: []*s3.Object{
				{Key: aws.String("processed.csv")},
				{Key: aws.String("unprocessed.csv")},
			},
			IsTruncated: aws.Bool(false),
		}, nil
	}

	findFileMeta := func(connectionId uint64, s3Path string) (*models.QDevS3FileMeta, error) {
		if s3Path == "processed.csv" {
			return &models.QDevS3FileMeta{Processed: true}, nil
		}
		if s3Path == "unprocessed.csv" {
			return &models.QDevS3FileMeta{Processed: false}, nil
		}
		return nil, errors.New("not found")
	}

	createdFiles := []string{}
	saveFileMeta := func(fileMeta *models.QDevS3FileMeta) error {
		createdFiles = append(createdFiles, fileMeta.S3Path)
		return nil
	}

	progress := func(increment int) {}
	logDebug := func(format string, args ...interface{}) {}

	err := collectS3FilesCore("bucket", "prefix", 1, listObjects, findFileMeta, saveFileMeta, progress, logDebug)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(createdFiles)) // No new files should be created
}
