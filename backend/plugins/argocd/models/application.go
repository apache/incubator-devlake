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
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.ToolLayerScope = (*ArgocdApplication)(nil)

type ArgocdApplication struct {
	common.Scope   `mapstructure:",squash"`
	Name           string     `gorm:"type:varchar(255);primaryKey" json:"name" mapstructure:"name" validate:"required"`
	Namespace      string     `gorm:"type:varchar(255)" json:"namespace" mapstructure:"namespace"`
	Project        string     `gorm:"type:varchar(255)" json:"project" mapstructure:"project"`
	RepoURL        string     `gorm:"type:varchar(500)" json:"repoUrl" mapstructure:"repoUrl"`
	Path           string     `gorm:"type:varchar(255)" json:"path" mapstructure:"path"`
	TargetRevision string     `gorm:"type:varchar(255)" json:"targetRevision" mapstructure:"targetRevision"`
	DestServer     string     `gorm:"type:varchar(255)" json:"destServer" mapstructure:"destServer"`
	DestNamespace  string     `gorm:"type:varchar(255)" json:"destNamespace" mapstructure:"destNamespace"`
	SyncStatus     string     `gorm:"type:varchar(100)" json:"syncStatus" mapstructure:"syncStatus"`     // Synced, OutOfSync, Unknown
	HealthStatus   string     `gorm:"type:varchar(100)" json:"healthStatus" mapstructure:"healthStatus"` // Healthy, Progressing, Degraded, Suspended, Missing, Unknown
	CreatedDate    *time.Time `json:"createdDate,omitempty" mapstructure:"createdDate,omitempty"`
	SummaryImages  []string   `gorm:"type:json;serializer:json" json:"summaryImages" mapstructure:"summaryImages"`
	common.NoPKModel
}

func (ArgocdApplication) TableName() string {
	return "_tool_argocd_applications"
}

func (a ArgocdApplication) ScopeId() string {
	return a.Name
}

func (a ArgocdApplication) ScopeName() string {
	return a.Name
}

func (a ArgocdApplication) ScopeFullName() string {
	return a.Name
}

func (a ArgocdApplication) ScopeParams() interface{} {
	return &ArgocdApiParams{
		ConnectionId: a.ConnectionId,
		Name:         a.Name,
	}
}

type ArgocdApiParams struct {
	ConnectionId uint64
	Name         string
}
