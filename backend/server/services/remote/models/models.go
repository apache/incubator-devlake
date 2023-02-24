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
)

const (
	PythonPoetryCmd PluginType      = "python-poetry"
	PythonCmd       PluginType      = "python"
	None            PluginExtension = ""
	Metric          PluginExtension = "metric"
	Datasource      PluginExtension = "datasource"
)

type (
	PluginType      string
	PluginExtension string
)

type PluginInfo struct {
	Type                     PluginType      `json:"type" validate:"required"`
	Name                     string          `json:"name" validate:"required"`
	Extension                PluginExtension `json:"extension"`
	ConnectionSchema         map[string]any  `json:"connection_schema" validate:"required"`
	TransformationRuleSchema map[string]any  `json:"transformation_rule_schema" validate:"required"`
	ScopeInfo                ScopeInfo       `json:"scope_info" validate:"dive"`
	Description              string          `json:"description"`
	PluginPath               string          `json:"plugin_path" validate:"required"`
	SubtaskMetas             []SubtaskMeta   `json:"subtask_metas" validate:"dive"`
}

type ScopeInfo struct {
	TableName   string         `json:"table_name" validate:"required"`
	ScopeSchema map[string]any `json:"scope_schema" validate:"required"`
}

type ScopeModel struct {
	Id                   string    `gorm:"primaryKey" json:"id"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	ConnectionId         uint64    `gorm:"primaryKey" json:"connection_id"`
	TransformationRuleId uint64    `json:"transformation_rule_id"`
}

type TransformationModel struct {
	Id        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string
}

type SubtaskMeta struct {
	Name             string   `json:"name" validate:"required"`
	EntryPointName   string   `json:"entry_point_name" validate:"required"`
	Arguments        []string `json:"arguments"`
	Required         bool     `json:"required"`
	EnabledByDefault bool     `json:"enabled_by_default"`
	Description      string   `json:"description" validate:"required"`
	DomainTypes      []string `json:"domain_types" validate:"required"`
}
