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
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go/service/identitystore"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/incubator-devlake/plugins/kiro/models"
)

// mockS3 records the bucket each call addressed, so tests can prove a request
// went to the right one.
type mockS3 struct {
	getObjectBody   string
	getObjectErr    error
	listOutputs     []*s3.ListObjectsV2Output
	listCallIdx     int
	seenGetBuckets  []string
	seenListBuckets []string
	seenGetKeys     []string
}

func (m *mockS3) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	m.seenGetBuckets = append(m.seenGetBuckets, *input.Bucket)
	m.seenGetKeys = append(m.seenGetKeys, *input.Key)
	if m.getObjectErr != nil {
		return nil, m.getObjectErr
	}
	return &s3.GetObjectOutput{
		Body: io.NopCloser(bytes.NewReader([]byte(m.getObjectBody))),
	}, nil
}

func (m *mockS3) ListObjectsV2(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	m.seenListBuckets = append(m.seenListBuckets, *input.Bucket)
	if m.listCallIdx >= len(m.listOutputs) {
		return &s3.ListObjectsV2Output{}, nil
	}
	out := m.listOutputs[m.listCallIdx]
	m.listCallIdx++
	return out, nil
}

// The fallback rules mean a single-bucket deployment (what real exports look
// like today) and the dual-bucket layout Kiro recommends both work without
// branching in collectors or extractors.
func TestKiroS3Clients_BucketRouting(t *testing.T) {
	t.Run("single bucket routes both file kinds to the same bucket", func(t *testing.T) {
		svc := &mockS3{}
		clients := &KiroS3Clients{
			Report:    &KiroS3Client{S3: svc, Bucket: "one-bucket"},
			PromptLog: &KiroS3Client{S3: svc, Bucket: "one-bucket"},
		}

		assert.Equal(t, "one-bucket", clients.ForFileType(models.FileTypeReport).Bucket)
		assert.Equal(t, "one-bucket", clients.ForFileType(models.FileTypeChatLog).Bucket)
		assert.Equal(t, "one-bucket", clients.ForFileType(models.FileTypeCompletionLog).Bucket)
		// Deduplicated, so a connection test checks access once rather than twice.
		assert.Equal(t, []string{"one-bucket"}, clients.Buckets())
	})

	t.Run("separate buckets route by file type", func(t *testing.T) {
		svc := &mockS3{}
		clients := &KiroS3Clients{
			Report:    &KiroS3Client{S3: svc, Bucket: "reports"},
			PromptLog: &KiroS3Client{S3: svc, Bucket: "logs"},
		}

		assert.Equal(t, "reports", clients.ForFileType(models.FileTypeReport).Bucket)
		assert.Equal(t, "logs", clients.ForFileType(models.FileTypeChatLog).Bucket)
		assert.Equal(t, "logs", clients.ForFileType(models.FileTypeCompletionLog).Bucket)
		assert.Equal(t, []string{"reports", "logs"}, clients.Buckets())
	})

	// An unrecognized file type must not silently read from the report bucket,
	// where it would find nothing; log data is the larger and more likely case.
	t.Run("unknown file type falls to the log bucket", func(t *testing.T) {
		svc := &mockS3{}
		clients := &KiroS3Clients{
			Report:    &KiroS3Client{S3: svc, Bucket: "reports"},
			PromptLog: &KiroS3Client{S3: svc, Bucket: "logs"},
		}
		assert.Equal(t, "logs", clients.ForFileType("something-new").Bucket)
	})
}

func TestNewKiroS3Clients_FallbackFromConnection(t *testing.T) {
	t.Run("no prompt log bucket falls back to the report bucket", func(t *testing.T) {
		conn := &models.KiroConnection{KiroConn: models.KiroConn{
			Region: "us-east-1",
			Bucket: "kiro-export-test",
		}}
		clients, err := NewKiroS3Clients(conn)
		require.Nil(t, err)
		assert.Equal(t, "kiro-export-test", clients.Report.Bucket)
		assert.Equal(t, "kiro-export-test", clients.PromptLog.Bucket)
		assert.Len(t, clients.Buckets(), 1)
	})

	t.Run("explicit prompt log bucket is used", func(t *testing.T) {
		conn := &models.KiroConnection{KiroConn: models.KiroConn{
			Region:          "us-east-1",
			Bucket:          "reports-bucket",
			PromptLogBucket: "logs-bucket",
		}}
		clients, err := NewKiroS3Clients(conn)
		require.Nil(t, err)
		assert.Equal(t, "reports-bucket", clients.Report.Bucket)
		assert.Equal(t, "logs-bucket", clients.PromptLog.Bucket)
		assert.Len(t, clients.Buckets(), 2)
	})
}

