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
	"time"

	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
)

type BitbucketPipelineStep struct {
	ConnectionId      uint64 `gorm:"primaryKey"`
	BitbucketId       string `gorm:"primaryKey"`
	PipelineId        string `gorm:"type:varchar(255)"`
	Name              string `gorm:"type:varchar(255)"`
	Trigger           string `gorm:"type:varchar(255)"`
	State             string `gorm:"type:varchar(255)"`
	Result            string `gorm:"type:varchar(255)"`
	MaxTime           int
	StartedOn         *time.Time
	CompletedOn       *time.Time
	DurationInSeconds int
	BuildSecondsUsed  int
	RunNumber         int
	archived.NoPKModel
}

func (BitbucketPipelineStep) TableName() string {
	return "_tool_bitbucket_pipeline_steps"
}

type BitbucketPipelineStep20230411 struct {
	RepoId string `gorm:"type:varchar(255)"`
}

func (BitbucketPipelineStep20230411) TableName() string {
	return "_tool_bitbucket_pipeline_steps"
}
