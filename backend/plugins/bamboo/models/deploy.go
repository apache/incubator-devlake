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

import "github.com/apache/incubator-devlake/core/models/common"

type BambooDeployEnvironment struct {
	ConnectionId        uint64 `json:"connection_id" gorm:"primaryKey"`
	EnvId               uint64 `json:"env_id" gorm:"primaryKey"`
	EnvKey              string `json:"key" gorm:"index"`
	Name                string `json:"name" gorm:"index"`
	PlanKey             string `json:"plan_key" gorm:"index"`
	ProjectKey          string `json:"project_key" gorm:"index"`
	Description         string `json:"description"`
	DeploymentProjectId uint64 `json:"deploymentProjectId"`
	Position            uint64 `json:"position"`
	ConfigurationState  string `json:"configurationState"`

	ApiBambooOperations
	common.NoPKModel
}

func (b *BambooDeployEnvironment) Convert(apiEnv *ApiBambooEnvironment) {
	b.EnvId = apiEnv.ID
	b.EnvKey = apiEnv.Key.Key
	b.Name = apiEnv.Name
	b.Description = apiEnv.Description
	b.DeploymentProjectId = apiEnv.DeploymentProjectId
	b.Position = apiEnv.Position
	b.ConfigurationState = apiEnv.ConfigurationState
	b.ApiBambooOperations = apiEnv.Operations
}

func (BambooDeployEnvironment) TableName() string {
	return "_tool_bamboo_deploy_environment"
}

type ApiBambooEnvironment struct {
	ID                  uint64              `json:"id"`
	Key                 ApiBambooKey        `json:"key"`
	Name                string              `json:"name"`
	Description         string              `json:"description"`
	DeploymentProjectId uint64              `json:"deploymentProjectId"`
	Operations          ApiBambooOperations `json:"operations"`
	Position            uint64              `json:"position"`
	ConfigurationState  string              `json:"configurationState"`
}

type ApiBambooDeployProject struct {
	ID           uint64                 `json:"id"`
	OID          string                 `json:"oid"`
	Key          ApiBambooKey           `json:"key"`
	Name         string                 `json:"name"`
	PlanKey      ApiBambooKey           `json:"planKey"`
	Description  string                 `json:"description"`
	Environments []ApiBambooEnvironment `json:"environments"`
	Operations   ApiBambooOperations    `json:"operations"`
}
