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
	"encoding/json"
	"time"

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
	PluginName   string          `json:"pluginName" mapstructure:"pluginName" gorm:"primaryKey;type:varchar(255)" validate:"required"`
	PluginOption json.RawMessage `json:"pluginOption" mapstructure:"pluginOption" gorm:"type:json"`
	Enable       bool            `json:"enable" mapstructure:"enable" gorm:"type:boolean"`
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
	Metrics     []*BaseMetric `json:"metrics" mapstructure:"metrics"`
	Blueprint   *Blueprint    `json:"blueprint" mapstructure:"blueprint"`
}

type ApiOutputProject struct {
	Project      `mapstructure:",squash"`
	Metrics      []*BaseMetric `json:"metrics" mapstructure:"metrics"`
	Blueprint    *Blueprint    `json:"blueprint" mapstructure:"blueprint"`
	LastPipeline *Pipeline     `json:"lastPipeline,omitempty" mapstructure:"lastPipeline"`
}

type ApiProjectCheck struct {
	Exist bool `json:"exist" mapstructure:"exist"`
}

type Store struct {
	StoreKey   string          `gorm:"primaryKey;type:varchar(255)"`
	StoreValue json.RawMessage `gorm:"type:json;serializer:json"`
	CreatedAt  time.Time       `json:"createdAt" mapstructure:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt" mapstructure:"updatedAt"`
}

func (Store) TableName() string {
	return "_devlake_store"
}

type ProjectScopeOutput struct {
	Projects []ProjectScope `json:"projects"`
	Count    int            `json:"count"`
}

type ProjectScope struct {
	Name        string `json:"name"`
	BlueprintId uint64 `json:"blueprintId"`
	Scopes      []struct {
		ScopeID   string `json:"scopeId"`
		ScopeName string `json:"scopeName"`
	} `json:"scopes"`
}
