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

package core

import (
	"encoding/json"
	"github.com/apache/incubator-devlake/errors"
)

// PipelineTask represents a smallest unit of execution inside a PipelinePlan
type PipelineTask struct {
	// Plugin name
	Plugin     string                 `json:"plugin" binding:"required"`
	SkipOnFail bool                   `json:"skipOnFail"`
	Subtasks   []string               `json:"subtasks"`
	Options    map[string]interface{} `json:"options"`
}

// PipelineStage consist of multiple PipelineTasks, they will be executed in parallel
type PipelineStage []*PipelineTask

// PipelinePlan consist of multiple PipelineStages, they will be executed in sequential order
type PipelinePlan []PipelineStage

// PluginBlueprintV100 is used to support Blueprint Normal model, for Plugin and Blueprint to
// collaboarte and generate a sophisticated Pipeline Plan based on User Settings.
// V100 doesn't support Project, and being deprecated, please use PluginBlueprintV200 instead
type PluginBlueprintV100 interface {
	// MakePipelinePlan generates `pipeline.tasks` based on `version` and `scope`
	//
	// `version` semver from `blueprint.settings.version`
	// `scope` arbitrary json.RawMessage, depends on `version`, for v1.0.0, it is an Array of Objects
	MakePipelinePlan(connectionId uint64, scope []*BlueprintScopeV100) (PipelinePlan, errors.Error)
}

// BlueprintConnectionV100 is the connection definition for protocol v1.0.0
type BlueprintConnectionV100 struct {
	Plugin       string                `json:"plugin" validate:"required"`
	ConnectionId uint64                `json:"connectionId" validate:"required"`
	SkipOnFail   bool                  `json:"skipOnFail"`
	Scope        []*BlueprintScopeV100 `json:"scope" validate:"required"`
}

// BlueprintScopeV100 is the scope definition for protocol v1.0.0
type BlueprintScopeV100 struct {
	Entities       []string        `json:"entities"`
	Options        json.RawMessage `json:"options"`
	Transformation json.RawMessage `json:"transformation"`
}

/*
PluginBlueprintV200 for project support


step 1: blueprint.settings like
	{
		"version": "2.0.0",
		"connections": [
			{
				"plugin": "github",
				"connectionId": 123,
				"scopes": [
					{ "id": null, "name": "apache/incubator-devlake" }
				]
			}
		]
	}

step 2: call plugin PluginBlueprintV200.MakePipelinePlan with
	[
		{ "id": "1", "name": "apache/incubator-devlake" }
	]
	plugin would return PipelinePlan like the following json, and config-ui should use scopeName for displaying
	[
		[
			{ "plugin": "github", "options": { "scopeId": "1", "scopeName": "apache/incubator-devlake" } }
		]
	]
	and []Scope for project_mapping:
	[
		&Repo{ "id": "github:GithubRepo:1:1", "name": "apache/incubator-devlake" },
		&Board{ "id": "github:GithubRepo:1:1", "name": "apache/incubator-devlake" }
	]

step 3: framework should maintain the project_mapping table based on the []Scope array
	[
		{ "projectName": "xxx", "table": "repos", "rowId": "github:GithubRepo:1:1" },
		{ "projectName": "xxx", "table": "boards", "rowId": "github:GithubRepo:1:1" },
	]
*/

// Scope represents the top level entity for a data source, i.e. github repo,
// gitlab project, jira board. They turn into repo, board in Domain Layer. In
// Apache Devlake, a Project is essentially a set of these top level entities,
// for the framework to maintain these relationships dynamically and
// automatically, all Domain Layer Top Level Entities should implement this
// interface
type Scope interface {
	ScopeId() string
	ScopeName() string
	TableName() string
}

// DataSourcePluginBlueprintV200 extends the V100 to provide support for
// Project, so that complex metrics like DORA can be implemented based on a set
// of Data Scopes
type DataSourcePluginBlueprintV200 interface {
	MakeDataSourcePipelinePlanV200(connectionId uint64, scopes []*BlueprintScopeV200) (PipelinePlan, []Scope, errors.Error)
}

// BlueprintConnectionV200 contains the pluginName/connectionId  and related Scopes,
type BlueprintConnectionV200 struct {
	Plugin       string                `json:"plugin" validate:"required"`
	ConnectionId uint64                `json:"connectionId" validate:"required"`
	Scopes       []*BlueprintScopeV200 `json:"scopes" validate:"required"`
}

// BlueprintScopeV200 contains the `id` and `name` for a specific scope
// transformationRuleId should be deduced by the ScopeId
type BlueprintScopeV200 struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// MetricPluginBlueprintV200 is similar to the DataSourcePluginBlueprintV200
// but for Metric Plugin, take dora as an example, it doens't have any scope,
// nor does it produce any, however, it does require other plugin to be
// executed beforehand, like calcuating refdiff before it can connect PR to the
// right Deployment keep in mind it would be called IFF the plugin was enabled
// for the project.
type MetricPluginBlueprintV200 interface {
	MakeMetricPluginPipelinePlanV200(projectName string, options json.RawMessage) (PipelinePlan, errors.Error)
}

// CompositeDataSourcePluginBlueprintV200 is for unit test
type CompositeDataSourcePluginBlueprintV200 interface {
	PluginMeta
	DataSourcePluginBlueprintV200
}

// CompositeMetricPluginBlueprintV200 is for unit test
type CompositeMetricPluginBlueprintV200 interface {
	PluginMeta
	MetricPluginBlueprintV200
}

// CompositeMetricPluginBlueprintV200 is for unit test
type CompositePluginBlueprintV200 interface {
	PluginMeta
	DataSourcePluginBlueprintV200
	MetricPluginBlueprintV200
}
