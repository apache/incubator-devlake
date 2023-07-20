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
)

type BaseProject struct {
	Name        string `json:"name" mapstructure:"name" gorm:"primaryKey;type:varchar(255)" validate:"required"`
	Description string `json:"description" mapstructure:"description" gorm:"type:text"`
}

type Project struct {
	BaseProject `mapstructure:",squash"`
	common.NoPKModel
}

func (Project) TableName() string {
	return "projects"
}

type BaseMetric struct {
	PluginName   string `json:"pluginName" mapstructure:"pluginName" gorm:"primaryKey;type:varchar(255)" validate:"required"`
	PluginOption string `json:"pluginOption" mapstructure:"pluginOption" gorm:"type:text"`
	Enable       bool   `json:"enable" mapstructure:"enable" gorm:"type:boolean"`
}

type BaseProjectMetricSetting struct {
	ProjectName string `json:"projectName" mapstructure:"projectName" gorm:"primaryKey;type:varchar(255)"`
	BaseMetric  `mapstructure:",squash"`
}

type ProjectMetricSetting struct {
	BaseProjectMetricSetting `mapstructure:",squash"`
	common.NoPKModel
}

func (ProjectMetricSetting) TableName() string {
	return "project_metric_settings"
}

type ApiInputProject struct {
	BaseProject `mapstructure:",squash"`
	Enable      *bool         `json:"enable" mapstructure:"enable"`
	Metrics     *[]BaseMetric `json:"metrics" mapstructure:"metrics"`
}

type ApiOutputProject struct {
	BaseProject    `mapstructure:",squash"`
	Metrics        *[]BaseMetric `json:"metrics" mapstructure:"metrics"`
	Blueprint      *Blueprint    `json:"blueprint" mapstructure:"blueprint"`
	LatestPipeLine *Pipeline     `json:"latest_pipeline,omitempty" mapstructure:"latest_pipeline"`
}
