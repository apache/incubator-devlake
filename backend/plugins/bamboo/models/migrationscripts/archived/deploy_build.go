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

type BambooDeployBuild struct {
	ConnectionId  uint64 `json:"connection_id" gorm:"primaryKey"`
	DeployBuildId uint64 `json:"deploy_build_id" gorm:"primaryKey"`

	DeploymentVersionName string     `json:"deploymentVersionName"`
	DeploymentState       string     `json:"deploymentState"`
	LifeCycleState        string     `json:"lifeCycleState"`
	StartedDate           *time.Time `json:"startedDate"`
	QueuedDate            *time.Time `json:"queuedDate"`
	ExecutedDate          *time.Time `json:"executedDate"`
	FinishedDate          *time.Time `json:"finishedDate"`
	ReasonSummary         string     `json:"reasonSummary"`

	ProjectKey string `json:"project_key" gorm:"index"`
	PlanKey    string `json:"plan_key" gorm:"index"`

	ApiBambooOperations
	archived.NoPKModel
}

func (BambooDeployBuild) TableName() string {
	return "_tool_bamboo_deploy_build"
}
