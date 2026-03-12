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
	"github.com/apache/incubator-devlake/core/models/common"
	"time"
)

// AzuredevopsRelease represents a release from Azure DevOps Release Pipelines (Classic)
type AzuredevopsRelease struct {
	common.NoPKModel

	ConnectionId          uint64 `gorm:"primaryKey"`
	AzuredevopsId         int    `json:"id" gorm:"primaryKey"`
	ProjectId             string
	Name                  string
	Status                string
	ReleaseDefinitionId   int
	ReleaseDefinitionName string
	Description           string
	CreatedOn             *time.Time
	ModifiedOn            *time.Time
}

func (AzuredevopsRelease) TableName() string {
	return "_tool_azuredevops_go_releases"
}

// AzuredevopsReleaseDeployment represents a deployment (environment) within a release
type AzuredevopsReleaseDeployment struct {
	common.NoPKModel

	ConnectionId      uint64 `gorm:"primaryKey"`
	AzuredevopsId     int    `json:"id" gorm:"primaryKey"`
	ReleaseId         int    `gorm:"index"`
	ProjectId         string
	Name              string
	Status            string
	OperationStatus   string
	DeploymentStatus  string
	DefinitionName    string
	DefinitionId      int
	EnvironmentId     int
	EnvironmentName   string
	AttemptNumber     int
	Reason            string
	QueuedOn          *time.Time
	StartedOn         *time.Time
	CompletedOn       *time.Time
	LastModifiedOn    *time.Time
}

func (AzuredevopsReleaseDeployment) TableName() string {
	return "_tool_azuredevops_go_release_deployments"
}

// AzuredevopsApiRelease is the API response structure from Azure DevOps Release API
type AzuredevopsApiRelease struct {
	Id                int        `json:"id"`
	Name              string     `json:"name"`
	Status            string     `json:"status"`
	CreatedOn         *time.Time `json:"createdOn"`
	ModifiedOn        *time.Time `json:"modifiedOn"`
	Description       string     `json:"description"`
	ReleaseDefinition struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Url  string `json:"url"`
		Path string `json:"path"`
	} `json:"releaseDefinition"`
	Environments []AzuredevopsApiReleaseEnvironment `json:"environments"`
	ProjectReference struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"projectReference"`
}

// AzuredevopsApiReleaseEnvironment represents an environment in the release
type AzuredevopsApiReleaseEnvironment struct {
	Id              int        `json:"id"`
	ReleaseId       int        `json:"releaseId"`
	Name            string     `json:"name"`
	Status          string     `json:"status"`
	DeploySteps     []AzuredevopsApiDeployStep `json:"deploySteps"`
	PreDeployApprovals []struct {
		Status string `json:"status"`
	} `json:"preDeployApprovals"`
	PostDeployApprovals []struct {
		Status string `json:"status"`
	} `json:"postDeployApprovals"`
}

// AzuredevopsApiDeployStep represents a deployment step
type AzuredevopsApiDeployStep struct {
	Id               int        `json:"id"`
	DeploymentId     int        `json:"deploymentId"`
	Attempt          int        `json:"attempt"`
	Reason           string     `json:"reason"`
	Status           string     `json:"status"`
	OperationStatus  string     `json:"operationStatus"`
	QueuedOn         *time.Time `json:"queuedOn"`
	LastModifiedOn   *time.Time `json:"lastModifiedOn"`
	HasStarted       bool       `json:"hasStarted"`
}

// AzuredevopsApiDeployment is the API response structure for deployments
type AzuredevopsApiDeployment struct {
	Id                int        `json:"id"`
	Release           struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"release"`
	ReleaseDefinition struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Path string `json:"path"`
	} `json:"releaseDefinition"`
	ReleaseEnvironment struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"releaseEnvironment"`
	ProjectReference struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"projectReference"`
	DefinitionEnvironmentId int        `json:"definitionEnvironmentId"`
	Attempt                 int        `json:"attempt"`
	Reason                  string     `json:"reason"`
	DeploymentStatus        string     `json:"deploymentStatus"`
	OperationStatus         string     `json:"operationStatus"`
	QueuedOn                *time.Time `json:"queuedOn"`
	StartedOn               *time.Time `json:"startedOn"`
	CompletedOn             *time.Time `json:"completedOn"`
	LastModifiedOn          *time.Time `json:"lastModifiedOn"`
}