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
	"github.com/apache/incubator-devlake/core/models/domainlayer"
)

type CicdDeploymentCommit struct {
	domainlayer.DomainEntity
	CicdScopeId         string `gorm:"index;type:varchar(255)"`
	CicdDeploymentId    string `gorm:"type:varchar(255)"` // if it is converted from a cicd_pipeline_commit
	Name                string `gorm:"type:varchar(255)"`
	Result              string `gorm:"type:varchar(100)"`
	Status              string `gorm:"type:varchar(100)"`
	OriginalStatus      string `gorm:"type:varchar(100)"`
	OriginalResult      string `gorm:"type:varchar(100)"`
	Environment         string `gorm:"type:varchar(255)"`
	OriginalEnvironment string `gorm:"type:varchar(255)"`
	TaskDatesInfo
	DurationSec                   *float64
	QueuedDurationSec             *float64
	CommitSha                     string `gorm:"primaryKey;type:varchar(255)"`
	CommitMsg                     string
	RefName                       string `gorm:"type:varchar(255)"` // to delete?
	RepoId                        string `gorm:"type:varchar(255)"`
	RepoUrl                       string `gorm:"index;not null"`
	PrevSuccessDeploymentCommitId string `gorm:"type:varchar(255)"`
}

func (cicdDeploymentCommit CicdDeploymentCommit) TableName() string {
	return "cicd_deployment_commits"
}

func (cicdDeploymentCommit CicdDeploymentCommit) ToDeployment() *CICDDeployment {
	return &CICDDeployment{
		DomainEntity: domainlayer.DomainEntity{
			Id:        cicdDeploymentCommit.CicdDeploymentId,
			NoPKModel: cicdDeploymentCommit.DomainEntity.NoPKModel,
		},
		CicdScopeId:         cicdDeploymentCommit.CicdScopeId,
		Name:                cicdDeploymentCommit.Name,
		Result:              cicdDeploymentCommit.Result,
		Status:              cicdDeploymentCommit.Status,
		OriginalStatus:      cicdDeploymentCommit.OriginalStatus,
		OriginalResult:      cicdDeploymentCommit.OriginalResult,
		Environment:         cicdDeploymentCommit.Environment,
		OriginalEnvironment: cicdDeploymentCommit.OriginalEnvironment,
		TaskDatesInfo: TaskDatesInfo{
			CreatedDate:  cicdDeploymentCommit.CreatedDate,
			QueuedDate:   cicdDeploymentCommit.QueuedDate,
			StartedDate:  cicdDeploymentCommit.StartedDate,
			FinishedDate: cicdDeploymentCommit.FinishedDate,
		},
		QueuedDurationSec: cicdDeploymentCommit.QueuedDurationSec,
		DurationSec:       cicdDeploymentCommit.DurationSec,
	}
}
