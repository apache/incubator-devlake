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
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"time"
)

type AzuredevopsRelease struct {
	archived.NoPKModel

	ConnectionId          uint64 `gorm:"primaryKey"`
	AzuredevopsId         int    `gorm:"primaryKey"`
	ProjectId             string `gorm:"type:varchar(255)"`
	Name                  string `gorm:"type:varchar(255)"`
	Status                string `gorm:"type:varchar(100)"`
	ReleaseDefinitionId   int
	ReleaseDefinitionName string `gorm:"type:varchar(255)"`
	Description           string `gorm:"type:text"`
	CreatedOn             *time.Time
	ModifiedOn            *time.Time
}

func (AzuredevopsRelease) TableName() string {
	return "_tool_azuredevops_go_releases"
}

type AzuredevopsReleaseDeployment struct {
	archived.NoPKModel

	ConnectionId     uint64 `gorm:"primaryKey"`
	AzuredevopsId    int    `gorm:"primaryKey"`
	ReleaseId        int    `gorm:"index"`
	ProjectId        string `gorm:"type:varchar(255)"`
	Name             string `gorm:"type:varchar(255)"`
	Status           string `gorm:"type:varchar(100)"`
	OperationStatus  string `gorm:"type:varchar(100)"`
	DeploymentStatus string `gorm:"type:varchar(100)"`
	DefinitionName   string `gorm:"type:varchar(255)"`
	DefinitionId     int
	EnvironmentId    int
	EnvironmentName  string `gorm:"type:varchar(255)"`
	AttemptNumber    int
	Reason           string `gorm:"type:varchar(100)"`
	QueuedOn         *time.Time
	StartedOn        *time.Time
	CompletedOn      *time.Time
	LastModifiedOn   *time.Time
}

func (AzuredevopsReleaseDeployment) TableName() string {
	return "_tool_azuredevops_go_release_deployments"
}
