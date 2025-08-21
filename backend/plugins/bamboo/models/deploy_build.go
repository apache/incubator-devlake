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

type BambooDeployBuild struct {
	ConnectionId          uint64     `json:"connection_id" gorm:"primaryKey"`
	DeployBuildId         uint64     `json:"deploy_build_id" gorm:"primaryKey"`
	PlanResultKey         string     `json:"planResultKey"`
	EnvKey                string     `json:"envKey"`
	DeploymentVersionName string     `json:"deploymentVersionName"`
	DeploymentState       string     `json:"deploymentState"`
	LifeCycleState        string     `json:"lifeCycleState"`
	StartedDate           *time.Time `json:"startedDate"`
	QueuedDate            *time.Time `json:"queuedDate"`
	ExecutedDate          *time.Time `json:"executedDate"`
	FinishedDate          *time.Time `json:"finishedDate"`
	ReasonSummary         string     `json:"reasonSummary"`
	ProjectKey            string     `json:"project_key" gorm:"index"`
	PlanKey               string     `json:"plan_key" gorm:"index"`
	Environment           string     `gorm:"type:varchar(255)"`
	PlanBranchName        string     `gorm:"type:varchar(255)"`
	ApiBambooOperations
	common.NoPKModel
}

func (BambooDeployBuild) TableName() string {
	return "_tool_bamboo_deploy_builds"
}

type ApiBambooDeployBuild struct {
	DeploymentVersion     ApiBambooDeploymentVersion `json:"deploymentVersion"`
	DeploymentVersionName string                     `json:"deploymentVersionName"`
	Id                    uint64                     `json:"id"`
	DeploymentState       string                     `json:"deploymentState"`
	LifeCycleState        string                     `json:"lifeCycleState"`
	StartedDate           int64                      `json:"startedDate"`
	QueuedDate            int64                      `json:"queuedDate"`
	ExecutedDate          int64                      `json:"executedDate"`
	FinishedDate          int64                      `json:"finishedDate"`
	ReasonSummary         string                     `json:"reasonSummary"`
	Key                   ApiBambooDeployBuildKey    `json:"key"`
	Agent                 ApiBambooDeployBuildAgent  `json:"agent"`
	Operations            ApiBambooOperations        `json:"operations"`
}

func (api *ApiBambooDeployBuild) Convert(op *BambooOptions, envName string) []*BambooDeployBuild {
	var result []*BambooDeployBuild
	tmpl := BambooDeployBuild{
		ConnectionId:          op.ConnectionId,
		PlanKey:               op.PlanKey,
		DeploymentVersionName: api.DeploymentVersionName,
		DeployBuildId:         api.Id,
		DeploymentState:       api.DeploymentState,
		LifeCycleState:        api.LifeCycleState,
		StartedDate:           unixForBambooDeployBuild(api.StartedDate),
		QueuedDate:            unixForBambooDeployBuild(api.QueuedDate),
		ExecutedDate:          unixForBambooDeployBuild(api.ExecutedDate),
		FinishedDate:          unixForBambooDeployBuild(api.FinishedDate),
		ReasonSummary:         api.ReasonSummary,
		ApiBambooOperations:   api.Operations,
		Environment:           envName,
		PlanBranchName:        api.DeploymentVersion.PlanBranchName,
		EnvKey:                api.Key.EntityKey.Key,
	}
	for _, item := range api.DeploymentVersion.Items {
		build := tmpl
		build.PlanResultKey = item.PlanResultKey.Key
		if build.PlanResultKey != "" {
			result = append(result, &build)
		}
	}
	return result
}

type ApiBambooDeployBuildAgent struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`

	Type   string `json:"type"`
	Active bool   `json:"active"`
	Enable bool   `json:"enabled"`
	Busy   bool   `json:"busy"`
}

type ApiBambooDeployBuildKey struct {
	Key          string       `json:"key"`
	EntityKey    ApiBambooKey `json:"entityKey"`
	ResultNumber uint64       `json:"resultNumber"`
}

type ApiBambooDeploymentVersion struct {
	ID              uint64 `json:"id"`
	Name            string `json:"name"`
	CreationDate    int64  `json:"creationDate"`
	CreatorUserName string `json:"creatorUserName"`
	Items           []struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		PlanResultKey struct {
			Key       string `json:"key"`
			EntityKey struct {
				Key string `json:"key"`
			} `json:"entityKey"`
			ResultNumber int `json:"resultNumber"`
		} `json:"planResultKey"`
		Type        string `json:"type"`
		Label       string `json:"label"`
		Location    string `json:"location"`
		CopyPattern string `json:"copyPattern"`
		Size        int    `json:"size"`
	} `json:"items"`
	PlanBranchName string `json:"planBranchName"`
}
