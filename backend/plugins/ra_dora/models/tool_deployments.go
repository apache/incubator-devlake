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

// Deployment represents the entire JSON structure
type Deployment struct {
	common.NoPKModel `swaggerignore:"true" json:"-" mapstructure:"-"`

	ConnectionId uint64   `gorm:"primaryKey"`
	Metadata     Metadata `json:"metadata"`
	Spec         Spec     `json:"spec"`
	Status       Status   `json:"status"`
}

// Metadata represents the metadata part of the JSON
type Metadata struct {
	Name              string `json:"name" gorm:"type:varchar(255)"`
	GenerateName      string `json:"generateName" gorm:"type:varchar(255)"`
	Namespace         string `json:"namespace" gorm:"type:varchar(255)"`
	UID               string `json:"uid" gorm:"type:varchar(255)"`
	ResourceVersion   string `json:"resourceVersion" gorm:"type:varchar(255)"`
	Generation        int    `json:"generation" gorm:"type:int"`
	CreationTimestamp string `json:"creationTimestamp" gorm:"type:varchar(255)"`
	Labels            string `json:"labels" gorm:"type:text"`
	Annotations       string `json:"annotations" gorm:"type:text"`
	ManagedFields     string `json:"managedFields" gorm:"type:text"`
}

// Spec represents the spec part of the JSON
type Spec struct {
	Entrypoint          string `json:"entrypoint" gorm:"type:varchar(255)"`
	Arguments           string `json:"arguments" gorm:"type:text"`
	WorkflowTemplateRef string `json:"workflowTemplateRef" gorm:"type:text"`
}

// Status represents the status part of the JSON
type Status struct {
	Phase      string `json:"phase" gorm:"type:varchar(255)"`
	StartedAt  string `json:"startedAt" gorm:"type:varchar(255)"`
	FinishedAt string `json:"finishedAt" gorm:"type:varchar(255)"`
	Progress   string `json:"progress" gorm:"type:varchar(255)"`
	Nodes      string `json:"nodes" gorm:"type:text"`
}

func (Deployment) TableName() string {
	return "_tool_argo_api_deployments"
}
