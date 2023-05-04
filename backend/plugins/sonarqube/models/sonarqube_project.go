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
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

var _ plugin.ToolLayerScope = (*SonarqubeProject)(nil)
var _ plugin.ApiScope = (*SonarqubeApiProject)(nil)

type SonarqubeProject struct {
	common.NoPKModel `json:"-" mapstructure:"-"`
	ConnectionId     uint64           `json:"connectionId" validate:"required" gorm:"primaryKey" mapstructure:"connectionId"`
	ProjectKey       string           `json:"projectKey" validate:"required" gorm:"type:varchar(255);primaryKey" mapstructure:"projectKey"`
	Name             string           `json:"name" gorm:"type:varchar(255)" mapstructure:"name"`
	Qualifier        string           `json:"qualifier" gorm:"type:varchar(255)" mapstructure:"qualifier"`
	Visibility       string           `json:"visibility" gorm:"type:varchar(64)" mapstructure:"visibility"`
	LastAnalysisDate *api.Iso8601Time `json:"lastAnalysisDate" mapstructure:"lastAnalysisDate"`
	Revision         string           `json:"revision" gorm:"type:varchar(128)" mapstructure:"revision"`
}

func (SonarqubeProject) TableName() string {
	return "_tool_sonarqube_projects"
}

func (p SonarqubeProject) ScopeId() string {
	return p.ProjectKey
}

func (p SonarqubeProject) ScopeName() string {
	return p.Name
}

type SonarqubeApiProject struct {
	ProjectKey       string           `json:"key"`
	Name             string           `json:"name"`
	Qualifier        string           `json:"qualifier"`
	Visibility       string           `json:"visibility"`
	LastAnalysisDate *api.Iso8601Time `json:"lastAnalysisDate"`
	Revision         string           `json:"revision"`
}

// Convert the API response to our DB model instance
func (sonarqubeApiProject SonarqubeApiProject) ConvertApiScope() plugin.ToolLayerScope {
	sonarqubeProject := SonarqubeProject{
		ProjectKey:       sonarqubeApiProject.ProjectKey,
		Name:             sonarqubeApiProject.Name,
		Qualifier:        sonarqubeApiProject.Qualifier,
		Visibility:       sonarqubeApiProject.Visibility,
		LastAnalysisDate: sonarqubeApiProject.LastAnalysisDate,
		Revision:         sonarqubeApiProject.Revision,
	}
	return sonarqubeProject
}
