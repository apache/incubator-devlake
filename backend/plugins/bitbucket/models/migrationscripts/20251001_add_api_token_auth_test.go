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

package migrationscripts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddApiTokenAuth_Version(t *testing.T) {
	script := &addApiTokenAuth{}
	assert.Equal(t, uint64(20251001000001), script.Version())
}

func TestAddApiTokenAuth_Name(t *testing.T) {
	script := &addApiTokenAuth{}
	assert.Equal(t, "add API token authentication support to Bitbucket connections", script.Name())
}

func TestBitbucketConnection20251001_TableName(t *testing.T) {
	conn := bitbucketConnection20251001{}
	assert.Equal(t, "_tool_bitbucket_connections", conn.TableName())
}

func TestBitbucketConnection20251001_Structure(t *testing.T) {
	// Test that the migration struct has the correct field
	conn := bitbucketConnection20251001{
		UsesApiToken: true,
	}
	assert.True(t, conn.UsesApiToken)

	conn2 := bitbucketConnection20251001{
		UsesApiToken: false,
	}
	assert.False(t, conn2.UsesApiToken)
}

// Note: Full integration test of the Up() method requires a test database setup.
// The migration is tested in practice when running the actual migrations against a database.
// For unit testing purposes, we verify the structure and metadata.
