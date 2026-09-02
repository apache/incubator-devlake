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
	"os"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/devlake/plugins/kiro/models"
)

func listOutput(truncated bool, keys ...string) *s3.ListObjectsV2Output {
	contents := make([]*s3.Object, 0, len(keys))
	for _, k := range keys {
		key := k
		contents = append(contents, &s3.Object{Key: &key})
	}
	return &s3.ListObjectsV2Output{
		Contents:    contents,
		IsTruncated: aws.Bool(truncated),
	}
}

func TestCollectCandidates_FileFiltering(t *testing.T) {
	spec := PrefixSpec{FileType: models.FileTypeChatLog}
	options := &KiroOptions{ConnectionId: 1, ScopeId: "s1"}

	t.Run("keeps csv and json.gz only", func(t *testing.T) {
		output := listOutput(false,
			"p/report.csv",
			"p/log.json.gz",
			// AWS writes small extension-less objects at the KiroLogs root as
			// permission probes; they must not enter the work queue.
			"p/26404955-bf00-40d3-b713-43d18edf0638",
			"p/notes.txt",
			"p/archive.zip",
		)
		candidates := collectCandidates(output, "bkt", spec, options)
		require.Len(t, candidates, 2)
		assert.Equal(t, "p/report.csv", candidates[0].S3Path)
		assert.Equal(t, "p/log.json.gz", candidates[1].S3Path)
	})

	t.Run("stores basename separately from the full key", func(t *testing.T) {
		key := "logging/AWSLogs/123456789012/KiroLogs/GenerateAssistantResponse/us-east-1/2026/07/27/23/" +
			"123456789012_GenerateAssistantResponse_202607272303_3tbIeIrGJNDFbfVx.json.gz"
		candidates := collectCandidates(listOutput(false, key), "bkt", spec, options)
		require.Len(t, candidates, 1)

		// The full key goes in S3Path, which is sized for it. FileName holds
		// only the basename - a full key there would eventually overflow.
		assert.Equal(t, key, candidates[0].S3Path)
		assert.Equal(t, "123456789012_GenerateAssistantResponse_202607272303_3tbIeIrGJNDFbfVx.json.gz",
			candidates[0].FileName)
		assert.Less(t, len(candidates[0].FileName), 255)
	})

	t.Run("records bucket, scope and file type", func(t *testing.T) {
		candidates := collectCandidates(listOutput(false, "p/a.csv"), "my-bucket", spec, options)
		require.Len(t, candidates, 1)
		assert.Equal(t, "my-bucket", candidates[0].Bucket)
		assert.Equal(t, "s1", candidates[0].ScopeId)
		assert.Equal(t, models.FileTypeChatLog, candidates[0].FileType)
		assert.Equal(t, uint64(1), candidates[0].ConnectionId)
		assert.False(t, candidates[0].Processed)
	})

	t.Run("nil key is skipped", func(t *testing.T) {
		output := &s3.ListObjectsV2Output{
			Contents:    []*s3.Object{{Key: nil}, {Key: aws.String("p/a.csv")}},
			IsTruncated: aws.Bool(false),
		}
		candidates := collectCandidates(output, "bkt", spec, options)
		assert.Len(t, candidates, 1)
	})

	t.Run("empty page yields nothing", func(t *testing.T) {
		candidates := collectCandidates(listOutput(false), "bkt", spec, options)
		assert.Empty(t, candidates)
	})
}

