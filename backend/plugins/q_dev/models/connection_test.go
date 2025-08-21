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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQDevConn_WithIdentityStore(t *testing.T) {
	conn := QDevConn{
		AccessKeyId:         "test-key",
		SecretAccessKey:     "test-secret",
		Region:              "us-east-1",
		Bucket:              "test-bucket",
		RateLimitPerHour:    20000,
		IdentityStoreId:     "d-1234567890",
		IdentityStoreRegion: "us-west-2",
	}

	assert.Equal(t, "d-1234567890", conn.IdentityStoreId)
	assert.Equal(t, "us-west-2", conn.IdentityStoreRegion)
	assert.Equal(t, "us-east-1", conn.Region) // S3 region
}

func TestQDevConn_RequiredFields(t *testing.T) {
	// Test that all required fields are present
	conn := QDevConn{
		AccessKeyId:         "test-key",
		SecretAccessKey:     "test-secret",
		Region:              "us-east-1",
		Bucket:              "test-bucket",
		RateLimitPerHour:    20000,
		IdentityStoreId:     "d-1234567890",
		IdentityStoreRegion: "us-west-2",
	}

	// All required fields should be non-empty
	assert.NotEmpty(t, conn.AccessKeyId)
	assert.NotEmpty(t, conn.SecretAccessKey)
	assert.NotEmpty(t, conn.Region)
	assert.NotEmpty(t, conn.Bucket)
	assert.NotEmpty(t, conn.IdentityStoreId)
	assert.NotEmpty(t, conn.IdentityStoreRegion)
	assert.Greater(t, conn.RateLimitPerHour, 0)
}

func TestQDevConn_Sanitize_PreservesIdentityStore(t *testing.T) {
	conn := QDevConn{
		SecretAccessKey:     "secret-key",
		IdentityStoreId:     "d-1234567890",
		IdentityStoreRegion: "us-west-2",
	}

	sanitized := conn.Sanitize()
	assert.NotEqual(t, "secret-key", sanitized.SecretAccessKey)
	assert.Equal(t, "d-1234567890", sanitized.IdentityStoreId)
	assert.Equal(t, "us-west-2", sanitized.IdentityStoreRegion)
}

func TestQDevConnection_WithIdentityStore(t *testing.T) {
	connection := QDevConnection{
		QDevConn: QDevConn{
			AccessKeyId:         "test-key",
			SecretAccessKey:     "test-secret",
			Region:              "us-east-1",
			Bucket:              "test-bucket",
			IdentityStoreId:     "d-1234567890",
			IdentityStoreRegion: "us-west-2",
		},
	}

	assert.Equal(t, "d-1234567890", connection.IdentityStoreId)
	assert.Equal(t, "us-west-2", connection.IdentityStoreRegion)
}

func TestQDevConnection_Sanitize_WithIdentityStore(t *testing.T) {
	connection := QDevConnection{
		QDevConn: QDevConn{
			SecretAccessKey:     "secret-key",
			IdentityStoreId:     "d-1234567890",
			IdentityStoreRegion: "us-west-2",
		},
	}

	sanitized := connection.Sanitize()
	assert.NotEqual(t, "secret-key", sanitized.SecretAccessKey)
	assert.Equal(t, "d-1234567890", sanitized.IdentityStoreId)
	assert.Equal(t, "us-west-2", sanitized.IdentityStoreRegion)
}

func TestQDevConnection_TableName(t *testing.T) {
	connection := QDevConnection{}
	assert.Equal(t, "_tool_q_dev_connections", connection.TableName())
}
