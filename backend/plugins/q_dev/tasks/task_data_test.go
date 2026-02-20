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

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestQDevTaskData_WithIdentityClient(t *testing.T) {
	taskData := &QDevTaskData{
		Options: &QDevOptions{
			ConnectionId: 1,
			S3Prefix:     "test-prefix/",
		},
		S3Client: &QDevS3Client{
			S3:     &s3.S3{},
			Bucket: "test-bucket",
		},
		IdentityClient: &QDevIdentityClient{
			StoreId: "d-1234567890",
			Region:  "us-west-2",
		},
		S3Prefixes: []string{"test-prefix/"},
	}

	assert.NotNil(t, taskData.IdentityClient)
	assert.Equal(t, "d-1234567890", taskData.IdentityClient.StoreId)
	assert.Equal(t, "us-west-2", taskData.IdentityClient.Region)
	assert.NotNil(t, taskData.S3Client)
	assert.NotNil(t, taskData.Options)
	assert.Equal(t, []string{"test-prefix/"}, taskData.S3Prefixes)
}

func TestQDevTaskData_WithoutIdentityClient(t *testing.T) {
	taskData := &QDevTaskData{
		Options: &QDevOptions{
			ConnectionId: 1,
			S3Prefix:     "test-prefix/",
		},
		S3Client: &QDevS3Client{
			S3:     &s3.S3{},
			Bucket: "test-bucket",
		},
		IdentityClient: nil, // No identity client configured
	}

	assert.Nil(t, taskData.IdentityClient)
	assert.NotNil(t, taskData.S3Client)
	assert.NotNil(t, taskData.Options)
	assert.Equal(t, uint64(1), taskData.Options.ConnectionId)
	assert.Equal(t, "test-prefix/", taskData.Options.S3Prefix)
}

func TestQDevTaskData_AllFields(t *testing.T) {
	month := 3
	options := &QDevOptions{
		ConnectionId: 123,
		S3Prefix:     "data/q-dev/",
		AccountId:    "034362076319",
		BasePath:     "user-report",
		Year:         2026,
		Month:        &month,
	}

	s3Client := &QDevS3Client{
		S3:     &s3.S3{},
		Bucket: "my-data-bucket",
	}

	identityClient := &QDevIdentityClient{
		StoreId: "d-9876543210",
		Region:  "eu-west-1",
	}

	taskData := &QDevTaskData{
		Options:        options,
		S3Client:       s3Client,
		IdentityClient: identityClient,
		S3Prefixes: []string{
			"user-report/AWSLogs/034362076319/KiroLogs/by_user_analytic/us-east-1/2026/03",
			"user-report/AWSLogs/034362076319/KiroLogs/user_report/us-east-1/2026/03",
		},
	}

	// Verify all fields are properly set
	assert.Equal(t, options, taskData.Options)
	assert.Equal(t, s3Client, taskData.S3Client)
	assert.Equal(t, identityClient, taskData.IdentityClient)

	// Verify nested field access
	assert.Equal(t, uint64(123), taskData.Options.ConnectionId)
	assert.Equal(t, "data/q-dev/", taskData.Options.S3Prefix)
	assert.Equal(t, "034362076319", taskData.Options.AccountId)
	assert.Equal(t, "user-report", taskData.Options.BasePath)
	assert.Equal(t, 2026, taskData.Options.Year)
	assert.Equal(t, &month, taskData.Options.Month)
	assert.Equal(t, "my-data-bucket", taskData.S3Client.Bucket)
	assert.Equal(t, "d-9876543210", taskData.IdentityClient.StoreId)
	assert.Equal(t, "eu-west-1", taskData.IdentityClient.Region)
	assert.Len(t, taskData.S3Prefixes, 2)
}

func TestQDevTaskData_EmptyStruct(t *testing.T) {
	taskData := &QDevTaskData{}

	assert.Nil(t, taskData.Options)
	assert.Nil(t, taskData.S3Client)
	assert.Nil(t, taskData.IdentityClient)
}

func TestQDevTaskData_PartialInitialization(t *testing.T) {
	taskData := &QDevTaskData{
		Options: &QDevOptions{
			ConnectionId: 456,
		},
		// S3Client and IdentityClient intentionally nil
	}

	assert.NotNil(t, taskData.Options)
	assert.Equal(t, uint64(456), taskData.Options.ConnectionId)
	assert.Equal(t, "", taskData.Options.S3Prefix) // Default empty string
	assert.Nil(t, taskData.S3Client)
	assert.Nil(t, taskData.IdentityClient)
}
