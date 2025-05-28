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

package archived

import (
	commonArchived "github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
)

// QDevConn holds the essential information to connect to AWS S3
type QDevConn struct {
	// AccessKeyId for AWS
	AccessKeyId string `mapstructure:"accessKeyId" json:"accessKeyId"`
	// SecretAccessKey for AWS
	SecretAccessKey string `mapstructure:"secretAccessKey" json:"secretAccessKey"`
	// Region for AWS
	Region string `mapstructure:"region" json:"region"`
	// Bucket for AWS S3
	Bucket string `mapstructure:"bucket" json:"bucket"`
	// RateLimitPerHour limits the API requests sent to AWS
	RateLimitPerHour int `mapstructure:"rateLimitPerHour" json:"rateLimitPerHour"`
}

// QDevConnection holds QDevConn plus ID/Name for database storage
type QDevConnection struct {
	commonArchived.BaseConnection `mapstructure:",squash"`
	QDevConn                      `mapstructure:",squash"`
}

func (QDevConnection) TableName() string {
	return "_tool_q_dev_connections"
}
