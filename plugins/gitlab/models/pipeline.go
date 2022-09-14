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

	"github.com/apache/incubator-devlake/models/common"
)

type GitlabPipeline struct {
	ConnectionId uint64 `gorm:"primaryKey"`

	GitlabId  int    `gorm:"primaryKey"`
	ProjectId int    `gorm:"index"`
	Status    string `gorm:"type:varchar(100)"`
	Ref       string `gorm:"type:varchar(255)"`
	Sha       string `gorm:"type:varchar(255)"`
	WebUrl    string `gorm:"type:varchar(255)"`
	Duration  int

	GitlabCreatedAt *time.Time
	GitlabUpdatedAt *time.Time
	StartedAt       *time.Time
	FinishedAt      *time.Time
	Coverage        string

	common.NoPKModel
}

func (GitlabPipeline) TableName() string {
	return "_tool_gitlab_pipelines"
}

type GitlabPipelineProject struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	PipelineId   int    `gorm:"primaryKey"`
	ProjectId    int    `gorm:"primaryKey"`
	Ref          string `gorm:"type:varchar(255)"`
	Sha          string `gorm:"type:varchar(255)"`
	common.NoPKModel
}

func (GitlabPipelineProject) TableName() string {
	return "_tool_gitlab_pipeline_projects"
}
