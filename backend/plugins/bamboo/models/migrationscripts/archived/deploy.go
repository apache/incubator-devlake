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

import "github.com/apache/incubator-devlake/core/models/migrationscripts/archived"

type ApiBambooOperations struct {
	CanView                   bool `json:"canView"`
	CanEdit                   bool `json:"canEdit"`
	CanDelete                 bool `json:"canDelete"`
	AllowedToExecute          bool `json:"allowedToExecute"`
	CanExecute                bool `json:"canExecute"`
	AllowedToCreateVersion    bool `json:"allowedToCreateVersion"`
	AllowedToSetVersionStatus bool `json:"allowedToSetVersionStatus"`
}

type BambooDeployEnvironment struct {
	ConnectionId        uint64 `json:"connection_id" gorm:"primaryKey"`
	EnvId               uint64 `json:"env_id" gorm:"primaryKey"`
	EnvKey              string `json:"key" gorm:"index;type:varchar(255)"`
	Name                string `json:"name" gorm:"index;type:varchar(255)"`
	PlanKey             string `json:"plan_key" gorm:"index;type:varchar(255)"`
	ProjectKey          string `json:"project_key" gorm:"index"`
	Description         string `json:"description"`
	DeploymentProjectId uint64 `json:"deploymentProjectId"`
	Position            uint64 `json:"position"`
	ConfigurationState  string `json:"configurationState"`

	ApiBambooOperations
	archived.NoPKModel
}

func (BambooDeployEnvironment) TableName() string {
	return "_tool_bamboo_deploy_environment"
}
