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
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/incubator-devlake/plugins/kiro/models"
)

// prefixMockS3 answers listings from a canned prefix tree, and records the
// requests so tests can assert that a delimiter was used.
type prefixMockS3 struct {
	// commonPrefixes maps a queried prefix to the child prefixes returned.
	commonPrefixes map[string][]string
	// objects maps a queried prefix to the object keys beneath it.
	objects       map[string][]string
	seenDelimiter []string
	seenPrefix    []string
}

func (m *prefixMockS3) ListObjectsV2(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	prefix := ""
	if input.Prefix != nil {
		prefix = *input.Prefix
	}
	m.seenPrefix = append(m.seenPrefix, prefix)
	if input.Delimiter != nil {
		m.seenDelimiter = append(m.seenDelimiter, *input.Delimiter)
	} else {
		m.seenDelimiter = append(m.seenDelimiter, "")
	}

	out := &s3.ListObjectsV2Output{IsTruncated: aws.Bool(false)}
	for _, child := range m.commonPrefixes[prefix] {
		full := prefix + child + "/"
		out.CommonPrefixes = append(out.CommonPrefixes, &s3.CommonPrefix{Prefix: aws.String(full)})
	}
	for _, key := range m.objects[prefix] {
		out.Contents = append(out.Contents, &s3.Object{Key: aws.String(prefix + key)})
	}
	return out, nil
}

func (m *prefixMockS3) GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return nil, nil
}

func discoveryFixture(svc S3API) *Discovery {
	conn := &models.KiroConnection{KiroConn: models.KiroConn{
		Region:          "us-east-1",
		Bucket:          "kiro-export-test",
		ReportPrefix:    "user-report",
		PromptLogPrefix: "logging",
	}}
	client := &KiroS3Client{S3: svc, Bucket: conn.Bucket}
	return &Discovery{
		clients:    &KiroS3Clients{Report: client, PromptLog: client},
		connection: conn,
	}
}

// Kiro encodes every scope dimension as a path segment, which is what allows a
// scope to be selected rather than typed. This matters beyond convenience: a
// mistyped prefix and a month with no data produce the same outcome - a
// successful run that collects nothing - so a hand-entered path cannot be
// verified from the result.
func TestDiscovery_ListsLayoutFromS3(t *testing.T) {
	// The real layout, as confirmed against the live bucket.
	svc := &prefixMockS3{commonPrefixes: map[string][]string{
		"user-report/AWSLogs/": {"123456789012"},
		"user-report/AWSLogs/123456789012/KiroLogs/user_report/us-east-1/":      {"2026"},
		"user-report/AWSLogs/123456789012/KiroLogs/user_report/us-east-1/2026/": {"02", "03", "07"},
	}}
	discovery := discoveryFixture(svc)

	accounts, err := discovery.ListAccounts()
	require.Nil(t, err)
	assert.Equal(t, []string{"123456789012"}, accounts)

	years, err := discovery.ListYears("123456789012")
	require.Nil(t, err)
	assert.Equal(t, []int{2026}, years)

	months, err := discovery.ListMonths("123456789012", 2026)
	require.Nil(t, err)
	// Zero-padded segments must parse to plain integers, and only months that
	// actually hold data are offered.
	assert.Equal(t, []int{2, 3, 7}, months)

	// A delimiter is required: without it S3 returns every object under the
	// prefix instead of just the segment names, which on a log prefix would be
	// hundreds of thousands of keys.
	for _, delimiter := range svc.seenDelimiter {
		assert.Equal(t, "/", delimiter)
	}
}

func TestDiscovery_IgnoresNonNumericSegments(t *testing.T) {
	svc := &prefixMockS3{commonPrefixes: map[string][]string{
		"user-report/AWSLogs/123456789012/KiroLogs/user_report/us-east-1/": {"2026", "unexpected", "2025"},
	}}
	years, err := discoveryFixture(svc).ListYears("123456789012")
	require.Nil(t, err)
	// S3 prefixes are strings; a stray directory must not surface as a year.
	assert.Equal(t, []int{2025, 2026}, years)
}

