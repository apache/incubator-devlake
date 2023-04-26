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

package devops

import (
	"time"

	"github.com/apache/incubator-devlake/core/models/domainlayer"
)

type CicdDeploymentCommit struct {
	domainlayer.DomainEntity
	CicdScopeId                   string `gorm:"index;type:varchar(255)"`
	CicdDeploymentId              string `gorm:"type:varchar(255)"` // if it is converted from a cicd_pipeline_commit
	Name                          string `gorm:"type:varchar(255)"`
	Result                        string `gorm:"type:varchar(100)"`
	Status                        string `gorm:"type:varchar(100)"`
	Environment                   string `gorm:"type:varchar(255)"`
	CreatedDate                   time.Time
	StartedDate                   *time.Time
	FinishedDate                  *time.Time
	DurationSec                   *uint64
	CommitSha                     string `gorm:"primaryKey;type:varchar(255)"`
	RefName                       string `gorm:"type:varchar(255)"` // to delete?
	RepoId                        string `gorm:"type:varchar(255)"`
	RepoUrl                       string `gorm:"index;not null"`
	PrevSuccessDeploymentCommitId string `gorm:"type:varchar(255)"`
}

func (CicdDeploymentCommit) TableName() string {
	return "cicd_deployment_commits"
}
