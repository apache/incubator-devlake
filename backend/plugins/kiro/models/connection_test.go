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

// The fallback rules exist so that a single-bucket setup (what real exports
// currently look like) and the dual-bucket layout Kiro recommends both work
// without branching anywhere else in the plugin.
func TestConnectionPrefixAndBucketFallback(t *testing.T) {
	t.Run("single bucket with defaults", func(t *testing.T) {
		conn := &KiroConn{Bucket: "kiro-export-test"}
		assert.Equal(t, "user-report", conn.GetReportPrefix())
		assert.Equal(t, "logging", conn.GetPromptLogPrefix())
		assert.Equal(t, "kiro-export-test", conn.GetPromptLogBucket())
		assert.False(t, conn.UsesSeparateBuckets())
	})

	t.Run("single bucket with explicit prefixes", func(t *testing.T) {
		conn := &KiroConn{
			Bucket:          "kiro-export-test",
			ReportPrefix:    "user-report",
			PromptLogPrefix: "logging",
		}
		assert.Equal(t, "user-report", conn.GetReportPrefix())
		assert.Equal(t, "logging", conn.GetPromptLogPrefix())
		assert.False(t, conn.UsesSeparateBuckets())
	})

	t.Run("separate buckets", func(t *testing.T) {
		conn := &KiroConn{
			Bucket:          "reports-bucket",
			PromptLogBucket: "logs-bucket",
		}
		assert.Equal(t, "reports-bucket", conn.Bucket)
		assert.Equal(t, "logs-bucket", conn.GetPromptLogBucket())
		assert.True(t, conn.UsesSeparateBuckets())
	})

	t.Run("prefixes are stripped of slashes and whitespace", func(t *testing.T) {
		conn := &KiroConn{
			Bucket:          " kiro-export-test ",
			ReportPrefix:    "/custom-report/",
			PromptLogPrefix: " /custom-logs/ ",
		}
		assert.Equal(t, "custom-report", conn.GetReportPrefix())
		assert.Equal(t, "custom-logs", conn.GetPromptLogPrefix())
		// Whitespace-only difference must not read as two buckets.
		assert.False(t, conn.UsesSeparateBuckets())
	})

	t.Run("whitespace-only prompt log bucket falls back", func(t *testing.T) {
		conn := &KiroConn{Bucket: "b", PromptLogBucket: "   "}
		assert.Equal(t, "b", conn.GetPromptLogBucket())
		assert.False(t, conn.UsesSeparateBuckets())
	})
}

func TestConnectionSanitize(t *testing.T) {
	conn := KiroConnection{
		KiroConn: KiroConn{
			AccessKeyId:     "AKIAEXAMPLE",
			SecretAccessKey: "super-secret-value",
		},
	}
	sanitized := conn.Sanitize()
	assert.NotEqual(t, "super-secret-value", sanitized.SecretAccessKey)
	// The key id is not a secret and stays readable for troubleshooting.
	assert.Equal(t, "AKIAEXAMPLE", sanitized.AccessKeyId)
}