func TestDiscovery_EmptyLayout(t *testing.T) {
	// An empty account list is the clearest signal that the bucket or report
	// prefix is wrong, so it must come back as an empty result rather than an
	// error.
	discovery := discoveryFixture(&prefixMockS3{})

	accounts, err := discovery.ListAccounts()
	require.Nil(t, err)
	assert.Empty(t, accounts)
}

// Per-stream counts separate three situations that a single pass/fail verdict
// cannot: a wrong prefix (zero everywhere), a dormant stream (zero for one type,
// which is what inline completions look like under agentic usage), and a
// permission gap on one bucket.
func TestDiscovery_CountStreams(t *testing.T) {
	base := "user-report/AWSLogs/123456789012/KiroLogs/user_report/us-east-1/2026/07/"
	logBase := "logging/AWSLogs/123456789012/KiroLogs/"
	svc := &prefixMockS3{objects: map[string][]string{
		base: {"KIRO_CLI_x_user_report_1.csv", "KIRO_IDE_x_user_report_1.csv"},
		logBase + "GenerateAssistantResponse/us-east-1/2026/07/": {"a.json.gz", "b.json.gz", "c.json.gz"},
		// GenerateCompletions intentionally absent: dormant, not misconfigured.
	}}

	month := 7
	counts := discoveryFixture(svc).CountStreams("123456789012", 2026, &month, 0)
	require.Len(t, counts, 3)

	byType := map[string]StreamCount{}
	for _, c := range counts {
		byType[c.FileType] = c
	}
	assert.Equal(t, 2, byType[models.FileTypeReport].Count)
	assert.Equal(t, 3, byType[models.FileTypeChatLog].Count)
	assert.Equal(t, 0, byType[models.FileTypeCompletionLog].Count)

	// The prefix is reported alongside the count so a zero can be checked
	// against the bucket directly.
	assert.Contains(t, byType[models.FileTypeReport].Prefix, "user_report/us-east-1/2026/07")
	for _, c := range counts {
		assert.NotEmpty(t, c.Bucket)
		assert.Empty(t, c.Error)
	}
}

func TestCountObjects(t *testing.T) {
	prefix := "p/"
	svc := &prefixMockS3{objects: map[string][]string{
		// Only .csv and .json.gz are collectable; the count must match what
		// would actually be collected, not every object present.
		prefix: {"a.csv", "b.json.gz", "c.txt", "d", "e.zip"},
	}}
	client := &KiroS3Client{S3: svc, Bucket: "b"}

	count, atLeast, err := client.CountObjects(prefix, 0)
	require.Nil(t, err)
	assert.Equal(t, 2, count)
	assert.False(t, atLeast)

	// Counting stops at the cap so the check stays cheap on a large bucket; the
	// flag marks the value as a floor.
	count, atLeast, err = client.CountObjects(prefix, 1)
	require.Nil(t, err)
	assert.Equal(t, 1, count)
	assert.True(t, atLeast)
}

func TestListSubPrefixes_TrailingSlashHandling(t *testing.T) {
	svc := &prefixMockS3{commonPrefixes: map[string][]string{
		"a/b/": {"x", "y"},
	}}
	client := &KiroS3Client{S3: svc, Bucket: "bkt"}

	// A caller-supplied prefix may or may not end in a slash; both must resolve
	// to the same listing, and the child names come back without the slash.
	withSlash, err := client.ListSubPrefixes("a/b/")
	require.Nil(t, err)
	withoutSlash, err := client.ListSubPrefixes("a/b")
	require.Nil(t, err)

	assert.Equal(t, []string{"x", "y"}, withSlash)
	assert.Equal(t, withSlash, withoutSlash)
}