func TestKiroS3Client_GetObjectBytes(t *testing.T) {
	t.Run("reads the body and addresses the bound bucket", func(t *testing.T) {
		svc := &mockS3{getObjectBody: "hello"}
		client := &KiroS3Client{S3: svc, Bucket: "my-bucket"}

		data, err := client.GetObjectBytes("some/key.csv")
		require.Nil(t, err)
		assert.Equal(t, "hello", string(data))
		assert.Equal(t, []string{"my-bucket"}, svc.seenGetBuckets)
		assert.Equal(t, []string{"some/key.csv"}, svc.seenGetKeys)
	})

	t.Run("propagates an S3 error", func(t *testing.T) {
		svc := &mockS3{getObjectErr: errors.New("access denied")}
		client := &KiroS3Client{S3: svc, Bucket: "my-bucket"}

		_, err := client.GetObjectBytes("k")
		assert.NotNil(t, err)
	})
}

// mockIdentityStore lets the optional display-name path be exercised without
// AWS.
type mockIdentityStore struct {
	displayName *string
	err         error
}

func (m *mockIdentityStore) DescribeUser(*identitystore.DescribeUserInput) (*identitystore.DescribeUserOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &identitystore.DescribeUserOutput{DisplayName: m.displayName}, nil
}

func TestKiroIdentityClient_ResolveDisplayName(t *testing.T) {
	name := "Some Developer"

	t.Run("resolves a display name", func(t *testing.T) {
		client := &KiroIdentityClient{IdentityStore: &mockIdentityStore{displayName: &name}, StoreId: "d-1"}
		got, err := client.ResolveDisplayName("user-1")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, name, *got)
	})

	// The column exists for human readability; falling back to the raw id would
	// make an unresolved value look like a resolved one.
	t.Run("empty display name yields nil, not the user id", func(t *testing.T) {
		empty := ""
		client := &KiroIdentityClient{IdentityStore: &mockIdentityStore{displayName: &empty}, StoreId: "d-1"}
		got, err := client.ResolveDisplayName("user-1")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	// Identity Store is optional, so an unconfigured client must be safe to
	// call rather than something every caller has to nil-check.
	t.Run("nil client is safe to call", func(t *testing.T) {
		var client *KiroIdentityClient
		got, err := client.ResolveDisplayName("user-1")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("error surfaces but yields no name", func(t *testing.T) {
		client := &KiroIdentityClient{IdentityStore: &mockIdentityStore{err: errors.New("throttled")}, StoreId: "d-1"}
		got, err := client.ResolveDisplayName("user-1")
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestNewKiroIdentityClient_OptionalConfiguration(t *testing.T) {
	// Missing configuration is not an error: collection works fully without
	// display names because identity comes from the report's email column.
	for _, tt := range []struct {
		name string
		conn models.KiroConn
	}{
		{"neither set", models.KiroConn{}},
		{"only store id", models.KiroConn{IdentityStoreId: "d-1"}},
		{"only region", models.KiroConn{IdentityStoreRegion: "us-east-1"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewKiroIdentityClient(&models.KiroConnection{KiroConn: tt.conn})
			require.NoError(t, err)
			assert.Nil(t, client)
		})
	}

	t.Run("fully configured returns a client", func(t *testing.T) {
		client, err := NewKiroIdentityClient(&models.KiroConnection{KiroConn: models.KiroConn{
			IdentityStoreId:     "d-1234567890",
			IdentityStoreRegion: "us-east-1",
		}})
		require.NoError(t, err)
		require.NotNil(t, client)
		assert.Equal(t, "d-1234567890", client.StoreId)
	})
}
