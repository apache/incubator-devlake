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
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
)

type BitbucketPipeline struct {
	ConnectionId      uint64 `gorm:"primaryKey"`
	BitbucketId       string `gorm:"primaryKey"`
	Status            string `gorm:"type:varchar(100)"`
	Result            string `gorm:"type:varchar(100)"`
	RefName           string `gorm:"type:varchar(255)"`
	RepoId            string `gorm:"type:varchar(255)"`
	CommitSha         string `gorm:"type:varchar(255)"`
	WebUrl            string `gorm:"type:varchar(255)"`
	Type              string `gorm:"type:varchar(255)"`
	Environment       string `gorm:"type:varchar(255)"`
	BuildNumber       int
	DurationInSeconds uint64

	BitbucketCreatedOn  *time.Time
	BitbucketCompleteOn *time.Time

	common.NoPKModel
}

func (BitbucketPipeline) TableName() string {
	return "_tool_bitbucket_pipelines"
}

const (
	// https://github.com/juan-carlos-duarte/bitbucket-api-client-lib/blob/main/docs/PipelineState.md
	// https://community.atlassian.com/t5/Bitbucket-questions/Possible-value-for-status-filtering-in-API-pipelines/qaq-p/1755200
	// https://github.com/juan-carlos-duarte/bitbucket-api-client-lib/blob/main/docs/PipelineState.md
	// https://developer.atlassian.com/server/bitbucket/rest/v815/api-group-builds-and-deployments/#api-api-latest-projects-projectkey-repos-repositoryslug-commits-commitid-deployments-get
	// https://github.com/juan-carlos-duarte/bitbucket-api-client-lib/blob/main/docs/DeploymentStateCompleted.md
	FAILED      = "FAILED"
	ERROR       = "ERROR"
	UNDEPLOYED  = "UNDEPLOYED"
	STOPPED     = "STOPPED"
	SKIPPED     = "SKIPPED"
	SUCCESSFUL  = "SUCCESSFUL"
	COMPLETED   = "COMPLETED"
	PAUSED      = "PAUSED"
	HALTED      = "HALTED"
	IN_PROGRESS = "IN_PROGRESS"
	PENDING     = "PENDING"
	BUILDING    = "BUILDING"
	EXPIRED     = "EXPIRED"
	RUNNING     = "RUNNING"
	READY       = "READY"
	PASSED      = "PASSED"
	NOT_RUN     = "NOT_RUN"
	CANCELLED   = "CANCELLED"
)