func TestBuildPrefixes(t *testing.T) {
	// Paths verified against real exports.
	t.Run("single bucket layout", func(t *testing.T) {
		conn := &models.KiroConnection{KiroConn: models.KiroConn{
			Region:          "us-east-1",
			Bucket:          "kiro-export-test",
			ReportPrefix:    "user-report",
			PromptLogPrefix: "logging",
		}}
		prefixes := BuildPrefixes(conn, "123456789012", "2026/07")

		require.Len(t, prefixes, 3)
		assert.Equal(t,
			"user-report/AWSLogs/123456789012/KiroLogs/user_report/us-east-1/2026/07",
			prefixes[0].Prefix)
		assert.Equal(t, models.FileTypeReport, prefixes[0].FileType)
		assert.Equal(t,
			"logging/AWSLogs/123456789012/KiroLogs/GenerateAssistantResponse/us-east-1/2026/07",
			prefixes[1].Prefix)
		assert.Equal(t, models.FileTypeChatLog, prefixes[1].FileType)
		assert.Equal(t,
			"logging/AWSLogs/123456789012/KiroLogs/GenerateCompletions/us-east-1/2026/07",
			prefixes[2].Prefix)
		assert.Equal(t, models.FileTypeCompletionLog, prefixes[2].FileType)

		for _, p := range prefixes {
			assert.Equal(t, "kiro-export-test", p.Bucket)
		}
	})

	t.Run("defaults apply when prefixes are unset", func(t *testing.T) {
		conn := &models.KiroConnection{KiroConn: models.KiroConn{
			Region: "us-east-1",
			Bucket: "b",
		}}
		prefixes := BuildPrefixes(conn, "acct", "2026/07")
		assert.Contains(t, prefixes[0].Prefix, "user-report/AWSLogs/acct/KiroLogs/user_report")
		assert.Contains(t, prefixes[1].Prefix, "logging/AWSLogs/acct/KiroLogs/GenerateAssistantResponse")
	})

	t.Run("separate buckets are assigned per file type", func(t *testing.T) {
		conn := &models.KiroConnection{KiroConn: models.KiroConn{
			Region:          "us-east-1",
			Bucket:          "reports",
			PromptLogBucket: "logs",
		}}
		prefixes := BuildPrefixes(conn, "acct", "2026/07")
		assert.Equal(t, "reports", prefixes[0].Bucket)
		assert.Equal(t, "logs", prefixes[1].Bucket)
		assert.Equal(t, "logs", prefixes[2].Bucket)
	})

	// A nil month widens the scope to the whole year, which is how a year-long
	// backfill is expressed.
	t.Run("year-only time path", func(t *testing.T) {
		conn := &models.KiroConnection{KiroConn: models.KiroConn{Region: "us-east-1", Bucket: "b"}}
		prefixes := BuildPrefixes(conn, "acct", "2026")
		assert.True(t, regexp.MustCompile(`/us-east-1/2026$`).MatchString(prefixes[0].Prefix))
	})
}

func TestWorkerCount(t *testing.T) {
	assert.Equal(t, DefaultWorkerCount, (&KiroTaskData{Options: &KiroOptions{}}).WorkerCount())
	assert.Equal(t, 5, (&KiroTaskData{Options: &KiroOptions{WorkerCount: 5}}).WorkerCount())
	// A zero or negative override is ignored rather than disabling concurrency.
	assert.Equal(t, DefaultWorkerCount, (&KiroTaskData{Options: &KiroOptions{WorkerCount: -1}}).WorkerCount())
	assert.Equal(t, DefaultWorkerCount, (&KiroTaskData{}).WorkerCount())
}

// This guards the predecessor defect that motivated the primary key choice: it queried
// its cursor table by s3_path while keying it on file_name, so the lookup falls
// back to scanning every row for the connection and collection never finishes at
// scale. The failure mode is a task that hangs rather than an error, so it is
// worth asserting structurally instead of hoping a reviewer notices.
func TestFileMetaQueriesMatchPrimaryKey(t *testing.T) {
	pkColumns := primaryKeyColumns(t, "s3_file_meta.go")
	require.Equal(t, []string{"ConnectionId", "S3Path"}, pkColumns,
		"the cursor table must be keyed on the connection and the full object path")

	source, err := os.ReadFile("s3_file_collector.go")
	require.NoError(t, err)

	whereClauses := regexp.MustCompile(`dal\.Where\(\s*"([^"]+)"`).FindAllStringSubmatch(string(source), -1)
	require.NotEmpty(t, whereClauses, "expected at least one filtered query")

	for _, clause := range whereClauses {
		condition := clause[1]
		assert.Contains(t, condition, "connection_id",
			"every cursor query must filter on connection_id, the first key column")
		assert.Contains(t, condition, "s3_path",
			"every cursor query must filter on s3_path, the second key column")
		assert.NotContains(t, condition, "file_name",
			"file_name is not indexed and must never appear in a lookup")
	}
}

// primaryKeyColumns extracts the fields tagged as primary keys from a model
// file, in declaration order - which is also the order of the composite index.
func primaryKeyColumns(t *testing.T, modelFile string) []string {
	t.Helper()
	source, err := os.ReadFile("../models/" + modelFile)
	require.NoError(t, err)

	fieldRe := regexp.MustCompile(`(?m)^\s*([A-Z][A-Za-z0-9]*)\s+\S+\s+` + "`" + `[^` + "`" + `]*primaryKey[^` + "`" + `]*` + "`")
	matches := fieldRe.FindAllStringSubmatch(string(source), -1)

	columns := make([]string, 0, len(matches))
	for _, m := range matches {
		columns = append(columns, m[1])
	}
	return columns
}
