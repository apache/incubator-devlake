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

// SubProjectDeployment attributes a deployment to a single sub-project of a monorepo,
// based on the name of the CI job that performed the deployment.
//
// One deployment may produce several rows when a single pipeline runs the deploy jobs
// of several sub-projects. That is not double counting: each sub-project really was
// deployed by that pipeline.
type SubProjectDeployment struct {
	common.NoPKModel
	// The four primary key columns are deliberately kept narrow: MySQL caps a composite
	// index at 3072 bytes, which is 768 characters under utf8mb4.
	ProjectName string `gorm:"primaryKey;type:varchar(100)"`
	SubProject  string `gorm:"primaryKey;type:varchar(100)"`
	// CicdDeploymentId is the id of the deployment (a cicd_pipelines.id when the
	// deployment was generated from a pipeline), taken from cicd_deployment_commits.
	CicdDeploymentId string `gorm:"primaryKey;type:varchar(255)"`
	// CommitSha is wide enough for a SHA-256 hash; the source column is varchar(255) but
	// only ever holds a git object id.
	CommitSha string `gorm:"primaryKey;type:varchar(64)"`
	// JobName is the cicd_tasks.name that matched this sub-project's DeployJobPattern.
	JobName      string `gorm:"type:varchar(255)"`
	Result       string `gorm:"type:varchar(100)"`
	Environment  string `gorm:"type:varchar(255)"`
	FinishedDate *time.Time
}

func (SubProjectDeployment) TableName() string {
	return "monorepo_subproject_deployments"
}
