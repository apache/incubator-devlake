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
	"github.com/aws/aws-sdk-go/service/s3"
)

type QDevApiParams struct {
	ConnectionId uint64 `json:"connectionId"`
}

type QDevOptions struct {
	ConnectionId uint64 `json:"connectionId"`
	S3Prefix     string `json:"s3Prefix"`
	ScopeId      string `json:"scopeId"`
}

type QDevTaskData struct {
	Options        *QDevOptions
	S3Client       *QDevS3Client
	IdentityClient *QDevIdentityClient // New field for Identity Center client
}

type QDevS3Client struct {
	S3     *s3.S3
	Bucket string
}

func (client *QDevS3Client) Close() {
	// S3客户端不需要特别关闭操作
}
