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

package impl

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apache/incubator-devlake/plugins/q_dev/tasks"
)

func TestQDev_BasicPluginMethods(t *testing.T) {
	plugin := &QDev{}

	assert.Equal(t, "q_dev", plugin.Name())
	assert.Equal(t, "To collect and enrich data from AWS Q Developer usage metrics", plugin.Description())
	assert.Equal(t, "github.com/apache/incubator-devlake/plugins/q_dev", plugin.RootPkgPath())

	// Test table info
	tables := plugin.GetTablesInfo()
	assert.Len(t, tables, 3)

	// Test subtask metas
	subtasks := plugin.SubTaskMetas()
	assert.Len(t, subtasks, 2)

	// Test API resources
	apiResources := plugin.ApiResources()
	assert.NotEmpty(t, apiResources)
	assert.Contains(t, apiResources, "test")
	assert.Contains(t, apiResources, "connections")
}

func TestQDev_TaskDataStructure(t *testing.T) {
	// Test that QDevTaskData has the expected structure
	taskData := &tasks.QDevTaskData{
		Options: &tasks.QDevOptions{
			ConnectionId: 1,
			S3Prefix:     "test/",
		},
		S3Client: &tasks.QDevS3Client{
			Bucket: "test-bucket",
		},
		IdentityClient: &tasks.QDevIdentityClient{
			StoreId: "d-1234567890",
			Region:  "us-west-2",
		},
	}

	assert.NotNil(t, taskData.Options)
	assert.NotNil(t, taskData.S3Client)
	assert.NotNil(t, taskData.IdentityClient)

	assert.Equal(t, uint64(1), taskData.Options.ConnectionId)
	assert.Equal(t, "test/", taskData.Options.S3Prefix)
	assert.Equal(t, "test-bucket", taskData.S3Client.Bucket)
	assert.Equal(t, "d-1234567890", taskData.IdentityClient.StoreId)
	assert.Equal(t, "us-west-2", taskData.IdentityClient.Region)
}

func TestQDev_TaskDataWithoutIdentityClient(t *testing.T) {
	// Test that QDevTaskData works without IdentityClient
	taskData := &tasks.QDevTaskData{
		Options: &tasks.QDevOptions{
			ConnectionId: 1,
		},
		S3Client: &tasks.QDevS3Client{
			Bucket: "test-bucket",
		},
		IdentityClient: nil, // No identity client
	}

	assert.NotNil(t, taskData.Options)
	assert.NotNil(t, taskData.S3Client)
	assert.Nil(t, taskData.IdentityClient)
}
