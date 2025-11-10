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

type ArgocdConnection struct {
	Name                    string `gorm:"type:varchar(100);uniqueIndex" json:"name"`
	archived.RestConnection `mapstructure:",squash"`
	archived.AccessToken    `mapstructure:",squash"`
	archived.Model
}

func (ArgocdConnection) TableName() string {
	return "_tool_argocd_connections"
}

type ArgocdApplication struct {
	ConnectionId   uint64 `gorm:"primaryKey"`
	Name           string `gorm:"primaryKey;type:varchar(255)"`
	Namespace      string `gorm:"type:varchar(255)"`
	Project        string `gorm:"type:varchar(255)"`
	RepoURL        string `gorm:"type:varchar(500)"`
	Path           string `gorm:"type:varchar(255)"`
	TargetRevision string `gorm:"type:varchar(255)"`
	DestServer     string `gorm:"type:varchar(255)"`
	DestNamespace  string `gorm:"type:varchar(255)"`
	SyncStatus     string `gorm:"type:varchar(100)"`
	HealthStatus   string `gorm:"type:varchar(100)"`
	CreatedDate    *time.Time
	SummaryImages  []string `gorm:"type:json;serializer:json"`
	ScopeConfigId  uint64
	archived.NoPKModel
}

func (ArgocdApplication) TableName() string {
	return "_tool_argocd_applications"
}

type ArgocdSyncOperation struct {
	ConnectionId    uint64 `gorm:"primaryKey"`
	ApplicationName string `gorm:"primaryKey;type:varchar(255)"`
	DeploymentId    int64  `gorm:"primaryKey"`
	Revision        string `gorm:"type:varchar(255)"`
	Kind            string `gorm:"type:varchar(100)"`
	StartedAt       *time.Time
	FinishedAt      *time.Time
	Phase           string `gorm:"type:varchar(100)"`
	Message         string `gorm:"type:text"`
	InitiatedBy     string `gorm:"type:varchar(255)"`
	SyncStatus      string `gorm:"type:varchar(100)"`
	HealthStatus    string `gorm:"type:varchar(100)"`
	ResourcesCount  int
	ContainerImages []string `gorm:"type:json;serializer:json"`
	archived.NoPKModel
}

func (ArgocdSyncOperation) TableName() string {
	return "_tool_argocd_sync_operations"
}

type ArgocdRevisionImage struct {
	ConnectionId    uint64   `gorm:"primaryKey"`
	ApplicationName string   `gorm:"primaryKey;type:varchar(255)"`
	Revision        string   `gorm:"primaryKey;type:varchar(255)"`
	Images          []string `gorm:"type:json;serializer:json"`
	archived.NoPKModel
}

func (ArgocdRevisionImage) TableName() string {
	return "_tool_argocd_revision_images"
}

type ArgocdScopeConfig struct {
	archived.ScopeConfig `mapstructure:",squash" json:",inline" gorm:"embedded"`
	ConnectionId         uint64   `gorm:"index"`
	Name                 string   `gorm:"type:varchar(255);index"`
	Entities             []string `gorm:"type:json;serializer:json"`
	DeploymentPattern    string   `gorm:"type:varchar(255)"`
	ProductionPattern    string   `gorm:"type:varchar(255)"`
	EnvNamePattern       string   `gorm:"type:varchar(255)"`
}

func (ArgocdScopeConfig) TableName() string {
	return "_tool_argocd_scope_configs"
}
