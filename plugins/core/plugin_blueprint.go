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
	Plugin   string                 `json:"plugin" binding:"required"`
	Subtasks []string               `json:"subtasks"`
	Options  map[string]interface{} `json:"options"`
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
	Scope        []*BlueprintScopeV100 `json:"scope" validate:"required"`
}

// BlueprintScopeV100 is the scope definition for protocol v1.0.0
type BlueprintScopeV100 struct {
	Entities       []string        `json:"entities"`
	Options        json.RawMessage `json:"options"`
	Transformation json.RawMessage `json:"transformation"`
}

/* PluginBlueprintV200 for project support */

// Scope represents the top level entity for a data source, i.e. github repo, gitlab project, jira board.
// They turn into repo, board in Domain Layer.
// In Apache Devlake, a Project is essentially a set of these top level entities, for the framework to
// maintain these relationships dynamically and automatically, all Domain Layer Top Level Entities should
// implement this interface
type Scope interface {
	ScopeId() string
	ScopeName() string
	TableName() string
}

// PluginBlueprintV200 extends the V100 to provide support for Project to support complex metrics
// like DORA
type PluginBlueprintV200 interface {
	MakePipelinePlan(scopes []*BlueprintScopeV200) (PipelinePlan, []Scope, errors.Error)
}

// BlueprintScopeV200 contains the Plugin name and related ScopeIds, connectionId and transformationRuleId should be
// deduced by the ScopeId
type BlueprintScopeV200 struct {
	Plugin   string `json:"plugin" validate:"required"`
	ScopeIds []string
}
